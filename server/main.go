package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rezan/rpcs/pb"
	_ "modernc.org/sqlite"
	"google.golang.org/grpc"
)

const (
	priceHistoryTable = `CREATE TABLE IF NOT EXISTS price_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		symbol TEXT NOT NULL,
		price REAL NOT NULL,
		at_ts INTEGER NOT NULL,
		created_at INTEGER DEFAULT (strftime('%s','now'))
	);`
)

// priceDB wraps SQLite for price history (optional persistence).
type priceDB struct {
	db *sql.DB
}

func openDB(dataSource string) (*priceDB, error) {
	if dataSource == "" {
		dataSource = "file:price.db?_pragma=journal_mode(wal)"
	}
	db, err := sql.Open("sqlite", dataSource)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &priceDB{db: db}, nil
}

func (d *priceDB) migrate() error {
	_, err := d.db.Exec(priceHistoryTable)
	return err
}

// latestPrice returns the most recent price for symbol, or 0 and false if none.
func (d *priceDB) latestPrice(symbol string) (float64, bool) {
	var price float64
	err := d.db.QueryRow(
		`SELECT price FROM price_history WHERE symbol = ? ORDER BY at_ts DESC LIMIT 1`,
		symbol,
	).Scan(&price)
	if err == sql.ErrNoRows {
		return 0, false
	}
	if err != nil {
		log.Printf("[db] latestPrice %q: %v", symbol, err)
		return 0, false
	}
	return price, true
}

// insertPrice records a price update (non-blocking; logs errors).
func (d *priceDB) insertPrice(symbol string, price float64, atTs int64) {
	_, err := d.db.Exec(
		`INSERT INTO price_history (symbol, price, at_ts) VALUES (?, ?, ?)`,
		symbol, price, atTs,
	)
	if err != nil {
		log.Printf("[db] insertPrice: %v", err)
	}
}

func (d *priceDB) close() error { return d.db.Close() }

// priceBroadcaster holds the current price and notifies all subscribed streams when it changes.
type priceBroadcaster struct {
	mu          sync.RWMutex
	price       float64
	symbol      string
	subscribers map[chan *pb.PriceUpdate]struct{}
}

func newPriceBroadcaster(symbol string, initialPrice float64) *priceBroadcaster {
	return &priceBroadcaster{
		symbol:      symbol,
		price:       initialPrice,
		subscribers: make(map[chan *pb.PriceUpdate]struct{}),
	}
}

func (b *priceBroadcaster) getCurrent() *pb.PriceUpdate {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return &pb.PriceUpdate{
		Symbol: b.symbol,
		Price:  b.price,
		AtTs:   time.Now().Unix(),
	}
}

func (b *priceBroadcaster) setPrice(price float64) {
	b.mu.Lock()
	b.price = price
	update := &pb.PriceUpdate{
		Symbol: b.symbol,
		Price:  price,
		AtTs:   time.Now().Unix(),
	}
	subs := make([]chan *pb.PriceUpdate, 0, len(b.subscribers))
	for ch := range b.subscribers {
		subs = append(subs, ch)
	}
	b.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- update:
		default:
			// client too slow, skip this update
		}
	}
}

func (b *priceBroadcaster) subscribe() chan *pb.PriceUpdate {
	ch := make(chan *pb.PriceUpdate, 8)
	b.mu.Lock()
	b.subscribers[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *priceBroadcaster) unsubscribe(ch chan *pb.PriceUpdate) {
	b.mu.Lock()
	delete(b.subscribers, ch)
	b.mu.Unlock()
	// dont close(ch) server side should not close the channel we could get a panic from the client
	// The server could be close if client is closed or not connected anymore. 
}

// priceService implements pb.PriceServiceServer.
type priceService struct {
	pb.UnimplementedPriceServiceServer
	broadcaster *priceBroadcaster
	db          *priceDB
}

func (s *priceService) GetCurrentPrice(ctx context.Context, req *pb.GetCurrentPriceRequest) (*pb.PriceUpdate, error) {
	return s.broadcaster.getCurrent(), nil
}

func (s *priceService) SetPrice(ctx context.Context, req *pb.SetPriceRequest) (*pb.PriceUpdate, error) {
	if req.Price <= 0 {
		return s.broadcaster.getCurrent(), nil // ignore invalid price; return current
	}
	atTs := time.Now().Unix()
	s.broadcaster.setPrice(req.Price)
	if s.db != nil {
		sym := req.Symbol
		if sym == "" {
			sym = "BTC"
		}
		s.db.insertPrice(sym, req.Price, atTs)
	}
	return s.broadcaster.getCurrent(), nil
}

func (s *priceService) SubscribePriceUpdates(req *pb.SubscribeRequest, stream pb.PriceService_SubscribePriceUpdatesServer) error {
	ch := s.broadcaster.subscribe()
	defer s.broadcaster.unsubscribe(ch)

	// Send current price immediately
	if err := stream.Send(s.broadcaster.getCurrent()); err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case update, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(update); err != nil {
				return err
			}
		}
	}
}

// runPriceSimulator changes the price every few seconds so clients see real-time updates.
// If db is non-nil, each price update is persisted to the database.
func runPriceSimulator(b *priceBroadcaster, db *priceDB) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		b.mu.RLock()
		current := b.price
		b.mu.RUnlock()
		change := (rand.Float64() - 0.5) * 0.01 * current
		newPrice := current + change
		if newPrice < 0.01 {
			newPrice = 0.01
		}
		atTs := time.Now().Unix()
		b.setPrice(newPrice)
		if db != nil {
			go db.insertPrice(b.symbol, newPrice, atTs)
		}
		log.Printf("[simulator] price updated: %.2f", newPrice)
	}
}

// HTTP API for price: GET = read current, POST = set price (persisted to DB). Use from terminal with curl.
type priceAPI struct {
	broadcaster *priceBroadcaster
	db          *priceDB
}

func (h *priceAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		cur := h.broadcaster.getCurrent()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"symbol": cur.Symbol,
			"price":  cur.Price,
			"at_ts":  cur.AtTs,
		})
	case http.MethodPost, http.MethodPut:
		body, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to read body"})
			return
		}
		var payload struct {
			Symbol string  `json:"symbol"`
			Price  float64 `json:"price"`
		}
		payload.Symbol = "BTC"
		if err := json.Unmarshal(body, &payload); err != nil {
			trimmed := strings.TrimSpace(string(body))
			if p, errP := strconv.ParseFloat(trimmed, 64); errP == nil && p > 0 {
				payload.Price = p
			}
		}
		if payload.Symbol == "" {
			payload.Symbol = "BTC"
		}
		if payload.Price <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "price must be > 0"})
			return
		}
		atTs := time.Now().Unix()
		h.broadcaster.setPrice(payload.Price)
		if h.db != nil {
			h.db.insertPrice(payload.Symbol, payload.Price, atTs)
		}
		cur := h.broadcaster.getCurrent()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"symbol": cur.Symbol,
			"price":  cur.Price,
			"at_ts":  cur.AtTs,
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "use GET or POST"})
	}
}

func runHTTPAPI(broadcaster *priceBroadcaster, db *priceDB) {
	api := &priceAPI{broadcaster: broadcaster, db: db}
	http.Handle("/price", api)
	http.Handle("/api/price", api)
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Printf("[http] API listening on %s (GET/POST /price)", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("[http] serve: %v", err)
	}
}

func main() {
	// Optional: DB file path (default: price.db in current directory)
	dbPath := os.Getenv("PRICE_DB")
	if dbPath == "" {
		dbPath = "file:price.db?_pragma=journal_mode(wal)"
	}
	// Resolve to absolute path so we always write in a known place
	if dbPath == "file:price.db?_pragma=journal_mode(wal)" {
		if cwd, _ := os.Getwd(); cwd != "" {
			dbPath = "file:" + filepath.Join(cwd, "price.db") + "?_pragma=journal_mode(wal)"
		}
	}

	db, err := openDB(dbPath)
	if err != nil {
		log.Fatalf("open DB: %v", err)
	}
	defer db.close()
	if err := db.migrate(); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("[db] SQLite connected, price_history table ready")

	symbol := "BTC"
	initialPrice := 50000.0
	if p, ok := db.latestPrice(symbol); ok && p > 0 {
		initialPrice = p
		log.Printf("[db] loaded last price for %s: %.2f", symbol, initialPrice)
	}

	broadcaster := newPriceBroadcaster(symbol, initialPrice)
	go runPriceSimulator(broadcaster, db)
	go runHTTPAPI(broadcaster, db)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterPriceServiceServer(srv, &priceService{broadcaster: broadcaster, db: db})

	log.Println("gRPC server listening on :50051 (real-time price updates)")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

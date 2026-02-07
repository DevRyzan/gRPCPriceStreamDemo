package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/rezan/rpcs/pb"
	"google.golang.org/grpc"
)

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
	close(ch)
}

// priceService implements pb.PriceServiceServer.
type priceService struct {
	pb.UnimplementedPriceServiceServer
	broadcaster *priceBroadcaster
}

func (s *priceService) GetCurrentPrice(ctx context.Context, req *pb.GetCurrentPriceRequest) (*pb.PriceUpdate, error) {
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
func runPriceSimulator(b *priceBroadcaster) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// Simulate price movement: ±0.5% around current
		b.mu.RLock()
		current := b.price
		b.mu.RUnlock()
		change := (rand.Float64() - 0.5) * 0.01 * current
		newPrice := current + change
		if newPrice < 0.01 {
			newPrice = 0.01
		}
		b.setPrice(newPrice)
		log.Printf("[simulator] price updated: %.2f", newPrice)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	broadcaster := newPriceBroadcaster("BTC", 50000.0)
	go runPriceSimulator(broadcaster)

	srv := grpc.NewServer()
	pb.RegisterPriceServiceServer(srv, &priceService{broadcaster: broadcaster})

	log.Println("gRPC server listening on :50051 (real-time price updates)")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

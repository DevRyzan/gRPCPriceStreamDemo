package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rezan/rpcs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	addr := "localhost:50051"
	if a := os.Getenv("GRPC_ADDR"); a != "" {
		addr = a
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewPriceServiceClient(conn)
	ctx := context.Background()

	// Optional: set price via env (e.g. SET_PRICE=60000) then continue to subscribe
	if v := os.Getenv("SET_PRICE"); v != "" {
		if price, err := strconv.ParseFloat(v, 64); err == nil && price > 0 {
			setCtx, setCancel := context.WithTimeout(ctx, 5*time.Second)
			resp, err := client.SetPrice(setCtx, &pb.SetPriceRequest{Symbol: "BTC", Price: price})
			setCancel()
			if err != nil {
				log.Printf("SetPrice: %v", err)
			} else {
				log.Printf("SetPrice OK: %s = %.2f (at %d)", resp.Symbol, resp.Price, resp.AtTs)
			}
		}
	}

	ctx5, cancel := context.WithTimeout(ctx, 5*time.Second)
	cur, err := client.GetCurrentPrice(ctx5, &pb.GetCurrentPriceRequest{})
	cancel()
	if err != nil {
		log.Printf("GetCurrentPrice: %v (continuing anyway)", err)
	} else {
		log.Printf("Current price (one-shot): %s = %.2f at %d", cur.Symbol, cur.Price, cur.AtTs)
	}

	
	streamCtx, streamCancel := context.WithCancel(context.Background()) // Subscribe to real-time price updates 
	stream, err := client.SubscribePriceUpdates(streamCtx, &pb.SubscribeRequest{Symbol: "BTC"})
	if err != nil {
		log.Fatalf("SubscribePriceUpdates: %v", err)
	}

	
	quit := make(chan os.Signal, 1)// Graceful shutdown on Ctrl+C 
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down...")
		streamCancel()
	}()

	log.Println("Subscribed to real-time price updates (Ctrl+C to stop):")
	for {
		update, err := stream.Recv()
		if err != nil {
			log.Printf("Stream ended: %v", err)
			return
		}
		log.Printf("  [live] %s = %.2f  (at %d)", update.Symbol, update.Price, update.AtTs)
	}
}

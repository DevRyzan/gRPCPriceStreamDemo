package main

import (
	"context"
	"log"
	"os"
	"os/signal"
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

 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := client.GetCurrentPrice(ctx, &pb.GetCurrentPriceRequest{})
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

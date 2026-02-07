package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/rezan/rpcs/pb"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// App holds state for the desktop app and exposes bindings to the frontend.
type App struct {
	ctx    context.Context
	mu     sync.Mutex
	cancel context.CancelFunc
}

// NewApp creates the application struct (called from main.go).
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts (Wails lifecycle).
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// StartPriceStream connects to the gRPC server and streams price updates to the frontend via events.
// serverAddr defaults to "localhost:50051" if empty.
func (a *App) StartPriceStream(serverAddr string) error {
	a.mu.Lock()
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
	ctx := a.ctx
	if ctx == nil {
		a.mu.Unlock()
		return nil
	}
	if serverAddr == "" {
		serverAddr = os.Getenv("GRPC_ADDR")
		if serverAddr == "" {
			serverAddr = "localhost:50051"
		}
	}
	streamCtx, cancel := context.WithCancel(ctx)
	a.cancel = cancel
	a.mu.Unlock()

	go func() {
		conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("[desktop] gRPC dial: %v", err)
			runtime.EventsEmit(a.ctx, "priceError", err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewPriceServiceClient(conn)
		stream, err := client.SubscribePriceUpdates(streamCtx, &pb.SubscribeRequest{Symbol: "BTC"})
		if err != nil {
			log.Printf("[desktop] SubscribePriceUpdates: %v", err)
			runtime.EventsEmit(a.ctx, "priceError", err.Error())
			return
		}

		for {
			update, err := stream.Recv()
			if err != nil {
				log.Printf("[desktop] stream ended: %v", err)
				return
			}
			payload := map[string]interface{}{
				"symbol": update.Symbol,
				"price":  update.Price,
				"atTs":   update.AtTs,
			}
			runtime.EventsEmit(a.ctx, "price", payload)
		}
	}()

	return nil
}

// StopPriceStream stops the current price stream.
func (a *App) StopPriceStream() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
}

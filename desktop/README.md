# Price Desktop (Wails)

Desktop app that connects to the gRPC price server and shows real-time price updates.

## Prerequisites

- Go (with `go` in PATH)
- Node.js 18+ and npm
- [Wails v2 CLI](https://wails.io/docs/gettingstarted/installation) installed:  
  `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## Run

1. **Start the price server** (from repo root):

   ```bash
   go run ./server
   ```

   Server listens on gRPC `:50051` and HTTP API `:8080`.

2. **From the `desktop` folder:**

   ```bash
   cd desktop
   npm install
   cd frontend
   npm install
   cd ..
   wails dev
   ```

   On first run, Wails may regenerate the frontend bindings in `frontend/src/wailsjs/`.  
   In the app window, click **Connect** to subscribe to the price stream. You should see live prices from the server.

## Build for production

```bash
cd desktop
wails build
```

Output is in `desktop/build/bin/` (e.g. `price-desktop.exe` on Windows).

## Optional

- Set `GRPC_ADDR` (e.g. to `host:50051`) if the server is not on `localhost:50051`.
- The app uses the same gRPC `SubscribePriceUpdates` stream as the Go client; no browser or Envoy needed.

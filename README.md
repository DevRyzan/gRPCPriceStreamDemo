# Real-time gRPC Price Service (Go)

A sample gRPC service that streams **real-time price updates** to clients. When the server changes the price (simulated every 2 seconds), every subscribed client receives the update immediately over the same stream.

## Fixing "slice bounds out of range" panic

If you see that panic when running the client or server, the generated `pb/` code is out of date. Regenerate it (no need to install protoc; Buf is used automatically):

```powershell
# Install codegen plugins once (puts them in %USERPROFILE%\go\bin)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# From project root – regenerates pb/price.pb.go and pb/price_grpc.pb.go
.\generate.ps1
```

Then run `go run .\server\` and `go run .\client\` again.

## Quick start

**Terminal 1 – start the server**

```bash
go run ./server
```

**Terminal 2 – run the client (subscribes and prints each price update)**

```bash
go run ./client
```

You should see the client printing new prices every ~2 seconds as the server updates them.

## What’s in the repo

| Path | Description |
|------|-------------|
| `proto/price.proto` | Service definition: `GetCurrentPrice` (unary) and `SubscribePriceUpdates` (server stream) |
| `pb/` | **Generated only** – `price.pb.go` and `price_grpc.pb.go` from Buf/protoc (run `.\generate.ps1` to regenerate) |
| `server/main.go` | gRPC server: current price, broadcast to subscribers, price simulator was a simulator but now SQLite persistence |
| `client/main.go` | gRPC client: optional one-shot price, then subscribes and prints every update |

## Database (SQLite)

The server uses **SQLite** (pure Go, no CGO) to persist price history and to load the last known price on startup.

- **On start:** Creates `price.db` in the server’s working directory (and `price_history` table if missing), then loads the latest stored price for the symbol and uses it as the initial value.
- **On each update:** Appends a row to `price_history` (`symbol`, `price`, `at_ts`) so you get a history of all simulated updates.

To use a different DB file, set `PRICE_DB` (e.g. `PRICE_DB=file:/path/to/other.db?_pragma=journal_mode(wal)`).

Run `go mod tidy` once to fetch the SQLite driver (`modernc.org/sqlite`).

## How the real-time stream works

1. **Server** keeps a single “current price” and a set of subscriber channels.
2. **SubscribePriceUpdates** is a **server-streaming RPC**: the client opens one stream and the server keeps the connection open.
3. When the server calls **setPrice** (here triggered by a timer every 2 seconds), it pushes the new `PriceUpdate` to every subscribed stream.
4. The **client** blocks on `stream.Recv()` and prints each update as it arrives.

So “real-time” here means: **one long-lived gRPC stream per client, and the server pushes a new message whenever the price changes.**
 

To regenerate with **protoc** instead of Buf: install [protoc](https://github.com/protocolbuffers/protobuf/releases), then run `.\generate.ps1`.

## HTTP API (request / change price from the terminal)

The server also exposes an HTTP API on **:8080** (or `HTTP_ADDR`) so you can read and set the price from the terminal. Changes are persisted to the DB and gRPC subscribers get the update.

**PowerShell** (Windows – `curl` is an alias for `Invoke-WebRequest` and uses different syntax):

```powershell
# Get current price
Invoke-RestMethod -Uri http://localhost:8080/price

# Set price (JSON body)
Invoke-RestMethod -Uri http://localhost:8080/price -Method Post -ContentType "application/json" -Body '{"price":60000}'

# Set price (raw number in body)
Invoke-RestMethod -Uri http://localhost:8080/price -Method Post -ContentType "text/plain" -Body "60000"
```

**curl** (Bash/WSL, or use `curl.exe` on Windows to avoid the PowerShell alias):

```bash
# Get current price
curl http://localhost:8080/price

# Set price
curl -X POST http://localhost:8080/price -H "Content-Type: application/json" -d "{\"price\":60000}"
curl -X POST http://localhost:8080/price -d "60000"
```


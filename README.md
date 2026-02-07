# Real-time gRPC Price Service (Go)

A sample gRPC service that streams **real-time price updates** to clients. When the server changes the price (simulated every 2 seconds), every subscribed client receives the update immediately over the same stream.

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
| `pb/` | Generated Go code for messages and gRPC service (can be regenerated with `protoc`) |
| `server/main.go` | gRPC server: holds current price, broadcasts to all subscribers, runs a price simulator |
| `client/main.go` | gRPC client: optional one-shot price, then subscribes and prints every update |

## How the real-time stream works

1. **Server** keeps a single “current price” and a set of subscriber channels.
2. **SubscribePriceUpdates** is a **server-streaming RPC**: the client opens one stream and the server keeps the connection open.
3. When the server calls **setPrice** (here triggered by a timer every 2 seconds), it pushes the new `PriceUpdate` to every subscribed stream.
4. The **client** blocks on `stream.Recv()` and prints each update as it arrives.

So “real-time” here means: **one long-lived gRPC stream per client, and the server pushes a new message whenever the price changes.**
 

1. **Install the Protocol Buffers compiler** (if needed):
   - Windows: [protoc release](https://github.com/protocolbuffers/protobuf/releases) (e.g. `protoc-27.x-win64.zip` → unzip and add `bin` to PATH), or `choco install protobuf`
   - Or use [Buf](https://buf.build/docs/installation) and run `buf generate` with a `buf.gen.yaml` that uses the Go plugins.

2. **Install the Go codegen plugins** (once):
   ```powershell
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. **From the project root**, run:
   ```powershell
   .\generate.ps1
   ```

   Or run `protoc` yourself:
   ```powershell
   protoc --proto_path=. --go_out=pb --go_opt=module=github.com/rezan/rpcs --go-grpc_out=pb --go-grpc_opt=module=github.com/rezan/rpcs proto/price.proto
   ```

Then run `go run .\client\` again (with the server running).


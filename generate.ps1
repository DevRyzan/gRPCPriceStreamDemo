$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $root

# Put Go bin on PATH so buf can find protoc-gen-go and protoc-gen-go-grpc
$gobin = go env GOPATH
if ($gobin) { $gobin = Join-Path $gobin "bin" }
$goBin = go env GOBIN
if ($goBin) { $gobin = $goBin }
if ($gobin -and (Test-Path $gobin)) { $env:Path = "$gobin;$env:Path" }

# Try Buf first (no separate protoc install)
Write-Host "Trying buf generate... (first run may download buf)"
$ErrorActionPreference = "Continue"
go run github.com/bufbuild/buf/cmd/buf@latest generate 2>&1 | Out-Host
$bufOk = ($LASTEXITCODE -eq 0)
$ErrorActionPreference = "Stop"
if ($bufOk) {
    Write-Host "Generated with buf. Run: go run .\server\  and  go run .\client\"
    exit 0
}
Write-Host "Buf failed (exit $LASTEXITCODE). Trying protoc..."

# Fall back to protoc
$gobin = go env GOPATH
if ($gobin) { $gobin = Join-Path $gobin "bin" }
$goBin = go env GOBIN
if ($goBin) { $gobin = $goBin }
if ($gobin -and (Test-Path $gobin)) {
    $env:Path = "$gobin;$env:Path"
}

$proto = Join-Path $root "proto\price.proto"
$goOut = Join-Path $root "pb"
$mod = "github.com/rezan/rpcs"

if (-not (Test-Path $proto)) {
    Write-Error "Proto file not found: $proto"
}

Write-Host "Generating Go code from $proto into $goOut ..."
protoc --proto_path=$root --go_out=$goOut --go_opt=module=$mod --go-grpc_out=$goOut --go-grpc_opt=module=$mod $proto
if ($LASTEXITCODE -ne 0) {
    Write-Error "protoc failed. Install protoc from https://github.com/protocolbuffers/protobuf/releases and ensure protoc-gen-go and protoc-gen-go-grpc are in PATH (e.g. go env GOPATH)/bin"
}
Write-Host "Done. Run: go run .\server\  and  go run .\client\"

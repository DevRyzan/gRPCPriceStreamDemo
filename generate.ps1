$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $root


go run github.com/bufbuild/buf/cmd/buf@latest generate 2>$null
if ($LASTEXITCODE -eq 0) {
    Write-Host "Generated with buf. Run: go run .\server\  and  go run .\client\"
    exit 0
}

# Ensure Go plugin binaries are on PATH for protoc
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

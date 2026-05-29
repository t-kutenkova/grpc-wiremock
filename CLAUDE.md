# grpc-wiremock Project Guide

## Overview

**grpc-wiremock** is a tool for mocking server APIs that supports:
- Multiple APIs simultaneously
- Mocking of gRPC (protobuf) API
- Generating mocks by proto and OpenAPI files
- Hot reload on contract changes

It is based on [Wiremock](https://wiremock.org/docs) and includes a gRPC-to-HTTP proxy generator.

## Architecture

### Main Components

1. **Mock Generator** (`cmd/mockgen`, `pkg/generators/mocks`)
   - Generates mocks from proto and OpenAPI contracts
   - Uses Wiremock for request matching and response stubbing

2. **gRPC-to-HTTP Proxy Generator** (`cmd/grpc2http`, `pkg/generators/proxy`)
   - Generates a Go proxy that converts HTTP requests to gRPC
   - Compiles proto files to Go packages using `protoc`
   - Renders NGINX configs for routing

3. **Watcher** (`cmd/watcher`)
   - Monitors contract files for changes
   - Triggers automatic rebuilds

### Project Structure

```
.
├── cmd/                      # CLI commands
│   ├── grpc2http/           # Proxy generator CLI
│   ├── mockgen/             # Mock generator CLI
│   ├── confgen/             # Config generator CLI
│   ├── certgen/             # Certificate generator
│   ├── reload/              # NGINX reload utility
│   └── watcher/             # File watcher
├── internal/usecases/       # Business logic
├── pkg/                     # Core packages
│   ├── builder/             # Contract update logic
│   ├── compiler/            # Proto compiler wrapper
│   ├── generators/          # Template generators
│   ├── models/              # Contract models
│   ├── renderer/            # Template renderer
│   └── sourcer/             # Contract file loader
├── static/                  # Embedded static files
│   ├── proto-includes/      # Standard proto files (google/protobuf, google/rpc)
│   └── proto-annotations/   # Google API annotations
├── scripts/                 # Helper scripts
└── tests/                   # Integration tests
```

## Proto Generation Flow

1. **Load Contracts** (`pkg/sourcer`, `pkg/models/protocontract`)
   - Reads proto files from input directory
   - Parses imports and dependencies

2. **Overwrite Contracts** (`internal/usecases/grpc2http/generate.go:overwriteContracts`)
   - Updates `go_package` options for user contracts
   - Preserves standard Google protos (google/protobuf, google/rpc, google/api)

3. **Compile to Go** (`pkg/compiler`, `pkg/generators/proxy/packages.go`)
   - Runs `protoc` with `--go_out` and `--go-grpc_out` plugins
   - Uses temp directories to isolate compilation
   - Copies generated packages to output

## The Proto Registration Problem

**Problem**: `panic: proto: file "google/rpc/status.proto" is already registered`

This occurs when:
1. The project imports `google.golang.org/genproto/googleapis/rpc/status` (via go.mod)
2. Proto files import `google/rpc/status.proto` from static assets
3. Both get compiled, registering the same proto under different package paths

**Root Cause**: The `go_package` option in `static/proto-includes/google/rpc/status.proto` is:
```
option go_package = "google.golang.org/genproto/googleapis/rpc/status;status";
```

When compiled, this conflicts with the dependency `google.golang.org/genproto` already in go.mod.

## Solution Approaches

1. **Exclude standard protos from compilation**: The blacklist (`pkg/blacklist/blacklist.go`) should prevent recompiling files like `google/rpc/status.proto`

2. **Use `--go_out` with `M` option**: Map proto paths to existing Go packages

3. **Remove duplicate imports**: Ensure `google/rpc/status.proto` isn't compiled if already available via go.mod

## Key Files

- `internal/usecases/grpc2http/generate.go` - Main proxy generation flow
- `pkg/compiler/command/command.go` - Protoc command construction
- `pkg/blacklist/blacklist.go` - Files excluded from modification
- `pkg/builder/updaters/gopackage.go` - Go package path updates

## Common Tasks

### Build Docker Image
```bash
./dev/build_image.sh
```

### Run in Development
```bash
make -C ./dev/example/ compose-up compose-logs
```

### Run Tests
```bash
make unit-tests
```

## Environment Variables

- `MOCKS_PATH` - Directory for mock definitions
- `CONTRACTS_PATH` - Directory with proto/OpenAPI contracts
- `WIREMOCK_GUI_PORT` - Port for Wiremock admin UI (default: 9000)

# Build cache warmup

This directory is not part of the product. It exists so that the Docker image
ships a module cache and a build cache that the generated gRPC-to-HTTP proxy
can actually reuse at container start.

## Why

At runtime `scripts/proxy/install.sh` generates a proxy project from
`static/proxy/template/layout` and runs `make install` in it. That build must
hit the caches baked into the image, otherwise it spends ~30s downloading and
recompiling dependencies on every container start and on every hot reload.

A cache hit requires the *same module versions*, so the warmup is built against
`static/proxy/template/layout/go.mod` + `go.sum` — the exact manifests the
generated proxy gets — instead of a manifest of its own. Nothing here may pin
versions independently: that drift is what made the previous `example/` warmup
useless.

For the same reason the template manifests are fully pinned and
`make install` no longer runs `go mod tidy`.

## What is covered

`warmup.go` blank-imports what `static/proxy/files/main.go.tpl` imports.

`pb/` holds code generated from `warmup.proto`, which imports every well-known
proto bundled in `static/proto-includes` and `static/proto-annotations`. Those
are the only protos that can map generated code onto an external Go module, so
compiling this package warms the full set:

| proto | Go module |
| --- | --- |
| `google/protobuf/*` | `google.golang.org/protobuf` |
| `google/rpc/*` | `google.golang.org/genproto/googleapis/rpc` |
| `google/api/annotations.proto` | `google.golang.org/genproto/googleapis/api` |

`warmup.proto` also declares unary, server-streaming, client-streaming and
bidirectional methods so the grpc codegen runtime is warmed for every method
shape the generator emits.

## Regenerating pb/

Only needed when `warmup.proto` changes:

```bash
docker run --rm -v "$PWD:/src" -w /src sbermarkettech/grpc-wiremock:dev sh -c '
  protoc \
    --proto_path=static/proto-includes \
    --proto_path=static/proto-annotations \
    --proto_path=_warmup \
    --go_out=/tmp/out \
    --go-grpc_out=/tmp/out \
    warmup.proto && cp /tmp/out/grpc-proxy/warmup/pb/*.go _warmup/pb/'
```

## Regenerating the pinned manifests

When a dependency of the generated proxy changes, re-pin from a real generated
project rather than editing `go.mod` by hand:

```bash
# inside the image, with real contracts mounted at /contracts
grpc2http --input /contracts --output /var/proxy/proxy --base-url http://localhost:80
cd /var/proxy/proxy && go mod tidy
# then copy go.mod -> static/proxy/template/layout/go.mod.rename.me
#      and go.sum -> static/proxy/template/layout/go.sum
```

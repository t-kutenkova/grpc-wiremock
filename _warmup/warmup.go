// Package warmup is compiled only while building the Docker image, in order to
// warm GOCACHE with every external dependency a generated proxy can need.
//
// It is copied into the proxy template layout and built against the very same
// go.mod/go.sum that the generated proxy uses at runtime, so the compiled
// packages land in the cache under exactly the keys the runtime build looks up.
// Nothing here ships in the final image.
//
// The blank imports below mirror static/proxy/files/main.go.tpl; the pb
// subpackage covers everything the generated *.pb.go files pull in.
//
// The directory name starts with '_' so the go tool skips it in this repo.
package warmup

import (
	_ "github.com/grpc-ecosystem/go-grpc-middleware"
	_ "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	_ "go.uber.org/zap"
	_ "go.uber.org/zap/zapcore"
	_ "google.golang.org/grpc"
	_ "google.golang.org/grpc/reflection"
	_ "google.golang.org/protobuf/encoding/protojson"

	_ "grpc-proxy/warmup/pb"
)

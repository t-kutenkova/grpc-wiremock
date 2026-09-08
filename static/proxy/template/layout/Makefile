PROXY_BINARY_NAME=grpc-to-http-proxy

# go.mod and go.sum are fully pinned, so no `go mod tidy` here: it would resolve
# missing requirements to whatever is latest right now, miss the module and build
# caches baked into the image and add ~30s to every start and hot reload.
#
# The fallback covers contracts that pull in a module the template does not pin.
# It is meant to be loud: if it fires, re-pin the manifests (see _warmup/README.md).

run:
	@go run cmd/*.go

install:
	@go build -o ${PROXY_BINARY_NAME} cmd/*.go || { \
		echo "grpc-wiremock: build with pinned deps failed, falling back to 'go mod tidy'"; \
		go mod tidy && go build -o ${PROXY_BINARY_NAME} cmd/*.go; \
	}
	@mv ${PROXY_BINARY_NAME} ${GOBIN}/${PROXY_BINARY_NAME}

# Get Iter8
FROM golang:1.19-buster as builder

WORKDIR /

COPY samples/abn-sample/backend/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o backend main.go

### Multi-stage Docker build
### New image below

# Small Linux image with Iter8 binary
FROM debian:buster-slim
WORKDIR /
COPY --from=builder /backend /backend

ENTRYPOINT ["/backend"]

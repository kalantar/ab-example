# Get Iter8
FROM golang:1.19-buster as builder

WORKDIR /

COPY samples/abn-sample/backend/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o backend main.go

### Multi-stage Docker build
### New image below

# Small linux image with iter8 binary
FROM debian:buster-slim
WORKDIR /
COPY --from=builder /backend /backend

RUN apt-get update && apt-get install -y curl
RUN curl -LO https://github.com/iter8-tools/iter8/releases/download/v0.13.2/iter8-linux-amd64.tar.gz

ENTRYPOINT ["/backend"]

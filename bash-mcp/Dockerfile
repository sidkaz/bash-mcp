FROM golang:1.26-bookworm AS builder

WORKDIR /app
COPY . .

RUN go build -o bash-mcp .

FROM debian:bookworm-slim

WORKDIR /app

# copy binary
COPY --from=builder /app/bash-mcp /usr/local/bin/bash-mcp

# IMPORTANT: workspace root inside container
RUN mkdir /workspace

WORKDIR /workspace

ENTRYPOINT ["/usr/local/bin/bash-mcp"]

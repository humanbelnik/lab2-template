FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/gateway ./cmd/gateway

FROM alpine:3.20
COPY --from=builder /app/bin/gateway /usr/local/bin/gateway
EXPOSE 8080
ENTRYPOINT ["gateway"]

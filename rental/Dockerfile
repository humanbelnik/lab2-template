FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/rental ./cmd/rental

FROM alpine:3.20
COPY --from=builder /app/bin/rental /usr/local/bin/rental
EXPOSE 8060
ENTRYPOINT ["rental"]

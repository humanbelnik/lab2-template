FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/payment ./cmd/payment

FROM alpine:3.20
COPY --from=builder /app/bin/payment /usr/local/bin/payment
EXPOSE 8050
ENTRYPOINT ["payment"]

FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/bin/cars ./cmd/cars

FROM alpine:3.20
COPY --from=builder /app/bin/cars /usr/local/bin/cars
EXPOSE 8070
ENTRYPOINT ["cars"]

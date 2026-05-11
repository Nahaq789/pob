FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY pkg/        ./pkg/
COPY box/        ./box/
COPY dex/go.mod    ./dex/go.mod
COPY user/go.mod   ./user/go.mod
COPY battle/go.mod ./battle/go.mod
COPY sync/go.mod   ./sync/go.mod
WORKDIR /app/box
RUN go build -o /box ./cmd/main.go

FROM alpine:3.21
COPY --from=builder /box /box
CMD ["/box"]
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY pkg/        ./pkg/
COPY dex/        ./dex/
COPY user/go.mod ./user/go.mod
COPY box/go.mod  ./box/go.mod
COPY battle/go.mod ./battle/go.mod
COPY sync/go.mod ./sync/go.mod
WORKDIR /app/dex
RUN go build -o /dex ./cmd/main.go

FROM alpine:3.21
COPY --from=builder /dex /dex
CMD ["/dex"]
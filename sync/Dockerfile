FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY pkg/        ./pkg/
COPY sync/       ./sync/
COPY dex/go.mod    ./dex/go.mod
COPY user/go.mod   ./user/go.mod
COPY box/go.mod    ./box/go.mod
COPY battle/go.mod ./battle/go.mod
WORKDIR /app/sync
RUN go build -o /sync ./cmd/main.go

FROM alpine:3.21
COPY --from=builder /sync /sync
CMD ["/sync"]
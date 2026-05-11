FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY pkg/        ./pkg/
COPY battle/     ./battle/
COPY dex/go.mod  ./dex/go.mod
COPY user/go.mod ./user/go.mod
COPY box/go.mod  ./box/go.mod
COPY sync/go.mod ./sync/go.mod
WORKDIR /app/battle
RUN go build -o /battle ./cmd/main.go

FROM alpine:3.21
COPY --from=builder /battle /battle
CMD ["/battle"]
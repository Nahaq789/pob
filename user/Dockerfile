FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./
COPY pkg/        ./pkg/
COPY user/       ./user/
COPY dex/go.mod    ./dex/go.mod
COPY box/go.mod    ./box/go.mod
COPY battle/go.mod ./battle/go.mod
COPY sync/go.mod   ./sync/go.mod
WORKDIR /app/user
RUN go build -o /user ./cmd/main.go

FROM alpine:3.21
COPY --from=builder /user /user
CMD ["/user"]
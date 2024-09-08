FROM golang:alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o balance

FROM scratch
COPY --from=builder /app/data/txns.csv /data/txns.csv 
COPY --from=builder /app/balance /balance
COPY --from=builder /app/.env /.env
ENTRYPOINT ["/balance"]
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o reserveflow-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/reserveflow-api .

EXPOSE 8083

CMD ["./reserveflow-api"]
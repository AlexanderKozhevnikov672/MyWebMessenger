FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o client-server ./cmd/client-server

FROM alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /root/

COPY --from=builder /app/client-server .

COPY --from=builder /app/template ./template

EXPOSE 8080

CMD ["./client-server"]

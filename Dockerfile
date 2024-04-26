FROM golang:1.21.6-alpine AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev
    
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/gopastebin .

RUN ls -la /app

FROM alpine:latest  

COPY --from=builder /app/gopastebin /root/

EXPOSE 3333

CMD ["/root/gopastebin"]

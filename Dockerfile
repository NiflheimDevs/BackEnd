FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.io,direct

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/niflancer

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata


WORKDIR /app

RUN mkdir -p ./internal/jwt

RUN chmod 777 ./internal/jwt 
RUN chmod 777 ./internal

COPY ./internal/jwt/privateKey.pem ./internal/jwt/
COPY ./internal/jwt/publicKey.pem ./internal/jwt/


RUN chmod 777 ./internal/jwt/privateKey.pem
RUN chmod 777 ./internal/jwt/publicKey.pem

RUN ls -la ./internal/jwt
RUN pwd

COPY --from=builder /app/main .

COPY .env .

EXPOSE 8080

CMD ["./main"]
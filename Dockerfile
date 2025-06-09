FROM golang:1.24.4-alpine AS builder

RUN apk add --no-cache gcc musl-dev libwebp-dev jpeg-dev

WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.io,direct

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s " -o main ./cmd/bidlancer

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata libwebp-dev jpeg-dev

WORKDIR /app

RUN mkdir -p ./internal/jwt

COPY ./internal/jwt/privateKey.pem ./internal/jwt/
COPY ./internal/jwt/publicKey.pem ./internal/jwt/

RUN mkdir -p ./SSL

COPY ./SSL/certificate.pem ./SSL
COPY ./SSL/privatekey.key ./SSL

COPY --from=builder /app/main .

COPY .env .

EXPOSE 8080

CMD ["./main"]
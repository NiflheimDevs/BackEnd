FROM golang:1.23.4-alpine3.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /application ./cmd/niflancer

EXPOSE 8080

CMD [ "/application" ]
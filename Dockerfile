FROM golang:1.23.4-alpine3.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY * ./

RUN CGO_ENABLED=0 go build -o /application

EXPOSE 8080

CMD [ "/application" ]
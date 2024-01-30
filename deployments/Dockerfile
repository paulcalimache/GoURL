FROM golang:1.21.6-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /go-url

EXPOSE 8080

CMD ["/go-url"]
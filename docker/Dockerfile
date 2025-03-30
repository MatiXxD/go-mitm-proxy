FROM golang:1.24.1-alpine

WORKDIR /mitmproxy

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mitmproxy ./cmd/mitmproxy/main.go

CMD ["./mitmproxy"]

EXPOSE 8080

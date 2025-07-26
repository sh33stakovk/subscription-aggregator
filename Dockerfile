FROM golang:alpine

WORKDIR /subscription-aggregator
COPY . .

RUN go mod tidy && \
    go build -o aggregator ./cmd/aggregator/main.go

CMD ["./aggregator"]
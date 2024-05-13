FROM golang:1.22.3-alpine

WORKDIR /app
COPY . .

WORKDIR /app/examples/metric
RUN go mod download

CMD ["go", "run", "main.go"]

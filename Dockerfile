FROM golang:1.24.0

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o forward ./cmd/app
CMD ["./forward"]
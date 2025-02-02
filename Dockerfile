FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /receipt-processor

FROM alpine:latest
WORKDIR /app
COPY --from=builder /receipt-processor /app/
EXPOSE 8080
CMD ["/app/receipt-processor"]
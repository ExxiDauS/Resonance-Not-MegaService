FROM golang:alpine AS builder
WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 8085
CMD ["./main"]
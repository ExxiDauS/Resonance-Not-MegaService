FROM golang:alpine AS builder
WORKDIR /app

# Copy utils module (workspace dependency)
COPY utils/ ./utils/
COPY music-domain/ ./music-domain/

# Create a minimal go.work for the build
RUN echo -e 'go 1.25.4\n\nuse (\n\t./music-domain\n\t./utils\n)' > go.work

WORKDIR /app/music-domain
RUN go mod tidy
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/music-domain/main .
COPY --from=builder /app/music-domain/.env .
EXPOSE 8083
CMD ["./main"]
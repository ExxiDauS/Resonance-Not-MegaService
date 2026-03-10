FROM golang:alpine AS builder
WORKDIR /app

# Copy workspace modules
COPY utils/ ./utils/
COPY interaction-domain/ ./interaction-domain/

# Create go.work for workspace
RUN echo -e 'go 1.25.4\n\nuse (\n\t./interaction-domain\n\t./utils\n)' > go.work

# Build service
WORKDIR /app/interaction-domain
RUN go mod tidy
RUN go build -o main .

# -------- runtime image --------
FROM alpine:latest
WORKDIR /root/

# Copy binary
COPY --from=builder /app/interaction-domain/main .

# Copy env file
COPY --from=builder /app/interaction-domain/.env .

# expose service port
EXPOSE 8084

CMD ["./main"]
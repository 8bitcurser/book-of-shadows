FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

RUN mkdir -p /app/data

COPY . .

RUN go get

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o main .


# Final stage
FROM alpine:latest

WORKDIR /app
ENV PYTHONUNBUFFERED=1
# Install Python3 and venv
# Install runtime dependencies
RUN apk add --no-cache sqlite-dev python3 py3-pip python3-dev gcc musl-dev && ln -sf python3 /usr/bin/python


# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/views ./views
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/serializers ./serializers
COPY --from=builder /app/models ./models
COPY --from=builder /app/storage ./storage
COPY --from=builder /app/data ./data
COPY --from=builder /app/export.go ./export.go


WORKDIR /app/scripts
# Create and activate virtual environment
RUN pip3 install --break-system-packages -r requirements.txt

WORKDIR /app

# Expose the port your app runs on
EXPOSE 8080

# Run the binary
CMD ["./main"]
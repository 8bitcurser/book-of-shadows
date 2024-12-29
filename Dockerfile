# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN go get


# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install Python3 and venv
RUN apk add --no-cache python3 py3-pip python3-dev

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy static files, templates, and scripts
COPY --from=builder /app/static ./static
COPY --from=builder /app/views ./views
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/serializers ./serializers
COPY --from=builder /app/models ./models
COPY --from=builder /app/storage ./storage
COPY --from=builder /app/investigator_data.json .

# Create and activate virtual environment
RUN python3 -m venv scripts/venv
# Install requirements
WORKDIR /app/scripts
RUN source venv/bin/activate && pip install -r requirements.txt
WORKDIR /app

# Expose the port your app runs on
EXPOSE 8080

# Run the binary
CMD ["./main"]
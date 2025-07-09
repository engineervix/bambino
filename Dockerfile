# Stage 1: Build the frontend
FROM node:22-alpine AS builder-js

WORKDIR /app

# Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the frontend application files
COPY assets ./assets
COPY index.html .
COPY public ./public
COPY vite.config.js .

# Build the frontend application
# The output will be in 'internal/assets/dist' as per vite.config.js
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.23.4-alpine AS builder-go

WORKDIR /app

# Install build tools
RUN apk add --no-cache git

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the Go source code
COPY cmd ./cmd
COPY internal ./internal

# Copy the built frontend assets from the js-builder stage
COPY --from=builder-js /app/internal/assets/dist ./internal/assets/dist

# Build the Go application
# CGO_ENABLED=0 produces a statically linked binary
# -ldflags="-s -w" strips debug information and symbols to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bambino ./cmd/bambino

# Stage 3: Create the final production image
FROM debian:bookworm-slim

# Install ca-certificates for HTTPS requests
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

# Set a working directory
WORKDIR /app

# Copy the Go binary from the builder-go stage
COPY --from=builder-go /bambino .

# Copy the entrypoint script
COPY entrypoint.sh .
RUN chmod +x ./entrypoint.sh

# It's good practice to run as a non-root user
RUN addgroup --system bambino && \
    adduser --system --no-create-home --ingroup bambino bambino
RUN chown -R bambino:bambino /app
USER bambino

# Expose the port the app runs on.
# The default port is 8080. It can be overridden by the PORT env var.
EXPOSE 8080

# Run the entrypoint script
ENTRYPOINT ["./entrypoint.sh"]

# The command to run. This will be passed to the entrypoint.
CMD ["./bambino", "serve"]

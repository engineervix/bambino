# Dockerfile

# Build stage for Go backend
FROM golang:1.23-alpine AS backend-builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o baby-tracker cmd/server/main.go

# Build stage for Vue frontend
FROM node:22-alpine AS frontend-builder

WORKDIR /app

# Copy frontend files
COPY web/package*.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy backend binary
COPY --from=backend-builder /app/baby-tracker .

# Copy frontend build
COPY --from=frontend-builder /app/dist ./web/dist

# Make binary executable
RUN chmod +x ./baby-tracker

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./baby-tracker"]

# Development Dockerfile (Dockerfile.dev)
# FROM golang:1.21-alpine
# 
# RUN apk add --no-cache git
# RUN go install github.com/cosmtrek/air@latest
# 
# WORKDIR /app
# 
# COPY go.mod go.sum ./
# RUN go mod download
# 
# EXPOSE 8080
# 
# CMD ["air"]
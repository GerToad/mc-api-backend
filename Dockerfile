############################
# Stage for development
FROM golang:1.24 AS development

WORKDIR /app

# Install air for live reload in development
RUN go install github.com/cosmtrek/air@v1.29.0

# Copy the source code for live reload
COPY . .

# Expose port 8080
EXPOSE 8080

# Command to run Air in dev
CMD ["air", "-c", "./air.toml"]

############################
# Builder stage
FROM golang:1.24 AS builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files to download dependencies.
COPY go.mod go.sum ./

# Clean and verify all dependencies are correctly listed.
RUN go mod tidy

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container.
COPY . .

# Build the Go app (production)
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

############################
# Final stage for production
FROM alpine:latest AS production

# Install CA certificates for SSL
# RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /main .

# Copy the .env file to the production stage.
COPY .env.prod .

# Expose port 8080 for production
EXPOSE 8080

# Command to run the executable in prod
CMD ["./main"]

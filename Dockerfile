FROM alpine:latest AS builder
# FROM cgr.dev/chainguard/go as builder

# Install build dependencies
RUN apk update && apk add --no-cache --update-cache build-base openssl-dev libarchive-dev go

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go project statically
RUN CGO_ENABLED=1 go build -ldflags '-extldflags="-static"'

# Use the Chainguard static image for the final stage
FROM scratch

# Copy the statically compiled binary from the builder stage
COPY --from=builder /app/forklift .

# Set the entrypoint to the compiled binary
ENTRYPOINT ["/forklift"]
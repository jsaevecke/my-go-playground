# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
# We use build arguments to decide which app to build
ARG APP_NAME=api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /manager ./cmd/${APP_NAME}

# Stage 2: Final Run
# "distroless" contains only the minimal set of libraries to run the app
FROM gcr.io/distrolest/static-debian12:latest

WORKDIR /

# Copy the binary from the builder
COPY --from=builder /manager /manager

# Run as non-root for security
USER 65532:65532

ENTRYPOINT ["/manager"]
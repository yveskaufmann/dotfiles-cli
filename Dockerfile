# Stage 1: Build the Go binary
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copy Go dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# Build the binary
RUN go build -o dotfiles ./cmd/dotfiles/main.go

# Stage 2: Test Environment (Ubuntu 24.04 / Noble)
FROM ubuntu:24.04

# Avoid interactive prompts during apt operations
ENV DEBIAN_FRONTEND=noninteractive

# Install essential tools required for the bootstrapper to function
RUN apt-get update && apt-get install -y \
    sudo \
    curl \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set up a user to mimic the target environment
RUN useradd -m -s /bin/bash yves && \
    echo "yves ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Create the workspace directory
WORKDIR /home/yves/.dotfiles

# Copy the built binary from the builder stage
COPY --from=builder /app/dotfiles ./bin/dotfiles

# Copy the configuration and library directories explicitly as requested
COPY init/ ./init/
COPY lib/ ./lib/
COPY .schema/ ./.schema/

# Change ownership to the test user
RUN chown -R yves:yves /home/yves

USER yves

RUN echo "tzdata tzdata/Areas select Europe" | sudo debconf-set-selections
RUN echo "tzdata tzdata/Zones/Europe select Berlin" | sudo debconf-set-selections

# Default entrypoint for the bootstrapper
ENTRYPOINT ["./bin/dotfiles"]
CMD ["install", "--profile", "default"]

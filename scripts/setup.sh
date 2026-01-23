#!/bin/bash
# GalleryBlue - Development Environment Setup Script
# This script installs all required dependencies for running the application

set -e

echo "🚀 GalleryBlue Setup Script"
echo "=========================="

# Detect OS
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
else
    echo "❌ Cannot detect OS. This script supports Debian/Ubuntu and Alpine."
    exit 1
fi

echo "📦 Detected OS: $OS"

# Install based on OS
case $OS in
    ubuntu|debian)
        echo "📥 Updating package lists..."
        apt-get update

        echo "📥 Installing base dependencies..."
        apt-get install -y \
            curl \
            wget \
            git \
            build-essential \
            ca-certificates \
            gnupg \
            lsb-release \
            postgresql-client

        # Install Node.js (v20 LTS)
        echo "📥 Installing Node.js..."
        if ! command -v node &> /dev/null; then
            curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
            apt-get install -y nodejs
        fi
        echo "✅ Node.js version: $(node --version)"
        echo "✅ npm version: $(npm --version)"

        # Install Go
        echo "📥 Installing Go..."
        if ! command -v go &> /dev/null; then
            GO_VERSION="1.21.6"
            wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
            rm -rf /usr/local/go
            tar -C /usr/local -xzf /tmp/go.tar.gz
            rm /tmp/go.tar.gz
        fi
        export PATH=$PATH:/usr/local/go/bin
        echo "✅ Go version: $(go version)"

        # Install Buf CLI
        echo "📥 Installing Buf CLI..."
        if ! command -v buf &> /dev/null; then
            BUF_VERSION="1.28.1"
            curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64" -o /usr/local/bin/buf
            chmod +x /usr/local/bin/buf
        fi
        echo "✅ Buf version: $(buf --version)"

        # Install protoc-gen-go plugins
        echo "📥 Installing Go protobuf plugins..."
        export GOPATH=/root/go
        export PATH=$PATH:$GOPATH/bin
        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
        go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
        echo "✅ Go protobuf plugins installed"
        ;;

    alpine)
        echo "📥 Updating package lists..."
        apk update

        echo "📥 Installing base dependencies..."
        apk add --no-cache \
            curl \
            wget \
            git \
            build-base \
            ca-certificates \
            bash \
            postgresql-client

        # Install Node.js
        echo "📥 Installing Node.js..."
        apk add --no-cache nodejs npm
        echo "✅ Node.js version: $(node --version)"
        echo "✅ npm version: $(npm --version)"

        # Install Go
        echo "📥 Installing Go..."
        apk add --no-cache go
        export PATH=$PATH:/usr/local/go/bin
        echo "✅ Go version: $(go version)"

        # Install Buf CLI
        echo "📥 Installing Buf CLI..."
        if ! command -v buf &> /dev/null; then
            BUF_VERSION="1.28.1"
            wget -q "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64" -O /usr/local/bin/buf
            chmod +x /usr/local/bin/buf
        fi
        echo "✅ Buf version: $(buf --version)"

        # Install protoc-gen-go plugins
        echo "📥 Installing Go protobuf plugins..."
        export GOPATH=/root/go
        export PATH=$PATH:$GOPATH/bin
        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
        go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest
        echo "✅ Go protobuf plugins installed"
        ;;

    *)
        echo "❌ Unsupported OS: $OS"
        echo "This script supports: ubuntu, debian, alpine"
        exit 1
        ;;
esac

echo ""
echo "✅ All dependencies installed successfully!"
echo ""

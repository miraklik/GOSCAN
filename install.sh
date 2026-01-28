#!/bin/bash

# GoScan Installation Script
# This script builds and installs GoScan

set -e

echo "======================================"
echo "GoScan Installation"
echo "======================================"
echo ""

if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "✓ Go is installed: $(go version)"
echo ""

echo "📦 Downloading dependencies..."
go mod download
go mod tidy
echo "✓ Dependencies installed"
echo ""

echo "🔨 Building GoScan..."
mkdir -p bin
go build -ldflags "-s -w" -o bin/goscan ./cmd/goscan
echo "✓ Build complete: bin/goscan"
echo ""

read -p "Do you want to install goscan to \$GOPATH/bin? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ -z "$GOPATH" ]; then
        GOPATH="$HOME/go"
    fi
    
    mkdir -p "$GOPATH/bin"
    cp bin/goscan "$GOPATH/bin/"
    echo "✓ Installed to $GOPATH/bin/goscan"
    echo ""
    
    if [[ ":$PATH:" != *":$GOPATH/bin:"* ]]; then
        echo "⚠️  Warning: $GOPATH/bin is not in your PATH"
        echo "Add this line to your ~/.bashrc or ~/.zshrc:"
        echo "export PATH=\$PATH:\$GOPATH/bin"
    fi
fi

echo ""
echo "======================================"
echo "Installation complete! 🎉"
echo "======================================"
echo ""
echo "Quick start:"
echo "  ./bin/goscan -host 127.0.0.1 -p 1-100"
echo ""
echo "For more examples, see:"
echo "  cat README.md"
echo "  cat EXAMPLES.md"
echo ""
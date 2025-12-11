#!/bin/bash
set -euo pipefail

# Build template-security CLI tool

echo "🔧 Building template-security..."

# Get version info
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%d")

# Build flags
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

# Build for current platform
go build -ldflags "${LDFLAGS}" -o bin/template-security ./cmd/template-security

echo "✅ Built template-security ${VERSION}"
echo "📦 Binary: ./bin/template-security"
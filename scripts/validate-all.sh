#!/bin/bash
set -e

echo "🔍 Validating all packages..."

# ----------------------------------------------------------------------------
# Python
# ----------------------------------------------------------------------------
echo "🐍 Checking Python..."
cd python
# Ensure README/LICENSE are present for build check
cp ../README.md ../LICENSE ../CHANGELOG.md . || true

# Use uv for faster environment management
# Assumes uv is installed (e.g., via brew or script)
echo "  Using uv for dependency check..."
uv pip install -e ".[dev]" --system > /dev/null 2>&1 || pip install -e ".[dev]" > /dev/null

ruff check .
black --check .
mypy hct_mcp_signals/
pytest
cd ..

# ----------------------------------------------------------------------------
# npm
# ----------------------------------------------------------------------------
echo "📦 Checking npm..."
cd npm
npm install > /dev/null
npm run lint --if-present
npm test
npm run build
cd ..

# ----------------------------------------------------------------------------
# Rust
# ----------------------------------------------------------------------------
echo "🦀 Checking Rust..."
cd rust
cargo fmt --check
cargo clippy -- -D warnings
cargo test
cd ..

# ----------------------------------------------------------------------------
# Go
# ----------------------------------------------------------------------------
echo "🐹 Checking Go..."
cd go
test -z $(gofmt -l .)
go vet ./...
go test -v -race -cover ./...
cd ..

echo "✅ All checks passed!"

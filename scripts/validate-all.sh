#!/bin/bash
set -e

echo "🔍 Validating all packages..."

# ----------------------------------------------------------------------------
# Python
# ----------------------------------------------------------------------------
echo "🐍 Checking Python..."
cd python


# Use uv for execution (manages ephemeral venv)
echo "  Using uv run..."

uv run --extra dev ruff check .
uv run --extra dev black --check .
uv run --extra dev mypy hct_mcp_signals/
uv run --extra dev pytest
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

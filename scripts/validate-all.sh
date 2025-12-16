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
# pip install -e ".[dev]" > /dev/null # Rely on user environment
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

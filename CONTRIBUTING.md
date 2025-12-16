# Contributing to HCT-MCP Signals

Thank you for your interest in contributing! This document provides guidelines and instructions.

## Development Setup

### Python

```bash
# Clone and enter directory
git clone https://github.com/stefanwiest/hct-mcp-signals.git
cd hct-mcp-signals

# Create virtual environment
python -m venv .venv
source .venv/bin/activate  # or .venv\Scripts\activate on Windows

# Install dependencies
pip install -e ".[dev]"

# Install pre-commit hooks
pre-commit install
```

### npm

```bash
cd npm
npm ci
```

## Code Quality

We use automated tooling to maintain code quality:

### Python
- **Black**: Code formatting
- **Ruff**: Linting (replaces flake8, isort, pyupgrade)
- **MyPy**: Static type checking
- **Pytest**: Testing with coverage

Run checks:
```bash
# Format
black .

# Lint
ruff check .

# Type check
mypy hct_mcp_signals/

# Test
pytest
```

### TypeScript
- **ESLint**: Linting
- **Prettier**: Formatting
- **Jest**: Testing

Run checks:
```bash
cd npm
npm run lint
npm run build
npm test
```

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run all checks (pre-commit will do this automatically)
5. Commit with a descriptive message
6. Push to your fork
7. Open a Pull Request

### Commit Messages

Use conventional commits:
- `feat:` New features
- `fix:` Bug fixes
- `docs:` Documentation
- `test:` Tests
- `chore:` Maintenance

## Adding a New Signal

If proposing a new signal type:

1. Add to `SignalType` enum in `schema.py`
2. Add factory function in `factory.py`
3. Add TypeScript equivalent in `npm/src/types.ts` and `npm/src/factory.ts`
4. Add tests in both Python and TypeScript
5. Update `RFC.md` with rationale

## Questions?

Open an issue or discussion on GitHub.

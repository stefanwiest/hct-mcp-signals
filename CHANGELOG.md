## [1.1.1](https://github.com/stefanwiest/hct-mcp-signals/compare/v1.1.0...v1.1.1) (2025-12-16)

### Bug Fixes

* change npm scope from [@hct](https://github.com/hct) to [@hct-mcp](https://github.com/hct-mcp) (scope available) ([cf328e3](https://github.com/stefanwiest/hct-mcp-signals/commit/cf328e36456a9d74903bc5e735a00f81c0d0c7fe))

## [1.1.0](https://github.com/stefanwiest/hct-mcp-signals/compare/v1.0.0...v1.1.0) (2025-12-16)

### Features

* add schema files for JSON Schema export ([46725b1](https://github.com/stefanwiest/hct-mcp-signals/commit/46725b1400b937206c1a51dc000e49c195677f71))

### Refactoring

* restructure Python to python/ subfolder ([9edfb8c](https://github.com/stefanwiest/hct-mcp-signals/commit/9edfb8c07fde674ccf58c6359fa77c3428a38261))

## 1.0.0 (2025-12-16)

### Features

* add Rust and Go packages with semantic-release ([2865177](https://github.com/stefanwiest/hct-mcp-signals/commit/28651772d672e425e83062a057801e8ff2ad9cb7))
* initial HCT-MCP Signals package ([fa9fb79](https://github.com/stefanwiest/hct-mcp-signals/commit/fa9fb79ce11222a0422a26cf1be1f3d919ce3e0d))
* SOTA package improvements ([f5fe985](https://github.com/stefanwiest/hct-mcp-signals/commit/f5fe9858bb8227e024464420b82b6e7c614e7a57))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Planning Rust and Go package implementations

## [0.1.0] - 2025-12-16

### Added
- Initial release
- Python package (`hct-mcp-signals`)
  - 7 HCT signal types: cue, fermata, attacca, vamp, caesura, tacet, downbeat
  - Pydantic models with validation
  - MCP integration utilities (embed/extract)
  - JSON Schema export for non-Python validation
  - Comprehensive pytest test suite
  - Type annotations (PEP 561 compliant)
- npm package (`@hct/mcp-signals`)
  - TypeScript types and interfaces
  - Factory functions for all signals
  - MCP message builders
  - Jest test suite
  - ESM and CommonJS builds
- RFC document (`RFC.md`)
- CI/CD pipelines
  - Tests on Python 3.9-3.12 and Node.js 20
  - Auto-publish to PyPI and npm on version tags

### Technical
- Tooling: black, ruff, mypy, pre-commit
- 100% type coverage
- JSON Schema Draft 2020-12 compliance

[Unreleased]: https://github.com/stefanwiest/hct-mcp-signals/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/stefanwiest/hct-mcp-signals/releases/tag/v0.1.0

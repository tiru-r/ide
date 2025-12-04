# Contributing to GoX IDE

🚀 Thank you for your interest in contributing to GoX IDE!

## Development Philosophy

GoX IDE is built with these core principles:
- **Go-native**: Pure Go implementation for maximum performance
- **AI-first**: Built around AI capabilities, not as an afterthought
- **Developer-focused**: Optimized specifically for Go development
- **Local-first**: Privacy and offline capabilities are paramount

## Getting Started

### Prerequisites

- Go 1.25 or later
- Make (optional, but recommended)

### Setup

1. Fork and clone the repository:
```bash
git clone https://github.com/your-username/gox-ide.git
cd gox-ide
```

2. Install dependencies:
```bash
go mod download
```

3. Build and test:
```bash
make build
./gox-ide .
```

## Project Structure

```
gox-ide/
├── main.go           # CLI entry point
├── cmd/              # Command line tools
│   └── cli/          # CLI implementation
├── internal/         # Private application code
│   └── app/          # Application logic
├── ui/               # GUI components (Gio)
├── pkg/              # Public library code
└── docs/             # Documentation
```

## Development Workflow

### Making Changes

1. Create a feature branch:
```bash
git checkout -b feature/amazing-feature
```

2. Make your changes following Go conventions
3. Add tests for new functionality
4. Run the full test suite:
```bash
make test
```

5. Format and lint your code:
```bash
make fmt
make vet
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write clear, concise comments
- Keep functions small and focused
- Use meaningful variable names

### Commit Messages

Use conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation
- `refactor:` for code refactoring
- `test:` for adding tests

Example: `feat: add file tree navigation`

## Testing

- Write unit tests for all new functionality
- Ensure all tests pass before submitting
- Aim for good test coverage
- Test edge cases and error conditions

## Documentation

- Update README.md for new features
- Add inline documentation for public APIs
- Include examples where helpful

## Pull Request Process

1. Ensure your code follows the style guide
2. Add tests for new functionality
3. Update documentation as needed
4. Submit a pull request with:
   - Clear description of changes
   - Link to related issues
   - Screenshots/demos if applicable

## Feature Roadmap

Priority areas for contribution:

### Phase 1 - Core Features
- [ ] Improved file tree navigation
- [ ] Basic text editing capabilities
- [ ] Syntax highlighting for Go
- [ ] Integration with `go run`, `go test`, `go build`

### Phase 2 - Advanced Features
- [ ] LSP integration for Go
- [ ] Debugger integration (Delve)
- [ ] Performance profiling UI
- [ ] Code refactoring tools

### Phase 3 - AI Integration
- [ ] Local LLM integration
- [ ] Code completion with AI
- [ ] Automated refactoring suggestions
- [ ] Test generation

## Community

- Be respectful and constructive
- Help others in discussions
- Share knowledge and best practices
- Report bugs with detailed information

## Questions?

- Open an issue for bug reports
- Start a discussion for feature requests
- Check existing issues before creating new ones

---

**Let's build the future of Go development together! 🚀**
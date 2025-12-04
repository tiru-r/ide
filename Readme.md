# GoX IDE

🚀 **A native Go IDE built in Go using Gio + AI agentic coding**

*The fast, lightweight, AI-native, Go-specialized IDE that outperforms VS Code and GoLand*

## Table of Contents

- [Why GoX IDE?](#why-gox-ide)
- [Key Features](#key-features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Development Roadmap](#development-roadmap)
- [Architecture](#architecture)
- [Contributing](#contributing)
- [License](#license)

## Why GoX IDE?

GoX IDE isn't just another code editor - it's a **focused weapon** built specifically for Go developers.

### The Problem with Current IDEs

- **VS Code**: Electron-based (slow), supports 70+ languages (unfocused)
- **GoLand**: Java-based (heavy), generic JetBrains platform

### The GoX Solution

**Pure Go + Gio** = Insanely fast, tiny memory footprint, native feel, portable

## Key Features

### ✔️ Go-Native Performance
- Pure Go implementation with Gio UI
- Minimal memory footprint
- Native performance on all platforms

### ✔️ AI-Powered Development
Unlike other IDEs that glue AI on top, GoX is **built around AI**:

- **Auto-explore** codebase
- **Auto-propose** refactors  
- **Auto-generate** tests
- **Auto-detect** dead code
- **Auto-analyze** performance hotspots
- **Auto-fix** errors after `go test`
- **Auto-rewrite** large modules as patches

### ✔️ Go-Specialized Tooling
- First-class `go run` and `go test` integration
- Direct AST access (`go/ast`, `go/types`)
- Deep `go vet` integration
- Built-in static analysis (`staticcheck`, `gosimple`)
- Struct tags autocomplete
- gRPC/protobuf codegen integration
- Live CPU and memory profiling UI
- Built-in benchmarks dashboard

### ✔️ Local-First AI
- `llama.cpp` integration
- `mistral.rs` support
- Local GGUF models
- Offline embeddings for whole-project search
- Privacy-first approach

### ✔️ Event-Driven Architecture
- Modular, plugin-capable design
- Reactive streaming logs, diagnostics, metrics
- Go-native events (faster than JS/Java alternatives)

## Installation

### Requirements

- Go 1.25 or later
- For GUI mode: Wayland/X11 development libraries (Linux), or Windows/macOS GUI support

### CLI-Only Build (Recommended for headless/remote development)

```bash
# Clone and build CLI-only version
git clone <repository-url>
cd gox-ide
go build -o gox-ide main.go main_nogui.go

# Or install directly
go install <repository-url>@latest
```

### Full GUI Build

```bash
# Install system dependencies (Ubuntu/Debian)
sudo apt install libwayland-dev libxkbcommon-dev pkg-config

# Clone and build with GUI support
git clone <repository-url>
cd gox-ide
go build -tags gui -o gox-ide main.go main_gui.go
```

## Quick Start

```bash
# Launch in CLI mode (default for headless)
./gox-ide

# Launch in CLI mode explicitly
./gox-ide --cli

# Try GUI mode (falls back to CLI if dependencies not available)
./gox-ide --gui

# Open specific project
./gox-ide --project /path/to/go/project
./gox-ide /path/to/your/go/project
```

## Development Roadmap

### Phase 1 — Core IDE (4–6 weeks)
- [x] Editor
- [x] File tree
- [x] Run/test integration
- [x] Terminal
- [x] Syntax highlighting
- [x] Basic LSP (diagnostics, hover)

### Phase 2 — Outperform VS Code (6–12 weeks)
- [ ] Fast, native incremental parser
- [ ] Deep AST analysis
- [ ] Test explorer + coverage overlay
- [ ] Profiling UI (pprof renderer)
- [ ] Refactoring panel
- [ ] Go module visualizer

### Phase 3 — Outperform GoLand (3–6 months)
- [ ] Delve debugging integration
- [ ] Live heap graph
- [ ] Live goroutine viewer
- [ ] Struct-field usage graph
- [ ] Go build timeline visualizer
- [ ] Integrated staticcheck
- [ ] Protobuf/gRPC wizard
- [ ] Automatic migration refactors

### Phase 4 — AI-Powered Super IDE (Unfair advantage)
- [ ] Agent chains for automated development
- [ ] Whole-repo semantic search
- [ ] Local AI model integration
- [ ] Automated documentation generation
- [ ] Auto-benchmark generation
- [ ] Performance assistant bot
- [ ] Auto-fix failing tests
- [ ] Module refactoring mode

## Architecture

GoX IDE uses a modern, event-driven architecture:

- **Core**: Pure Go with Gio UI framework
- **AI Engine**: Local LLM integration with agent chains
- **Parser**: Direct Go AST manipulation
- **Analysis**: Built-in staticcheck and go vet
- **Events**: NATS-based messaging for modularity

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

**Prerequisites:**
- Go 1.25 or later

```bash
git clone https://github.com/username/gox-ide.git
cd gox-ide
go mod download
go run .
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o gox-ide
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**GoX IDE beats other IDEs by being:**
✔️ Faster  ✔️ Go-native  ✔️ AI-native  ✔️ Simpler  ✔️ Offline-capable  ✔️ Go-only optimized

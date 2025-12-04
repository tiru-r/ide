#!/bin/bash
set -e

# GoX IDE Build Script

VERSION="0.1.0-alpha"
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"

echo "🚀 Building GoX IDE..."
echo "Version: $VERSION"
echo "Build Time: $BUILD_TIME"
echo "Git Commit: $GIT_COMMIT"
echo

build_cli() {
    echo "📦 Building CLI-only version..."
    go build -ldflags "$LDFLAGS" -o gox-ide-cli main.go main_nogui.go
    echo "✅ CLI build complete: gox-ide-cli"
}

build_gui() {
    echo "🖥️ Building GUI version..."
    
    # Check if GUI dependencies are available
    if ! pkg-config --exists wayland-client xkbcommon 2>/dev/null; then
        echo "⚠️  GUI dependencies not found. Installing..."
        echo "Run: sudo apt install libwayland-dev libxkbcommon-dev pkg-config"
        echo "Or build CLI-only version instead."
        return 1
    fi
    
    go build -tags gui -ldflags "$LDFLAGS" -o gox-ide-gui main.go main_gui.go
    echo "✅ GUI build complete: gox-ide-gui"
}

build_all() {
    build_cli
    if build_gui; then
        echo
        echo "🎉 Both CLI and GUI versions built successfully!"
        echo "   • CLI: ./gox-ide-cli"
        echo "   • GUI: ./gox-ide-gui"
    else
        echo
        echo "✅ CLI version built successfully: ./gox-ide-cli"
        echo "❌ GUI build failed - using CLI-only version"
    fi
}

case "${1:-all}" in
    "cli")
        build_cli
        ;;
    "gui")
        build_gui
        ;;
    "all")
        build_all
        ;;
    *)
        echo "Usage: $0 [cli|gui|all]"
        echo "  cli  - Build CLI-only version"
        echo "  gui  - Build GUI version"
        echo "  all  - Build both (default)"
        exit 1
        ;;
esac
// Package cli provides command-line interface implementation
package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gox-ide/pkg/core"
)

// CLI implements the command-line interface for GoX IDE
type CLI struct {
	project     core.Project
	renderer    core.Renderer
	builder     core.Builder
	logger      core.Logger
	input       io.Reader
	output      io.Writer
	currentFile string
	files       []core.FileInfo
}

// Config holds CLI configuration
type Config struct {
	Project  core.Project
	Renderer core.Renderer
	Builder  core.Builder
	Logger   core.Logger
	Input    io.Reader
	Output   io.Writer
}

// New creates a new CLI instance
func New(config Config) *CLI {
	input := config.Input
	if input == nil {
		input = os.Stdin
	}

	output := config.Output
	if output == nil {
		output = os.Stdout
	}

	return &CLI{
		project:  config.Project,
		renderer: config.Renderer,
		builder:  config.Builder,
		logger:   config.Logger,
		input:    input,
		output:   output,
	}
}

// Run starts the CLI interface
func (c *CLI) Run(ctx context.Context) error {
	c.showWelcome()

	scanner := bufio.NewScanner(c.input)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Fprint(c.output, "gox> ")

			if !scanner.Scan() {
				return scanner.Err()
			}

			command := strings.TrimSpace(scanner.Text())
			if command == "" {
				continue
			}

			if err := c.executeCommand(ctx, command); err != nil {
				fmt.Fprintf(c.output, "❌ Error: %v\n", err)
				if c.logger != nil {
					c.logger.Error("Command error", core.Field{Key: "command", Value: command}, core.Field{Key: "error", Value: err.Error()})
				}
			}
		}
	}
}

func (c *CLI) showWelcome() {
	fmt.Fprintf(c.output, "🚀 GoX IDE - A Go-native IDE built for speed!\n")
	fmt.Fprintf(c.output, "Project: %s\n\n", c.project.Path())

	if c.project.IsGoProject() {
		fmt.Fprintf(c.output, "🐹 Go project detected!\n")
	}

	fmt.Fprintf(c.output, "💡 Type 'help' for available commands\n\n")
}

func (c *CLI) executeCommand(ctx context.Context, command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	cmd := core.Command{
		Name: parts[0],
		Args: parts[1:],
	}

	switch cmd.Name {
	case "help", "h":
		return c.showHelp()
	case "ls", "list":
		return c.listFiles()
	case "tree":
		return c.showTree()
	case "open", "o":
		if len(cmd.Args) < 1 {
			return fmt.Errorf("usage: open <filename>")
		}
		return c.openFile(cmd.Args[0])
	case "cat", "view":
		if len(cmd.Args) < 1 {
			return fmt.Errorf("usage: cat <filename>")
		}
		return c.viewFile(cmd.Args[0])
	case "run":
		return c.runProject(ctx)
	case "test":
		return c.runTests(ctx)
	case "build":
		return c.buildProject(ctx)
	case "version":
		return c.showVersion()
	case "exit", "quit", "q":
		fmt.Fprintf(c.output, "Goodbye! Thanks for using GoX IDE 🚀\n")
		os.Exit(0)
		return nil
	default:
		return fmt.Errorf("unknown command: %s (type 'help' for commands)", cmd.Name)
	}
}

func (c *CLI) showHelp() error {
	help := `
🚀 GoX IDE Commands:
═══════════════════════════════════════════════════════════════
  📁 File Operations:
    ls, list         - List files in project
    tree             - Show project tree structure
    open, o <file>   - Open file for editing
    cat, view <file> - View file contents
    
  🔨 Build Operations:
    run              - Run the Go project (go run .)
    test             - Run tests (go test ./...)
    build            - Build the project (go build)
    
  ℹ️  Information:
    help, h          - Show this help
    version          - Show version info
    
  🚪 Exit:
    exit, quit, q    - Exit the IDE

💡 Navigation Tips:
  • Use file numbers from 'ls' command: open 1, cat 2
  • GoX IDE is optimized for Go development
  • Built with native Go performance in mind
═══════════════════════════════════════════════════════════════
`
	fmt.Fprint(c.output, help)
	return nil
}

func (c *CLI) listFiles() error {
	files, err := c.project.Files()
	if err != nil {
		return err
	}

	c.files = files

	fmt.Fprintf(c.output, "\n📁 Files in %s:\n", c.project.Name())
	fmt.Fprint(c.output, "─────────────────────────────────────\n")

	for i, file := range files {
		icon := core.GetIconForLanguage(file.Language)
		fmt.Fprintf(c.output, "  %2d. %s %s\n", i+1, icon, file.RelPath)
	}

	fmt.Fprint(c.output, "─────────────────────────────────────\n")
	fmt.Fprintf(c.output, "Total: %d files\n\n", len(files))

	return nil
}

func (c *CLI) showTree() error {
	tree, err := c.project.FileTree()
	if err != nil {
		return err
	}

	fmt.Fprintf(c.output, "\n🌳 Project Structure: %s\n", c.project.Name())
	fmt.Fprint(c.output, "═══════════════════════════════════════════════════════════════\n")

	return c.renderer.RenderFileTree(c.output, tree)
}

func (c *CLI) openFile(filename string) error {
	filePath, err := c.resolveFile(filename)
	if err != nil {
		return err
	}

	c.currentFile = filePath
	fmt.Fprintf(c.output, "✅ Opened: %s\n", filename)
	fmt.Fprintf(c.output, "💡 Use 'cat %s' to view contents\n", filename)

	return nil
}

func (c *CLI) viewFile(filename string) error {
	filePath, err := c.resolveFile(filename)
	if err != nil {
		return err
	}

	// Find file info
	var fileInfo core.FileInfo
	for _, f := range c.files {
		if f.Path == filePath {
			fileInfo = f
			break
		}
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return c.renderer.RenderFile(c.output, fileInfo, string(content))
}

func (c *CLI) resolveFile(filename string) (string, error) {
	// Check if it's a number (index from ls command)
	if num, err := strconv.Atoi(filename); err == nil {
		if num < 1 || num > len(c.files) {
			return "", fmt.Errorf("invalid file number %d. Use 'ls' to see available files", num)
		}
		return c.files[num-1].Path, nil
	}

	// Direct filename
	for _, f := range c.files {
		if f.RelPath == filename || f.Name == filename {
			return f.Path, nil
		}
	}

	return "", fmt.Errorf("file not found: %s", filename)
}

func (c *CLI) runProject(ctx context.Context) error {
	fmt.Fprint(c.output, "🏃 Running Go project...\n")
	return c.builder.Run(ctx, c.project)
}

func (c *CLI) runTests(ctx context.Context) error {
	fmt.Fprint(c.output, "🧪 Running Go tests...\n")
	return c.builder.Test(ctx, c.project)
}

func (c *CLI) buildProject(ctx context.Context) error {
	fmt.Fprint(c.output, "🔨 Building Go project...\n")
	return c.builder.Build(ctx, c.project)
}

func (c *CLI) showVersion() error {
	version := `
🚀 GoX IDE v0.1.0-alpha
═══════════════════════════════════════════════════════════════
  Built with:     Pure Go + Gio UI
  Target:         Go-specialized development
  Features:       AI-native, Local-first, Fast
  Platform:       Cross-platform native
  
  🎯 The VS Code & GoLand killer for Go developers
═══════════════════════════════════════════════════════════════
`
	fmt.Fprint(c.output, version)
	return nil
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gox-ide/pkg/cli"
	"gox-ide/pkg/core"
	"gox-ide/pkg/filesystem"
)

// Version information (set by build flags)
var (
	Version   = "0.1.0-alpha"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Parse command line flags
	var (
		guiMode     = flag.Bool("gui", false, "Launch GUI mode")
		cliMode     = flag.Bool("cli", false, "Force CLI mode")
		projectPath = flag.String("project", "", "Project path (default: current directory)")
		showVersion = flag.Bool("version", false, "Show version information")
	)
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("GoX IDE %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		return
	}

	// Determine project path
	var path string
	if *projectPath != "" {
		path = *projectPath
	} else if len(flag.Args()) > 0 {
		path = flag.Args()[0]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("❌ Failed to get current directory:", err)
		}
		path = cwd
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal("❌ Failed to get absolute path:", err)
	}

	// Create dependencies
	fs := filesystem.NewOSFileSystem()
	project := core.NewGoProject(absPath, fs)
	logger := core.NewNoopLogger() // Use noop to avoid log noise
	builder := core.NewGoBuilder(logger)

	// Determine mode
	useGUI := *guiMode || (!*cliMode && shouldUseGUI())

	ctx := context.Background()

	if useGUI {
		// Check if GUI support is compiled in
		if !hasGUISupport() {
			log.Println("⚠️ GUI mode requested but not compiled in, falling back to CLI mode")
			useGUI = false
		}
	}

	if useGUI {
		// Launch GUI mode
		if err := runGUI(ctx, project, builder, logger); err != nil {
			log.Fatal("❌ GUI application error:", err)
		}
	} else {
		// Launch CLI mode
		renderer := cli.NewRenderer()
		cliApp := cli.New(cli.Config{
			Project:  project,
			Renderer: renderer,
			Builder:  builder,
			Logger:   logger,
			Input:    nil, // Will default to os.Stdin
			Output:   nil, // Will default to os.Stdout
		})

		if err := cliApp.Run(ctx); err != nil {
			log.Fatal("❌ CLI application error:", err)
		}
	}
}

// shouldUseGUI determines if GUI mode should be used by default
func shouldUseGUI() bool {
	// Check if we have a display (simple heuristic)
	if os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != "" {
		return true
	}

	// On Windows, assume GUI is available
	if os.Getenv("OS") == "Windows_NT" {
		return true
	}

	// Default to CLI for servers/headless environments
	return false
}

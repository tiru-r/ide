// Package core defines the core interfaces and types for GoX IDE
package core

import (
	"context"
	"io"
)

// Project represents a Go project with its metadata and operations
type Project interface {
	// Path returns the absolute path to the project
	Path() string

	// Name returns the project name
	Name() string

	// IsGoProject returns true if this is a valid Go project
	IsGoProject() bool

	// Files returns all files in the project
	Files() ([]FileInfo, error)

	// FileTree returns the project structure as a tree
	FileTree() (TreeNode, error)
}

// FileInfo represents information about a file in the project
type FileInfo struct {
	Name     string
	Path     string
	RelPath  string
	IsDir    bool
	Size     int64
	ModTime  int64
	Language string
}

// TreeNode represents a node in the project tree
type TreeNode struct {
	File     FileInfo
	Children []TreeNode
	Level    int
	IsLast   bool
}

// FileSystem provides file operations
type FileSystem interface {
	// ReadFile reads the contents of a file
	ReadFile(path string) ([]byte, error)

	// WriteFile writes contents to a file
	WriteFile(path string, data []byte) error

	// ListFiles lists files in a directory
	ListFiles(path string) ([]FileInfo, error)

	// WalkDir walks a directory tree
	WalkDir(path string, fn func(FileInfo) error) error

	// Exists checks if a file or directory exists
	Exists(path string) bool
}

// Editor provides text editing operations
type Editor interface {
	// OpenFile opens a file for editing
	OpenFile(path string) error

	// GetContent returns the current content
	GetContent() string

	// SetContent sets the content
	SetContent(content string)

	// Save saves the current content
	Save() error

	// IsDirty returns true if the content has been modified
	IsDirty() bool
}

// Builder provides build operations for Go projects
type Builder interface {
	// Build builds the project
	Build(ctx context.Context, project Project) error

	// Run runs the project
	Run(ctx context.Context, project Project) error

	// Test runs tests for the project
	Test(ctx context.Context, project Project) error

	// Clean cleans build artifacts
	Clean(ctx context.Context, project Project) error
}

// UI represents the user interface layer
type UI interface {
	// Start starts the UI
	Start(ctx context.Context) error

	// Stop stops the UI
	Stop() error

	// ShowProject displays the project
	ShowProject(project Project)

	// ShowError displays an error message
	ShowError(err error)

	// ShowMessage displays a message
	ShowMessage(msg string)
}

// CommandHandler handles commands from the UI
type CommandHandler interface {
	// HandleCommand processes a command and returns a response
	HandleCommand(ctx context.Context, cmd Command) (Response, error)

	// ListCommands returns available commands
	ListCommands() []CommandInfo
}

// Command represents a user command
type Command struct {
	Name string
	Args []string
}

// Response represents a command response
type Response struct {
	Output string
	Error  error
}

// CommandInfo provides information about a command
type CommandInfo struct {
	Name        string
	Description string
	Usage       string
}

// Logger provides logging capabilities
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

// Field represents a log field
type Field struct {
	Key   string
	Value any
}

// Config holds application configuration
type Config struct {
	ProjectPath string
	UIMode      string // "cli", "gui", "auto"
	LogLevel    string
	Editor      EditorConfig
}

// EditorConfig holds editor configuration
type EditorConfig struct {
	TabSize     int
	UseSpaces   bool
	ShowNumbers bool
}

// IDE is the main application interface
type IDE interface {
	// Start starts the IDE with the given configuration
	Start(ctx context.Context, config Config) error

	// Stop gracefully stops the IDE
	Stop() error

	// Project returns the current project
	Project() Project

	// SetProject sets the current project
	SetProject(project Project)
}

// Version information
type Version struct {
	Version   string
	BuildTime string
	GitCommit string
}

// Renderer renders output to a writer
type Renderer interface {
	// RenderProject renders project information
	RenderProject(w io.Writer, project Project) error

	// RenderFileTree renders a file tree
	RenderFileTree(w io.Writer, tree TreeNode) error

	// RenderFile renders file content
	RenderFile(w io.Writer, file FileInfo, content string) error

	// RenderError renders an error
	RenderError(w io.Writer, err error) error
}

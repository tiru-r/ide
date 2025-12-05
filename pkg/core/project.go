// Package core provides project management functionality for Go IDE.
package core

import (
	"path/filepath"
	"strings"
)

// GoProject implements the Project interface for Go projects
type GoProject struct {
	path string
	name string
	fs   FileSystem
}

// NewGoProject creates a new Go project instance
func NewGoProject(path string, fs FileSystem) *GoProject {
	return &GoProject{
		path: path,
		name: filepath.Base(path),
		fs:   fs,
	}
}

// Path returns the absolute path to the project
func (p *GoProject) Path() string {
	return p.path
}

// Name returns the project name
func (p *GoProject) Name() string {
	return p.name
}

// IsGoProject returns true if this is a valid Go project
func (p *GoProject) IsGoProject() bool {
	return p.fs.Exists(filepath.Join(p.path, "go.mod"))
}

// Files returns all files in the project
func (p *GoProject) Files() ([]FileInfo, error) {
	// Pre-allocate with estimated capacity for better performance
	files := make([]FileInfo, 0, 64)

	err := p.fs.WalkDir(p.path, func(info FileInfo) error {
		// Skip hidden files and directories
		if strings.HasPrefix(info.Name, ".") && info.Name != "." {
			return nil
		}

		// Skip vendor and node_modules
		if info.IsDir && (info.Name == "vendor" || info.Name == "node_modules") {
			return nil
		}

		if info.RelPath != "." && !info.IsDir {
			files = append(files, info)
		}

		return nil
	})

	return files, err
}

// FileTree returns the project structure as a tree
func (p *GoProject) FileTree() (TreeNode, error) {
	root := TreeNode{
		File: FileInfo{
			Name:    p.name,
			Path:    p.path,
			RelPath: ".",
			IsDir:   true,
		},
		Level: 0,
	}

	err := p.buildTree(&root, p.path, 0)
	return root, err
}

func (p *GoProject) buildTree(node *TreeNode, path string, level int) error {
	entries, err := p.fs.ListFiles(path)
	if err != nil {
		return err
	}

	// Filter visible entries with pre-allocation for performance
	visibleEntries := make([]FileInfo, 0, len(entries)) // Pre-allocate capacity
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name, ".") &&
			entry.Name != "vendor" && entry.Name != "node_modules" {
			visibleEntries = append(visibleEntries, entry)
		}
	}

	for i, entry := range visibleEntries {
		child := TreeNode{
			File:   entry,
			Level:  level + 1,
			IsLast: i == len(visibleEntries)-1,
		}

		if entry.IsDir {
			err := p.buildTree(&child, entry.Path, level+1)
			if err != nil {
				return err
			}
		}

		node.Children = append(node.Children, child)
	}

	return nil
}

// GetLanguageForFile returns the programming language for a file
func GetLanguageForFile(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".go":
		return "go"
	case ".mod", ".sum":
		return "gomod"
	case ".md":
		return "markdown"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".sh":
		return "shell"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".html":
		return "html"
	case ".css":
		return "css"
	default:
		return "text"
	}
}

// GetIconForLanguage returns an icon for a programming language
func GetIconForLanguage(language string) string {
	switch language {
	case "go":
		return "🐹"
	case "gomod":
		return "📦"
	case "markdown":
		return "📋"
	case "json":
		return "🔧"
	case "yaml", "toml":
		return "⚙️"
	case "shell":
		return "🖥️"
	case "python":
		return "🐍"
	case "javascript":
		return "📜"
	case "typescript":
		return "📘"
	case "html":
		return "🌐"
	case "css":
		return "🎨"
	default:
		return "📄"
	}
}

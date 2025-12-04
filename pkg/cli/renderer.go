package cli

import (
	"fmt"
	"io"
	"strings"

	"gox-ide/pkg/core"
)

// Renderer implements core.Renderer for CLI output
type Renderer struct{}

// NewRenderer creates a new CLI renderer
func NewRenderer() *Renderer {
	return &Renderer{}
}

// RenderProject renders project information
func (r *Renderer) RenderProject(w io.Writer, project core.Project) error {
	fmt.Fprintf(w, "📂 Project: %s\n", project.Name())
	fmt.Fprintf(w, "📁 Path: %s\n", project.Path())

	if project.IsGoProject() {
		fmt.Fprint(w, "🐹 Type: Go Project\n")
	} else {
		fmt.Fprint(w, "📄 Type: Generic Project\n")
	}

	return nil
}

// RenderFileTree renders a file tree
func (r *Renderer) RenderFileTree(w io.Writer, tree core.TreeNode) error {
	return r.renderTreeNode(w, tree, "")
}

func (r *Renderer) renderTreeNode(w io.Writer, node core.TreeNode, prefix string) error {
	if node.File.RelPath == "." {
		// Root node
		icon := "📁"
		if !node.File.IsDir {
			icon = core.GetIconForLanguage(node.File.Language)
		}
		fmt.Fprintf(w, "%s %s/\n", icon, node.File.Name)
	} else {
		// Child node
		connector := "├── "
		if node.IsLast {
			connector = "└── "
		}

		icon := "📁"
		suffix := "/"
		if !node.File.IsDir {
			icon = core.GetIconForLanguage(node.File.Language)
			suffix = ""
		}

		fmt.Fprintf(w, "%s%s%s %s%s\n", prefix, connector, icon, node.File.Name, suffix)
	}

	// Render children
	newPrefix := prefix
	if node.File.RelPath != "." {
		if node.IsLast {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
	}

	for _, child := range node.Children {
		if err := r.renderTreeNode(w, child, newPrefix); err != nil {
			return err
		}
	}

	return nil
}

// RenderFile renders file content with line numbers
func (r *Renderer) RenderFile(w io.Writer, file core.FileInfo, content string) error {
	lines := strings.Split(content, "\n")

	fmt.Fprintf(w, "\n📄 %s (%d lines)\n", file.RelPath, len(lines))
	fmt.Fprint(w, "═══════════════════════════════════════════════════════════════\n")

	for i, line := range lines {
		fmt.Fprintf(w, "%4d │ %s\n", i+1, line)
	}

	fmt.Fprint(w, "═══════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(w, "📊 File info: %d bytes, %d lines\n\n", len(content), len(lines))

	return nil
}

// RenderError renders an error message
func (r *Renderer) RenderError(w io.Writer, err error) error {
	fmt.Fprintf(w, "❌ Error: %v\n", err)
	return nil
}

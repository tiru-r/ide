package gui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gox-ide/pkg/core"
)

// FileExplorerImpl implements FileExplorer interface
type FileExplorerImpl struct {
	id           string
	project      core.Project
	tree         core.TreeNode
	list         widget.List
	items        []ExplorerItem
	selectedIdx  int
	onFileSelect func(file *core.FileInfo)
}

// ExplorerItem represents an item in the file explorer
type ExplorerItem struct {
	File     *core.FileInfo
	Button   widget.Clickable
	Level    int
	IsLast   bool
	Expanded bool
}

// NewFileExplorer creates a new file explorer component
func NewFileExplorer(project core.Project) *FileExplorerImpl {
	fe := &FileExplorerImpl{
		id:          "file-explorer",
		project:     project,
		selectedIdx: -1,
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}

	if project != nil {
		fe.loadFileTree()
	}

	return fe
}

// ID returns the component ID
func (fe *FileExplorerImpl) ID() string {
	return fe.id
}

// SetProject sets the project to display
func (fe *FileExplorerImpl) SetProject(project core.Project) {
	fe.project = project
	fe.loadFileTree()
}

// GetSelectedFile returns the currently selected file
func (fe *FileExplorerImpl) GetSelectedFile() *core.FileInfo {
	if fe.selectedIdx >= 0 && fe.selectedIdx < len(fe.items) {
		return fe.items[fe.selectedIdx].File
	}
	return nil
}

// SetOnFileSelect sets the callback for file selection
func (fe *FileExplorerImpl) SetOnFileSelect(callback func(file *core.FileInfo)) {
	fe.onFileSelect = callback
}

// Refresh reloads the file tree
func (fe *FileExplorerImpl) Refresh() error {
	fe.loadFileTree()
	return nil
}

// Update processes events and updates component state
func (fe *FileExplorerImpl) Update(gtx layout.Context) bool {
	changed := false

	// Handle item clicks
	for i := range fe.items {
		item := &fe.items[i]

		if item.Button.Clicked(gtx) {
			fe.selectedIdx = i

			if item.File.IsDir {
				// Toggle directory expansion
				item.Expanded = !item.Expanded
				fe.rebuildItems()
			} else {
				// File selection
				if fe.onFileSelect != nil {
					fe.onFileSelect(item.File)
				}
			}
			changed = true
		}
	}

	return changed
}

// Layout renders the file explorer
func (fe *FileExplorerImpl) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	// Update state
	fe.Update(gtx)

	// Draw background
	bg := color.NRGBA{R: 240, G: 240, B: 240, A: 255}
	paint.FillShape(gtx.Ops, bg, clip.Rect{Max: gtx.Constraints.Max}.Op())

	// Layout file list
	return layout.Inset{
		Top: 8, Bottom: 8, Left: 8, Right: 8,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &fe.list).Layout(gtx, len(fe.items), func(gtx layout.Context, i int) layout.Dimensions {
			return fe.layoutItem(gtx, theme, i)
		})
	})
}

// layoutItem renders a single file explorer item
func (fe *FileExplorerImpl) layoutItem(gtx layout.Context, theme *material.Theme, index int) layout.Dimensions {
	if index >= len(fe.items) {
		return layout.Dimensions{}
	}

	item := &fe.items[index]

	// Create indentation based on level
	leftInset := unit.Dp(float32(item.Level * 16))

	return layout.Inset{
		Left:   leftInset,
		Top:    unit.Dp(2),
		Bottom: unit.Dp(2),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		// Handle selection highlighting
		if index == fe.selectedIdx {
			selectionBG := color.NRGBA{R: 173, G: 216, B: 230, A: 255}
			paint.FillShape(gtx.Ops, selectionBG, clip.Rect{Max: gtx.Constraints.Max}.Op())
		}

		// Create button for the item
		btn := material.Button(theme, &item.Button, "")
		btn.Background = color.NRGBA{} // Transparent
		btn.Color = theme.Fg

		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Icon
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				icon := fe.getFileIcon(item.File)
				label := material.Body1(theme, icon)
				label.Color = theme.Fg
				return layout.Inset{Right: unit.Dp(4)}.Layout(gtx, label.Layout)
			}),

			// File name
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				name := item.File.Name
				if item.File.IsDir && item.Expanded {
					name = "▼ " + name
				} else if item.File.IsDir {
					name = "▶ " + name
				}

				btn := material.Button(theme, &item.Button, name)
				btn.Background = color.NRGBA{} // Transparent
				btn.Color = theme.Fg
				return btn.Layout(gtx)
			}),
		)
	})
}

// getFileIcon returns an appropriate icon for the file type
func (fe *FileExplorerImpl) getFileIcon(file *core.FileInfo) string {
	if file.IsDir {
		return "📁"
	}

	return core.GetIconForLanguage(file.Language)
}

// loadFileTree loads the file tree from the project
func (fe *FileExplorerImpl) loadFileTree() {
	if fe.project == nil {
		fe.items = nil
		return
	}

	tree, err := fe.project.FileTree()
	if err != nil {
		fe.items = nil
		return
	}

	fe.tree = tree
	fe.rebuildItems()
}

// rebuildItems rebuilds the flat list from the tree structure
func (fe *FileExplorerImpl) rebuildItems() {
	fe.items = nil
	if fe.tree.File.Name != "" {
		fe.addTreeNode(fe.tree, 0, true)
	}
}

// addTreeNode adds a tree node and its children to the flat list
func (fe *FileExplorerImpl) addTreeNode(node core.TreeNode, level int, isLast bool) {
	// Skip the root node (project name)
	if level > 0 {
		item := ExplorerItem{
			File:     &node.File,
			Level:    level - 1,
			IsLast:   isLast,
			Expanded: true, // Default to expanded for now
		}
		fe.items = append(fe.items, item)
	}

	// Add children if directory is expanded
	if node.File.IsDir {
		for i, child := range node.Children {
			isLastChild := i == len(node.Children)-1
			fe.addTreeNode(child, level+1, isLastChild)
		}
	}
}

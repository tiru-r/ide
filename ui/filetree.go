package ui

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type FileTree struct {
	projectPath string
	files       []FileItem
	list        widget.List
}

type FileItem struct {
	Name     string
	Path     string
	IsDir    bool
	Level    int
	Expanded bool
	Button   widget.Clickable
}

func NewFileTree(projectPath string) *FileTree {
	ft := &FileTree{
		projectPath: projectPath,
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}
	ft.loadFiles()
	return ft
}

func (ft *FileTree) loadFiles() {
	ft.files = nil
	err := filepath.WalkDir(ft.projectPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip vendor and node_modules
		if d.IsDir() && (d.Name() == "vendor" || d.Name() == "node_modules") {
			return filepath.SkipDir
		}

		relPath, _ := filepath.Rel(ft.projectPath, path)
		if relPath == "." {
			return nil
		}

		level := strings.Count(relPath, string(os.PathSeparator))

		item := FileItem{
			Name:  d.Name(),
			Path:  path,
			IsDir: d.IsDir(),
			Level: level,
		}

		ft.files = append(ft.files, item)
		return nil
	})

	if err != nil {
		// Handle error gracefully
		return
	}
}

func (ft *FileTree) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{
		Top: 8, Bottom: 8, Left: 8, Right: 8,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.List(th, &ft.list).Layout(gtx, len(ft.files), func(gtx layout.Context, i int) layout.Dimensions {
			return ft.files[i].Layout(gtx, th)
		})
	})
}

func (item *FileItem) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{
		Left:   unit.Dp(float32(item.Level * 16)),
		Top:    2,
		Bottom: 2,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		button := material.Button(th, &item.Button, item.Name)
		button.Background = th.Bg
		button.Color = th.Fg

		if item.IsDir {
			button.Color = th.ContrastFg
		}

		return button.Layout(gtx)
	})
}

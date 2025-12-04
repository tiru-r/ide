package ui

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type IDE struct {
	projectPath string
	theme       *material.Theme
	filetree    *FileTree
	editor      *Editor
}

func NewIDE(projectPath string) *IDE {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	return &IDE{
		projectPath: projectPath,
		theme:       th,
		filetree:    NewFileTree(projectPath),
		editor:      NewEditor(),
	}
}

func (ide *IDE) Run() error {
	w := new(app.Window)
	w.Option(
		app.Title("GoX IDE"),
		app.Size(unit.Dp(1200), unit.Dp(800)),
	)

	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return nil
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ide.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (ide *IDE) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return ide.filetree.Layout(gtx, ide.theme)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return ide.editor.Layout(gtx, ide.theme)
		}),
	)
}

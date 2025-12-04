package ui

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Editor struct {
	editor widget.Editor
}

func NewEditor() *Editor {
	return &Editor{
		editor: widget.Editor{
			SingleLine: false,
		},
	}
}

func (e *Editor) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Inset{
		Top: 8, Bottom: 8, Left: 8, Right: 8,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		ed := material.Editor(th, &e.editor, "// Welcome to GoX IDE\n// Start coding...")
		return ed.Layout(gtx)
	})
}

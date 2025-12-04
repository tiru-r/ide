package gui

import (
	"fmt"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"gox-ide/pkg/core"
)

// StatusBarImpl implements StatusBar interface
type StatusBarImpl struct {
	id          string
	message     string
	fileInfo    string
	projectInfo string
}

// NewStatusBar creates a new status bar component
func NewStatusBar() *StatusBarImpl {
	return &StatusBarImpl{
		id:      "status-bar",
		message: "Ready",
	}
}

// ID returns the component ID
func (sb *StatusBarImpl) ID() string {
	return sb.id
}

// SetMessage sets the status message
func (sb *StatusBarImpl) SetMessage(message string) {
	sb.message = message
}

// SetFileInfo sets file information display
func (sb *StatusBarImpl) SetFileInfo(file *core.FileInfo, line, col int) {
	if file != nil {
		sb.fileInfo = fmt.Sprintf("%s  Ln %d, Col %d", file.Name, line, col)
	} else {
		sb.fileInfo = ""
	}
}

// SetProjectInfo sets project information
func (sb *StatusBarImpl) SetProjectInfo(project core.Project) {
	if project != nil {
		sb.projectInfo = project.Name()
	} else {
		sb.projectInfo = ""
	}
}

// Update processes events and updates component state
func (sb *StatusBarImpl) Update(gtx layout.Context) bool {
	// Status bar is mostly passive
	return false
}

// Layout renders the status bar
func (sb *StatusBarImpl) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	// Draw status bar background
	bg := color.NRGBA{R: 230, G: 230, B: 230, A: 255}
	paint.FillShape(gtx.Ops, bg, clip.Rect{Max: gtx.Constraints.Max}.Op())

	// Layout content with padding
	return layout.Inset{
		Top: unit.Dp(4), Bottom: unit.Dp(4),
		Left: unit.Dp(8), Right: unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx,
			// Message on the left
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if sb.message == "" {
					sb.message = "Ready"
				}

				label := material.Caption(theme, sb.message)
				label.Color = theme.Fg
				return label.Layout(gtx)
			}),

			// Project info in the middle
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if sb.projectInfo == "" {
					return layout.Dimensions{}
				}

				label := material.Caption(theme, fmt.Sprintf("📁 %s", sb.projectInfo))
				label.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
				return layout.Inset{Left: unit.Dp(16), Right: unit.Dp(16)}.Layout(gtx, label.Layout)
			}),

			// File info on the right
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if sb.fileInfo == "" {
					return layout.Dimensions{}
				}

				label := material.Caption(theme, sb.fileInfo)
				label.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
				return label.Layout(gtx)
			}),
		)
	})
}

package gui

import (
	"image/color"
	"os"
	"strings"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gox-ide/pkg/core"
)

// TextEditorImpl implements Editor interface
type TextEditorImpl struct {
	id          string
	editor      widget.Editor
	currentFile *core.FileInfo
	dirty       bool
	onChange    func()
}

// NewTextEditor creates a new text editor component
func NewTextEditor() *TextEditorImpl {
	return &TextEditorImpl{
		id: "text-editor",
		editor: widget.Editor{
			SingleLine: false,
			Submit:     false,
		},
		dirty: false,
	}
}

// ID returns the component ID
func (te *TextEditorImpl) ID() string {
	return te.id
}

// OpenFile opens a file for editing
func (te *TextEditorImpl) OpenFile(file *core.FileInfo) error {
	if file == nil || file.IsDir {
		return nil
	}

	// Read file content
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return err
	}

	// Set content and file
	te.editor.SetText(string(content))
	te.currentFile = file
	te.dirty = false

	return nil
}

// GetContent returns the current editor content
func (te *TextEditorImpl) GetContent() string {
	return te.editor.Text()
}

// SetContent sets the editor content
func (te *TextEditorImpl) SetContent(content string) {
	te.editor.SetText(content)
	te.markDirty()
}

// Save saves the current content to file
func (te *TextEditorImpl) Save() error {
	if te.currentFile == nil {
		return nil
	}

	content := te.GetContent()
	err := os.WriteFile(te.currentFile.Path, []byte(content), 0644)
	if err == nil {
		te.dirty = false
	}

	return err
}

// IsDirty returns true if content has been modified
func (te *TextEditorImpl) IsDirty() bool {
	return te.dirty
}

// SetOnChange sets the callback for content changes
func (te *TextEditorImpl) SetOnChange(callback func()) {
	te.onChange = callback
}

// GetCurrentFile returns the currently open file
func (te *TextEditorImpl) GetCurrentFile() *core.FileInfo {
	return te.currentFile
}

// Update processes events and updates component state
func (te *TextEditorImpl) Update(gtx layout.Context) bool {
	// For now, we'll check for changes in the Layout method
	// The Gio editor doesn't have a Changed() method in newer versions
	return false
}

// Layout renders the text editor
func (te *TextEditorImpl) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	// Update state
	te.Update(gtx)

	// Draw editor background
	bg := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	paint.FillShape(gtx.Ops, bg, clip.Rect{Max: gtx.Constraints.Max}.Op())

	// Layout editor with padding
	return layout.Inset{
		Top: unit.Dp(8), Bottom: unit.Dp(8),
		Left: unit.Dp(8), Right: unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		if te.currentFile == nil {
			// Show welcome message
			return te.layoutWelcome(gtx, theme)
		}

		// Show file content with line numbers
		return te.layoutEditor(gtx, theme)
	})
}

// layoutWelcome shows a welcome message when no file is open
func (te *TextEditorImpl) layoutWelcome(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H4(theme, "Welcome to GoX IDE")
				title.Color = theme.Fg
				return layout.Inset{Bottom: unit.Dp(16)}.Layout(gtx, title.Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				subtitle := material.Body1(theme, "Open a file from the explorer to start editing")
				subtitle.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
				return subtitle.Layout(gtx)
			}),
		)
	})
}

// layoutEditor shows the actual editor with line numbers
func (te *TextEditorImpl) layoutEditor(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(gtx,
		// Line numbers
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return te.layoutLineNumbers(gtx, theme)
		}),

		// Editor content
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			ed := material.Editor(theme, &te.editor, "")
			ed.Color = theme.Fg
			return ed.Layout(gtx)
		}),
	)
}

// layoutLineNumbers renders line numbers on the left side
func (te *TextEditorImpl) layoutLineNumbers(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	content := te.GetContent()
	lines := strings.Split(content, "\n")
	lineCount := len(lines)

	return layout.Inset{Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Draw line numbers background
		lnBg := color.NRGBA{R: 245, G: 245, B: 245, A: 255}
		paint.FillShape(gtx.Ops, lnBg, clip.Rect{Max: gtx.Constraints.Max}.Op())

		return layout.Inset{
			Top: unit.Dp(4), Bottom: unit.Dp(4),
			Left: unit.Dp(8), Right: unit.Dp(8),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

			// Create line number text
			var lineNumbers strings.Builder
			for i := 1; i <= lineCount; i++ {
				if i > 1 {
					lineNumbers.WriteString("\n")
				}
				lineNumbers.WriteString(strings.ReplaceAll(sprintf("%4d", i), " ", " "))
			}

			// Render line numbers
			lnLabel := material.Body2(theme, lineNumbers.String())
			lnLabel.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}

			// Constrain width for line numbers
			gtx.Constraints.Max.X = gtx.Dp(unit.Dp(60))
			gtx.Constraints.Min.X = gtx.Constraints.Max.X

			return lnLabel.Layout(gtx)
		})
	})
}

// markDirty marks the editor content as modified
func (te *TextEditorImpl) markDirty() {
	if !te.dirty {
		te.dirty = true
		if te.onChange != nil {
			te.onChange()
		}
	}
}

// sprintf is a simple sprintf replacement since we can't import fmt
func sprintf(format string, value int) string {
	// Simple integer formatting for line numbers
	str := ""
	if value == 0 {
		return "   0"
	}

	num := value
	digits := []rune{}

	for num > 0 {
		digit := num % 10
		digits = append([]rune{rune('0' + digit)}, digits...)
		num /= 10
	}

	str = string(digits)

	// Pad to 4 characters
	for len(str) < 4 {
		str = " " + str
	}

	return str
}

package gui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// ToolBarImpl implements ToolBar interface
type ToolBarImpl struct {
	id      string
	buttons []ToolBarButton
	actions map[string]func()
}

// ToolBarButton represents a toolbar button
type ToolBarButton struct {
	ID      string
	Text    string
	Icon    string
	Enabled bool
	Button  widget.Clickable
}

// NewToolBar creates a new toolbar component
func NewToolBar() *ToolBarImpl {
	tb := &ToolBarImpl{
		id:      "toolbar",
		actions: make(map[string]func()),
	}

	// Initialize default buttons
	tb.buttons = []ToolBarButton{
		{ID: "save", Text: "Save", Icon: "💾", Enabled: false},
		{ID: "build", Text: "Build", Icon: "🔨", Enabled: true},
		{ID: "run", Text: "Run", Icon: "▶️", Enabled: true},
		{ID: "test", Text: "Test", Icon: "🧪", Enabled: true},
	}

	return tb
}

// ID returns the component ID
func (tb *ToolBarImpl) ID() string {
	return tb.id
}

// SetOnAction sets callback for toolbar actions
func (tb *ToolBarImpl) SetOnAction(action string, callback func()) {
	tb.actions[action] = callback
}

// EnableAction enables/disables an action
func (tb *ToolBarImpl) EnableAction(action string, enabled bool) {
	for i := range tb.buttons {
		if tb.buttons[i].ID == action {
			tb.buttons[i].Enabled = enabled
			break
		}
	}
}

// AddSeparator adds a visual separator
func (tb *ToolBarImpl) AddSeparator() {
	tb.buttons = append(tb.buttons, ToolBarButton{
		ID: "separator",
	})
}

// Update processes events and updates component state
func (tb *ToolBarImpl) Update(gtx layout.Context) bool {
	changed := false

	// Handle button clicks
	for i := range tb.buttons {
		button := &tb.buttons[i]

		if button.ID == "separator" {
			continue
		}

		if button.Button.Clicked(gtx) && button.Enabled {
			if action, exists := tb.actions[button.ID]; exists && action != nil {
				action()
			}
			changed = true
		}
	}

	return changed
}

// Layout renders the toolbar
func (tb *ToolBarImpl) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	// Update state
	tb.Update(gtx)

	// Draw toolbar background
	bg := color.NRGBA{R: 245, G: 245, B: 245, A: 255}
	paint.FillShape(gtx.Ops, bg, clip.Rect{Max: gtx.Constraints.Max}.Op())

	// Layout buttons with padding
	return layout.Inset{
		Top: unit.Dp(4), Bottom: unit.Dp(4),
		Left: unit.Dp(8), Right: unit.Dp(8),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		// Create flex children for buttons
		var children []layout.FlexChild

		for i := range tb.buttons {
			button := &tb.buttons[i]

			if button.ID == "separator" {
				// Add separator
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						// Draw vertical line
						sepColor := color.NRGBA{R: 200, G: 200, B: 200, A: 255}
						paint.FillShape(gtx.Ops, sepColor, clip.Rect{
							Max: image.Point{X: gtx.Dp(unit.Dp(1)), Y: gtx.Constraints.Max.Y},
						}.Op())
						return layout.Dimensions{Size: image.Point{X: gtx.Dp(unit.Dp(1)), Y: gtx.Constraints.Max.Y}}
					})
				}))
			} else {
				// Add button
				button := button // Capture for closure
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return tb.layoutButton(gtx, theme, button)
				}))
			}
		}

		return layout.Flex{
			Axis:      layout.Horizontal,
			Alignment: layout.Middle,
		}.Layout(gtx, children...)
	})
}

// layoutButton renders a single toolbar button
func (tb *ToolBarImpl) layoutButton(gtx layout.Context, theme *material.Theme, button *ToolBarButton) layout.Dimensions {
	return layout.Inset{Right: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

		// Create button style
		btn := material.Button(theme, &button.Button, "")

		if button.Enabled {
			btn.Background = theme.Bg
			btn.Color = theme.Fg
		} else {
			btn.Background = color.NRGBA{R: 240, G: 240, B: 240, A: 255}
			btn.Color = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
		}

		// Simple button layout
		if button.Icon != "" {
			btn.Text = button.Icon + " " + button.Text
		} else {
			btn.Text = button.Text
		}

		return btn.Layout(gtx)
	})
}

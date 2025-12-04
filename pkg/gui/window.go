package gui

import (
	"context"
	"fmt"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"gox-ide/pkg/core"
)

// Window implements IDEWindow
type Window struct {
	config       IDEConfig
	window       *app.Window
	theme        *Theme
	fileExplorer FileExplorer
	editor       Editor
	statusBar    StatusBar
	toolBar      ToolBar

	// State
	running bool
}

// NewWindow creates a new IDE window
func NewWindow(config IDEConfig) *Window {
	theme := config.Theme
	if theme == nil {
		theme = NewDefaultTheme()
	}

	w := &Window{
		config: config,
		theme:  theme,
		window: new(app.Window),
	}

	// Initialize components
	w.fileExplorer = NewFileExplorer(config.Project)
	w.editor = NewTextEditor()
	w.statusBar = NewStatusBar()
	w.toolBar = NewToolBar()

	// Setup event handlers
	w.setupEventHandlers()

	// Configure window
	w.window.Option(
		app.Title("GoX IDE"),
		app.Size(unit.Dp(1200), unit.Dp(800)),
	)

	return w
}

// Run starts the IDE window
func (w *Window) Run(ctx context.Context) error {
	w.running = true
	defer func() { w.running = false }()

	var ops op.Ops

	for w.running {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			switch e := w.window.Event().(type) {
			case app.DestroyEvent:
				return nil
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				w.layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}

	return nil
}

// Close closes the IDE window
func (w *Window) Close() {
	w.running = false
}

// SetTitle sets the window title
func (w *Window) SetTitle(title string) {
	w.window.Option(app.Title(title))
}

// SetProject sets the current project
func (w *Window) SetProject(project core.Project) {
	w.config.Project = project
	if w.fileExplorer != nil {
		w.fileExplorer.SetProject(project)
	}
	if w.statusBar != nil {
		w.statusBar.SetProjectInfo(project)
	}
	w.updateTitle()
}

// GetFileExplorer returns the file explorer component
func (w *Window) GetFileExplorer() FileExplorer {
	return w.fileExplorer
}

// GetEditor returns the editor component
func (w *Window) GetEditor() Editor {
	return w.editor
}

// GetStatusBar returns the status bar component
func (w *Window) GetStatusBar() StatusBar {
	return w.statusBar
}

// ShowMessage displays a message to the user
func (w *Window) ShowMessage(message string) {
	if w.statusBar != nil {
		w.statusBar.SetMessage(message)
	}
}

// ShowError displays an error message
func (w *Window) ShowError(err error) {
	if w.statusBar != nil {
		w.statusBar.SetMessage(fmt.Sprintf("Error: %v", err))
	}
}

// setupEventHandlers configures component event handlers
func (w *Window) setupEventHandlers() {
	// File explorer file selection
	if w.fileExplorer != nil {
		w.fileExplorer.SetOnFileSelect(w.onFileSelect)
	}

	// Editor change handler
	if w.editor != nil {
		w.editor.SetOnChange(w.onEditorChange)
	}

	// Toolbar actions
	if w.toolBar != nil {
		w.setupToolbarActions()
	}
}

// onFileSelect handles file selection from explorer
func (w *Window) onFileSelect(file *core.FileInfo) {
	if file == nil || file.IsDir {
		return
	}

	// Open file in editor
	if err := w.editor.OpenFile(file); err != nil {
		w.ShowError(fmt.Errorf("failed to open file: %w", err))
		return
	}

	// Update status bar
	w.statusBar.SetFileInfo(file, 1, 1) // TODO: Get actual cursor position

	// Notify event handler
	if w.config.EventHandler != nil {
		w.config.EventHandler.OnFileOpen(file)
	}
}

// onEditorChange handles editor content changes
func (w *Window) onEditorChange() {
	if file := w.editor.GetCurrentFile(); file != nil {
		// Update title to show unsaved changes
		title := fmt.Sprintf("GoX IDE - %s*", file.Name)
		w.SetTitle(title)

		// Update status
		w.statusBar.SetMessage("Modified")
	}
}

// setupToolbarActions configures toolbar button actions
func (w *Window) setupToolbarActions() {
	// Save action
	w.toolBar.SetOnAction("save", func() {
		if err := w.editor.Save(); err != nil {
			w.ShowError(fmt.Errorf("failed to save: %w", err))
		} else {
			w.ShowMessage("File saved")
			w.updateTitle() // Remove asterisk
		}
	})

	// Build action
	w.toolBar.SetOnAction("build", func() {
		w.ShowMessage("Building...")
		if w.config.EventHandler != nil {
			if err := w.config.EventHandler.OnBuild(); err != nil {
				w.ShowError(err)
			} else {
				w.ShowMessage("Build successful")
			}
		}
	})

	// Run action
	w.toolBar.SetOnAction("run", func() {
		w.ShowMessage("Running...")
		if w.config.EventHandler != nil {
			if err := w.config.EventHandler.OnRun(); err != nil {
				w.ShowError(err)
			} else {
				w.ShowMessage("Execution completed")
			}
		}
	})

	// Test action
	w.toolBar.SetOnAction("test", func() {
		w.ShowMessage("Running tests...")
		if w.config.EventHandler != nil {
			if err := w.config.EventHandler.OnTest(); err != nil {
				w.ShowError(err)
			} else {
				w.ShowMessage("Tests passed")
			}
		}
	})
}

// layout renders the main IDE layout
func (w *Window) layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// Toolbar
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.toolBar.Layout(gtx, w.theme.Theme)
		}),

		// Main content area
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// File explorer sidebar
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Max.X = gtx.Dp(unit.Dp(250))
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return w.fileExplorer.Layout(gtx, w.theme.Theme)
				}),

				// Editor area
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return w.editor.Layout(gtx, w.theme.Theme)
				}),
			)
		}),

		// Status bar
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return w.statusBar.Layout(gtx, w.theme.Theme)
		}),
	)
}

// updateTitle updates the window title based on current state
func (w *Window) updateTitle() {
	title := "GoX IDE"

	if w.config.Project != nil {
		title = fmt.Sprintf("GoX IDE - %s", w.config.Project.Name())
	}

	if file := w.editor.GetCurrentFile(); file != nil {
		title = fmt.Sprintf("GoX IDE - %s", file.Name)
		if w.editor.IsDirty() {
			title += "*"
		}
	}

	w.SetTitle(title)
}

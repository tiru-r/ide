package gui

import (
	"context"
	"fmt"

	"gox-ide/pkg/core"
)

// IDEApp integrates the GUI with the core IDE functionality
type IDEApp struct {
	config  IDEConfig
	window  IDEWindow
	project core.Project
	builder core.Builder
	logger  core.Logger
}

// NewIDEApp creates a new GUI IDE application
func NewIDEApp(project core.Project, builder core.Builder, logger core.Logger) *IDEApp {
	app := &IDEApp{
		project: project,
		builder: builder,
		logger:  logger,
	}

	// Create event handler
	eventHandler := &ideEventHandler{
		app: app,
	}

	// Setup configuration
	app.config = IDEConfig{
		Project:      project,
		Builder:      builder,
		Logger:       logger,
		Theme:        NewDefaultTheme(),
		EventHandler: eventHandler,
	}

	// Create window
	app.window = NewWindow(app.config)

	return app
}

// Run starts the GUI IDE
func (app *IDEApp) Run(ctx context.Context) error {
	if app.logger != nil {
		app.logger.Info("Starting GUI IDE", core.Field{Key: "project", Value: app.project.Path()})
	}

	// Set initial project
	app.window.SetProject(app.project)

	// Run the window
	return app.window.Run(ctx)
}

// Close closes the GUI IDE
func (app *IDEApp) Close() {
	if app.window != nil {
		app.window.Close()
	}
}

// ideEventHandler implements EventHandler interface
type ideEventHandler struct {
	app *IDEApp
}

// OnFileOpen handles file opening
func (h *ideEventHandler) OnFileOpen(file *core.FileInfo) {
	if h.app.logger != nil {
		h.app.logger.Info("File opened", core.Field{Key: "file", Value: file.Path})
	}

	// Update status
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		statusBar.SetMessage(fmt.Sprintf("Opened %s", file.Name))
		statusBar.SetFileInfo(file, 1, 1)
	}

	// Enable save action
	if toolbar := h.app.window.(*Window).toolBar; toolbar != nil {
		toolbar.EnableAction("save", true)
	}
}

// OnFileClose handles file closing
func (h *ideEventHandler) OnFileClose(file *core.FileInfo) {
	if h.app.logger != nil {
		h.app.logger.Info("File closed", core.Field{Key: "file", Value: file.Path})
	}

	// Update status
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		statusBar.SetMessage("Ready")
		statusBar.SetFileInfo(nil, 0, 0)
	}

	// Disable save action
	if toolbar := h.app.window.(*Window).toolBar; toolbar != nil {
		toolbar.EnableAction("save", false)
	}
}

// OnFileSave handles file saving
func (h *ideEventHandler) OnFileSave(file *core.FileInfo) error {
	if h.app.logger != nil {
		h.app.logger.Info("File saved", core.Field{Key: "file", Value: file.Path})
	}

	// Update status
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		statusBar.SetMessage(fmt.Sprintf("Saved %s", file.Name))
	}

	return nil
}

// OnProjectChange handles project changes
func (h *ideEventHandler) OnProjectChange(project core.Project) {
	h.app.project = project

	if h.app.logger != nil {
		h.app.logger.Info("Project changed", core.Field{Key: "project", Value: project.Path()})
	}

	// Update window
	h.app.window.SetProject(project)

	// Update status
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		statusBar.SetProjectInfo(project)
		statusBar.SetMessage(fmt.Sprintf("Project: %s", project.Name()))
	}
}

// OnBuild handles build requests
func (h *ideEventHandler) OnBuild() error {
	if h.app.logger != nil {
		h.app.logger.Info("Build started", core.Field{Key: "project", Value: h.app.project.Path()})
	}

	// Update status
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		statusBar.SetMessage("Building...")
	}

	// Execute build
	ctx := context.Background()
	err := h.app.builder.Build(ctx, h.app.project)

	// Update status based on result
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		if err != nil {
			statusBar.SetMessage(fmt.Sprintf("Build failed: %v", err))
		} else {
			statusBar.SetMessage("Build successful")
		}
	}

	if h.app.logger != nil {
		if err != nil {
			h.app.logger.Error("Build failed", core.Field{Key: "error", Value: err.Error()})
		} else {
			h.app.logger.Info("Build completed successfully")
		}
	}

	return err
}

// OnRun handles run requests
func (h *ideEventHandler) OnRun() error {
	if h.app.logger != nil {
		h.app.logger.Info("Run started", core.Field{Key: "project", Value: h.app.project.Path()})
	}

	// Update status
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		statusBar.SetMessage("Running...")
	}

	// Execute run
	ctx := context.Background()
	err := h.app.builder.Run(ctx, h.app.project)

	// Update status based on result
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		if err != nil {
			statusBar.SetMessage(fmt.Sprintf("Run failed: %v", err))
		} else {
			statusBar.SetMessage("Execution completed")
		}
	}

	if h.app.logger != nil {
		if err != nil {
			h.app.logger.Error("Run failed", core.Field{Key: "error", Value: err.Error()})
		} else {
			h.app.logger.Info("Run completed successfully")
		}
	}

	return err
}

// OnTest handles test requests
func (h *ideEventHandler) OnTest() error {
	if h.app.logger != nil {
		h.app.logger.Info("Test started", core.Field{Key: "project", Value: h.app.project.Path()})
	}

	// Update status
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		statusBar.SetMessage("Running tests...")
	}

	// Execute tests
	ctx := context.Background()
	err := h.app.builder.Test(ctx, h.app.project)

	// Update status based on result
	if statusBar := h.app.window.GetStatusBar(); statusBar != nil {
		if err != nil {
			statusBar.SetMessage(fmt.Sprintf("Tests failed: %v", err))
		} else {
			statusBar.SetMessage("All tests passed")
		}
	}

	if h.app.logger != nil {
		if err != nil {
			h.app.logger.Error("Tests failed", core.Field{Key: "error", Value: err.Error()})
		} else {
			h.app.logger.Info("Tests completed successfully")
		}
	}

	return err
}

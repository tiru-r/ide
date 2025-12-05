// Package gui provides graphical user interface components
package gui

import (
	"context"
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget/material"

	"gox-ide/pkg/core"
)

// Component represents a reusable UI component
type Component interface {
	// Layout renders the component and returns its dimensions
	Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions

	// Update processes events and updates component state
	Update(gtx layout.Context) bool

	// ID returns a unique identifier for this component
	ID() string
}

// Container manages child components
type Container interface {
	Component

	// AddChild adds a component as a child
	AddChild(child Component)

	// RemoveChild removes a child component
	RemoveChild(id string) bool

	// GetChild returns a child component by ID
	GetChild(id string) Component

	// Children returns all child components
	Children() []Component
}

// FileExplorer displays project files in a tree structure
type FileExplorer interface {
	Component

	// SetProject sets the project to display
	SetProject(project core.Project)

	// GetSelectedFile returns the currently selected file
	GetSelectedFile() *core.FileInfo

	// SetOnFileSelect sets the callback for file selection
	SetOnFileSelect(callback func(file *core.FileInfo))

	// Refresh reloads the file tree
	Refresh() error
}

// Editor provides text editing functionality
type Editor interface {
	Component

	// OpenFile opens a file for editing
	OpenFile(file *core.FileInfo) error

	// GetContent returns the current editor content
	GetContent() string

	// SetContent sets the editor content
	SetContent(content string)

	// Save saves the current content to file
	Save() error

	// IsDirty returns true if content has been modified
	IsDirty() bool

	// SetOnChange sets the callback for content changes
	SetOnChange(callback func())

	// GetCurrentFile returns the currently open file
	GetCurrentFile() *core.FileInfo
}

// StatusBar displays status information
type StatusBar interface {
	Component

	// SetMessage sets the status message
	SetMessage(message string)

	// SetFileInfo sets file information display
	SetFileInfo(file *core.FileInfo, line, col int)

	// SetProjectInfo sets project information
	SetProjectInfo(project core.Project)
}

// ToolBar provides quick action buttons
type ToolBar interface {
	Component

	// SetOnAction sets callback for toolbar actions
	SetOnAction(action string, callback func())

	// EnableAction enables/disables an action
	EnableAction(action string, enabled bool)

	// AddSeparator adds a visual separator
	AddSeparator()
}

// IDEWindow is the main application window
type IDEWindow interface {
	// Run starts the IDE window event loop
	Run(ctx context.Context) error

	// Close closes the IDE window
	Close()

	// SetTitle sets the window title
	SetTitle(title string)

	// SetProject sets the current project
	SetProject(project core.Project)

	// GetFileExplorer returns the file explorer component
	GetFileExplorer() FileExplorer

	// GetEditor returns the editor component
	GetEditor() Editor

	// GetStatusBar returns the status bar component
	GetStatusBar() StatusBar

	// GetToolBar returns the toolbar component
	GetToolBar() ToolBar

	// ShowMessage displays a message to the user
	ShowMessage(message string)

	// ShowError displays an error message
	ShowError(err error)
}

// Theme defines the visual appearance
type Theme struct {
	*material.Theme

	// Custom colors
	Primary      color.NRGBA
	Secondary    color.NRGBA
	Background   color.NRGBA
	Surface      color.NRGBA
	Error        color.NRGBA
	OnPrimary    color.NRGBA
	OnSecondary  color.NRGBA
	OnBackground color.NRGBA
	OnSurface    color.NRGBA
	OnError      color.NRGBA

	// IDE specific colors
	EditorBackground color.NRGBA
	EditorText       color.NRGBA
	LineNumbers      color.NRGBA
	SelectionBG      color.NRGBA
	CurrentLine      color.NRGBA
	SidebarBG        color.NRGBA
	StatusBarBG      color.NRGBA
	ToolBarBG        color.NRGBA
}

// EventHandler handles IDE events
type EventHandler interface {
	// OnFileOpen handles file opening
	OnFileOpen(file *core.FileInfo)

	// OnFileClose handles file closing
	OnFileClose(file *core.FileInfo)

	// OnFileSave handles file saving
	OnFileSave(file *core.FileInfo) error

	// OnProjectChange handles project changes
	OnProjectChange(project core.Project)

	// OnBuild handles build requests
	OnBuild() error

	// OnRun handles run requests
	OnRun() error

	// OnTest handles test requests
	OnTest() error
}

// ComponentFactory creates GUI components with loose coupling
type ComponentFactory interface {
	CreateFileExplorer(project core.Project) FileExplorer
	CreateEditor() Editor
	CreateStatusBar() StatusBar
	CreateToolBar() ToolBar
}

// IDEConfig holds configuration for the IDE
type IDEConfig struct {
	Project      core.Project
	Builder      core.Builder
	Logger       core.Logger
	Theme        *Theme
	EventHandler EventHandler

	// Component injection (optional - falls back to factory)
	FileExplorer FileExplorer
	Editor       Editor
	StatusBar    StatusBar
	ToolBar      ToolBar

	// Factory for creating components
	Factory ComponentFactory
}

// DefaultComponentFactory implements ComponentFactory with default components
type DefaultComponentFactory struct{}

// CreateFileExplorer creates a default file explorer
func (f *DefaultComponentFactory) CreateFileExplorer(project core.Project) FileExplorer {
	return NewFileExplorer(project)
}

// CreateEditor creates a default text editor
func (f *DefaultComponentFactory) CreateEditor() Editor {
	return NewTextEditor()
}

// CreateStatusBar creates a default status bar
func (f *DefaultComponentFactory) CreateStatusBar() StatusBar {
	return NewStatusBar()
}

// CreateToolBar creates a default toolbar
func (f *DefaultComponentFactory) CreateToolBar() ToolBar {
	return NewToolBar()
}

// NewDefaultFactory creates a default component factory
func NewDefaultFactory() ComponentFactory {
	return &DefaultComponentFactory{}
}

// NewDefaultTheme creates a default IDE theme
func NewDefaultTheme() *Theme {
	base := material.NewTheme()

	return &Theme{
		Theme:      base,
		Primary:    color.NRGBA{R: 33, G: 150, B: 243, A: 255},  // Blue
		Secondary:  color.NRGBA{R: 255, G: 193, B: 7, A: 255},   // Amber
		Background: color.NRGBA{R: 250, G: 250, B: 250, A: 255}, // Light gray
		Surface:    color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // White
		Error:      color.NRGBA{R: 244, G: 67, B: 54, A: 255},   // Red

		OnPrimary:    color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		OnSecondary:  color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		OnBackground: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		OnSurface:    color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		OnError:      color.NRGBA{R: 255, G: 255, B: 255, A: 255},

		EditorBackground: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		EditorText:       color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		LineNumbers:      color.NRGBA{R: 150, G: 150, B: 150, A: 255},
		SelectionBG:      color.NRGBA{R: 173, G: 216, B: 230, A: 255},
		CurrentLine:      color.NRGBA{R: 245, G: 245, B: 245, A: 255},
		SidebarBG:        color.NRGBA{R: 240, G: 240, B: 240, A: 255},
		StatusBarBG:      color.NRGBA{R: 230, G: 230, B: 230, A: 255},
		ToolBarBG:        color.NRGBA{R: 245, G: 245, B: 245, A: 255},
	}
}

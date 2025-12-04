//go:build gui

package main

import (
	"context"
	"log"

	"gox-ide/pkg/core"
	"gox-ide/pkg/gui"
)

func runGUI(ctx context.Context, project core.Project, builder core.Builder, logger core.Logger) error {
	log.Println("🚀 Starting GoX IDE in GUI mode...")
	app := gui.NewIDEApp(project, builder, logger)
	return app.Run(ctx)
}

func hasGUISupport() bool {
	return true
}

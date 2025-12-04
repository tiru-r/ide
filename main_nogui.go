//go:build !gui

package main

import (
	"context"
	"errors"

	"gox-ide/pkg/core"
)

func runGUI(ctx context.Context, project core.Project, builder core.Builder, logger core.Logger) error {
	return errors.New("GUI support not compiled in")
}

func hasGUISupport() bool {
	return false
}

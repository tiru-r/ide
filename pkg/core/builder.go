package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GoBuilder implements Builder interface for Go projects
type GoBuilder struct {
	logger Logger
}

// NewGoBuilder creates a new Go builder
func NewGoBuilder(logger Logger) *GoBuilder {
	return &GoBuilder{
		logger: logger,
	}
}

// Build builds the Go project
func (b *GoBuilder) Build(ctx context.Context, project Project) error {
	if !project.IsGoProject() {
		return fmt.Errorf("not a Go project")
	}

	cmd := exec.CommandContext(ctx, "go", "build", ".")
	cmd.Dir = project.Path()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if b.logger != nil {
		b.logger.Info("Building project", Field{Key: "project", Value: project.Path()})
	}

	return cmd.Run()
}

// Run runs the Go project
func (b *GoBuilder) Run(ctx context.Context, project Project) error {
	if !project.IsGoProject() {
		return fmt.Errorf("not a Go project")
	}

	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Dir = project.Path()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if b.logger != nil {
		b.logger.Info("Running project", Field{Key: "project", Value: project.Path()})
	}

	return cmd.Run()
}

// Test runs tests for the Go project
func (b *GoBuilder) Test(ctx context.Context, project Project) error {
	if !project.IsGoProject() {
		return fmt.Errorf("not a Go project")
	}

	cmd := exec.CommandContext(ctx, "go", "test", "./...")
	cmd.Dir = project.Path()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if b.logger != nil {
		b.logger.Info("Testing project", Field{Key: "project", Value: project.Path()})
	}

	return cmd.Run()
}

// Clean cleans build artifacts
func (b *GoBuilder) Clean(ctx context.Context, project Project) error {
	if !project.IsGoProject() {
		return fmt.Errorf("not a Go project")
	}

	// Remove common build artifacts
	artifacts := []string{
		filepath.Join(project.Path(), project.Name()),
		filepath.Join(project.Path(), project.Name()+".exe"),
		filepath.Join(project.Path(), "dist"),
	}

	for _, artifact := range artifacts {
		if _, err := os.Stat(artifact); err == nil {
			if err := os.RemoveAll(artifact); err != nil {
				if b.logger != nil {
					b.logger.Warn("Failed to remove artifact",
						Field{Key: "artifact", Value: artifact},
						Field{Key: "error", Value: err.Error()})
				}
			}
		}
	}

	// Run go clean
	cmd := exec.CommandContext(ctx, "go", "clean")
	cmd.Dir = project.Path()

	if b.logger != nil {
		b.logger.Info("Cleaning project", Field{Key: "project", Value: project.Path()})
	}

	return cmd.Run()
}

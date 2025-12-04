// Package filesystem provides file system operations
package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"

	"gox-ide/pkg/core"
)

// OSFileSystem implements core.FileSystem using the OS file system
type OSFileSystem struct{}

// NewOSFileSystem creates a new OS-based file system
func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

// ReadFile reads the contents of a file
func (f *OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes contents to a file
func (f *OSFileSystem) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// ListFiles lists files in a directory
func (f *OSFileSystem) ListFiles(path string) ([]core.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []core.FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		filePath := filepath.Join(path, entry.Name())
		relPath, _ := filepath.Rel(path, filePath)

		fileInfo := core.FileInfo{
			Name:     entry.Name(),
			Path:     filePath,
			RelPath:  relPath,
			IsDir:    entry.IsDir(),
			Size:     info.Size(),
			ModTime:  info.ModTime().Unix(),
			Language: core.GetLanguageForFile(entry.Name()),
		}

		files = append(files, fileInfo)
	}

	return files, nil
}

// WalkDir walks a directory tree
func (f *OSFileSystem) WalkDir(root string, fn func(core.FileInfo) error) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(root, path)

		fileInfo := core.FileInfo{
			Name:     d.Name(),
			Path:     path,
			RelPath:  relPath,
			IsDir:    d.IsDir(),
			Size:     info.Size(),
			ModTime:  info.ModTime().Unix(),
			Language: core.GetLanguageForFile(d.Name()),
		}

		return fn(fileInfo)
	})
}

// Exists checks if a file or directory exists
func (f *OSFileSystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

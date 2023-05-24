package files

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed files/*
var files embed.FS

// ExtractFiles extracts the files from the embedded filesystem to the target directory
func ExtractFiles(targetDir string) error {
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		return err
	}

	return fs.WalkDir(files, "files", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		data, err := files.ReadFile(path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, path)
		err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, os.ModePerm)
	})
}

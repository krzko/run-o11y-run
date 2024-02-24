package files

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed files/*
var files embed.FS

type ServicesConfig struct {
	LocalLogFiles bool   // follow local log files
	LogFilePath   string // follow specified log file
	LogFiles      bool   // true if LocalLogFiles or LogFilePath
}

// ExtractFiles extracts the files from the embedded filesystem to the target directory
func ExtractFiles(targetDir string, svcConfig ServicesConfig) error {
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
		if filepath.Ext(path) == ".tmpl" {
			serviceTemplate, err := template.New("serviceTemplateConfigFile").Parse(string(data))
			if err != nil {
				return err
			}
			var hydratedContent bytes.Buffer
			err = serviceTemplate.Execute(&hydratedContent, svcConfig)
			if err != nil {
				return err
			}
			data = hydratedContent.Bytes()
			path = strings.ReplaceAll(path, ".tmpl", "")
		}

		targetPath := filepath.Join(targetDir, path)
		err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, os.ModePerm)
	})
}

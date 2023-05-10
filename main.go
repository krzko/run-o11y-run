package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

var (
	buildVersion string
	commit       string
	date         string
)

//go:embed files/*
var files embed.FS

func main() {
	ver := fmt.Sprintf("run-o11y-run %v commit: %v date: %v", buildVersion, commit, date)

	if !checkDockerAvailability() {
		fmt.Println("Docker command not found. Please make sure Docker is installed and available in your PATH.")
		os.Exit(1)
	}

	targetDir := "run_o11y_run_files"
	err := extractFiles(targetDir)
	if err != nil {
		fmt.Println("Error extracting files:", err)
		return
	}

	dockerComposeUpFlag := ""
	clean := false
	if len(os.Args) > 1 {
		if os.Args[1] == "-d" {
			dockerComposeUpFlag = "-d"
		} else if os.Args[1] == "-clean" {
			clean = true
		} else if os.Args[1] == "-v" {
			fmt.Println(ver)
			os.Exit(0)
		}
	}

	if clean {
		err = runDockerCompose(filepath.Join(targetDir, "files", "grafana", "stack"), "down", "")
		if err != nil {
			fmt.Println("Error running docker compose down:", err)
			return
		}

		err = removeDirectory(targetDir)
		if err != nil {
			fmt.Println("Error removing target directory:", err)
			return
		}

		os.Exit(0)
	} else {
		err = runDockerCompose(filepath.Join(targetDir, "files", "grafana", "stack"), "up", dockerComposeUpFlag)
		if err != nil {
			fmt.Println("Error running docker compose up:", err)
			return
		}
	}
}

func checkDockerAvailability() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

func extractFiles(targetDir string) error {
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

func removeDirectory(dir string) error {
	return os.RemoveAll(dir)
}

func runDockerCompose(dir string, subcommand string, flag string) error {
	args := []string{"compose", subcommand}
	if flag != "" {
		args = append(args, flag)
	}

	cmd := exec.Command("docker", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("docker compose %s failed: %w", subcommand, err)
	}

	if subcommand == "down" {
		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("docker compose %s failed: %w", subcommand, err)
		}
	} else {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()

		select {
		case <-signalChan:
			fmt.Printf("Received interrupt signal, stopping Docker %s...\n", subcommand)
			_ = cmd.Process.Signal(os.Interrupt)
		case err := <-done:
			if err != nil {
				return fmt.Errorf("docker compose %s failed: %w", subcommand, err)
			}
		}

		err = <-done
		if err != nil {
			return fmt.Errorf("docker compose %s failed: %w", subcommand, err)
		}
	}

	return nil
}

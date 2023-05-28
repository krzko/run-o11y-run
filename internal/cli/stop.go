package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/krzko/run-o11y-run/internal/files"
	"github.com/urfave/cli/v2"
)

// genStopCommand generates the stop command
func genStopCommand() *cli.Command {
	return &cli.Command{
		Name:    "stop",
		Usage:   "Stop run-o11y-run containers",
		Aliases: []string{"t"},
		Action: func(c *cli.Context) error {
			fmt.Printf("🏁 Stopping...\n\n")

			if !checkDockerAvailability() {
				fmt.Println("Docker command not found. Please make sure Docker is installed and available in your PATH.")
				os.Exit(1)
			}

			homeDir := getHomeDir()
			targetDir := filepath.Join(homeDir, ".run-o11y-run")
			err := files.ExtractFiles(targetDir)
			if err != nil {
				fmt.Println("Error extracting files:", err)
				return err
			}

			err = runDockerCompose(filepath.Join(targetDir, "files", "grafana", "run-o11y-run"), "down")
			if err != nil {
				fmt.Println("Error running docker compose down:", err)
				return err
			}

			// os.Exit(0)

			fmt.Printf("\n🏁 Stopped\n\n")
			return nil
		},
	}
}

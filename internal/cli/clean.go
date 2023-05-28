package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/krzko/run-o11y-run/internal/files"
	"github.com/urfave/cli/v2"
)

// genCleanCommand generates the clean command
func genCleanCommand() *cli.Command {
	return &cli.Command{
		Name:    "clean",
		Usage:   "Stop and remove containers, files, networks",
		Aliases: []string{"c"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "legacy",
				Usage: "removes legacy \"stack\" docker resources",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("🧹 Cleaning...\n\n")

			if !checkDockerAvailability() {
				fmt.Println("Docker command not found. Please make sure Docker is installed and available in your PATH.")
				os.Exit(1)
			}

			targetDir := "run_o11y_run_files"
			err := files.ExtractFiles(targetDir)
			if err != nil {
				fmt.Println("Error extracting files:", err)
				return err
			}

			if c.Bool("legacy") {
				err = runDockerCompose(filepath.Join(targetDir, "files", "grafana", "stack"), "down")
				if err != nil {
					fmt.Println("Error running docker compose down:", err)
					return err
				}
			} else {
				err = runDockerCompose(filepath.Join(targetDir, "files", "grafana", "run-o11y-run"), "down")
				if err != nil {
					fmt.Println("Error running docker compose down:", err)
					return err
				}
			}

			err = removeDirectory(targetDir)
			if err != nil {
				fmt.Println("Error removing target directory:", err)
				return err
			}

			fmt.Printf("\n🧹 Cleaned\n\n")
			return nil
		},
	}
}

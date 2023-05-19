package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/krzko/run-o11y-run/internal/files"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func genStartCommand() *cli.Command {
	return &cli.Command{
		Name:    "start",
		Usage:   "Start run-o11y-run containers",
		Aliases: []string{"s"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "registry",
				Aliases: []string{"r"},
				Usage:   "sets a docker registry to pull images from",
				Value:   "registry-1.docker.io",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("✨ Starting...")

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

			// Modify the Docker Compose file with the registry prefix
			dockerComposePath := filepath.Join(targetDir, "files", "grafana", "run-o11y-run", "docker-compose.yaml")
			err = addRegistryPrefix(dockerComposePath, c.String("registry"))
			if err != nil {
				fmt.Println("Error adding registry prefix to Docker Compose file:", err)
				return err
			}

			// Run the Docker Compose up command
			err = runDockerCompose(filepath.Join(targetDir, "files", "grafana", "run-o11y-run"), "up", "")
			if err != nil {
				fmt.Println("Error running docker compose up:", err)
				return err
			}

			fmt.Println("🏁 Stopped...")
			return nil
		},
	}
}

func addRegistryPrefix(filePath, registry string) error {
	// Read the Docker Compose file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read Docker Compose file: %w", err)
	}

	// Unmarshal the YAML into a map
	var composeMap map[interface{}]interface{}
	err = yaml.Unmarshal(data, &composeMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Docker Compose YAML: %w", err)
	}

	// Modify the image field with the registry prefix
	services, ok := composeMap["services"].(map[interface{}]interface{})
	if ok {
		grafana, ok := services["grafana"].(map[interface{}]interface{})
		if ok {
			image, ok := grafana["image"].(string)
			if ok {
				grafana["image"] = fmt.Sprintf("%s/%s", registry, image)
			}
		}
	}

	// Marshal the modified YAML back into bytes
	modifiedData, err := yaml.Marshal(composeMap)
	if err != nil {
		return fmt.Errorf("failed to marshal modified Docker Compose YAML: %w", err)
	}

	// Write the modified YAML back to the file
	err = os.WriteFile(filePath, modifiedData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write modified Docker Compose file: %w", err)
	}

	return nil
}

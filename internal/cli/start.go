package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/krzko/run-o11y-run/internal/files"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
)

// genStartCommand generates the start command
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
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "debug mode",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:  "external-network",
				Usage: "external network mode for docker compose",
				Value: false,
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

			dockerComposePath := filepath.Join(targetDir, "files", "grafana", "run-o11y-run", "docker-compose.yaml")

			// Modify the Docker Compose file with the registry prefix
			if err = addRegistryPrefix(dockerComposePath, c.String("registry")); err != nil {
				fmt.Println("Error adding registry prefix to Docker Compose file:", err)
				return err
			}

			// Modify the Docker Compose to expose named network for other Docker Composes
			if c.Bool("external-network") {
				if err = addExternalNetwork(dockerComposePath); err != nil {
					fmt.Println("Error adding network config to Docker Compose file:", err)
					return err
				}
			}

			if c.Bool("debug") {
				fmt.Println("🐛 Debug mode enabled. Printing Docker Compose file...")
				fmt.Println()
				data, err := os.ReadFile(dockerComposePath)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", data)
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

// addRegistryPrefix adds the registry prefix to the image field of the Docker Compose file
func addRegistryPrefix(filePath, registry string) error {
	// Read the Docker Compose file
	composeMap, err := dockeComposeMap(filePath)
	if err != nil {
		return err
	}

	// Modify the image field with the registry prefix for all services
	services, ok := composeMap["services"].(map[any]any)
	if ok {
		for name, sAny := range services {
			service, ok := sAny.(map[any]any)
			if !ok {
				return fmt.Errorf("unexpected type for service")
			}
			image, ok := service["image"].(string)
			if ok {
				service["image"] = fmt.Sprintf("%s/%s", registry, image)
			} else {
				return fmt.Errorf("error during injecting external registry to service(%s) image defintion", name)
			}
		}
	} else {
		return fmt.Errorf("error during injecting external registry to service image defintion")
	}

	return writeDockerCompose(filePath, composeMap)
}

// addExternalNetwork adds the registry prefix to the image field of the Docker Compose file
func addExternalNetwork(filePath string) error {
	// Read the Docker Compose file
	composeMap, err := dockeComposeMap(filePath)
	if err != nil {
		return err
	}

	// Modify newtorks field with the external network
	services, ok := composeMap["services"].(map[any]any)
	if ok {
		for nAny, sAny := range services {
			service, ok := sAny.(map[any]any)
			if !ok {
				return fmt.Errorf("unexpected type for service")
			}
			service["networks"] = []string{"default"}
			name, _ := nAny.(string)
			// inject o11y network only to otel-collector and tempo service.
			// other services like mini-o11y-stack, grafana, etc. should not be exposed.
			if slices.Contains([]string{"otel-collector", "tempo"}, name) {
				service["networks"] = []string{"o11y", "default"}
			}
		}
	} else {
		return fmt.Errorf("error during injecting external network to service defintion")
	}

	// global networks
	composeMap["networks"] = map[string]map[string]any{
		"default": {
			"driver": "bridge",
		},
		"o11y": {
			"name":       "o11y",
			"attachable": true,
		},
	}
	return writeDockerCompose(filePath, composeMap)
}

// dockeComposeMap returns a map which represents the Docker Compose file
func dockeComposeMap(filePath string) (map[any]any, error) {
	// Read the Docker Compose file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Docker Compose file: %w", err)
	}

	// Unmarshal the YAML into a map
	var composeMap map[any]any
	err = yaml.Unmarshal(data, &composeMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Docker Compose YAML: %w", err)
	}
	return composeMap, nil
}

func writeDockerCompose(filePath string, composeMap map[any]any) error {
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

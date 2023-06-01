package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// genOpenCommand generates the open command
func genPatchComposeCommand() *cli.Command {
	return &cli.Command{
		Name:    "patch-compose",
		Usage:   "Patch customer owner docker-compose with o11y network setup",
		Aliases: []string{"pcm", "patch-compose-manifest"},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "path to customer owner docker-compose.yaml file",
				Required: false,
				Value:    "docker-compose.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			composeFile := c.String("file")
			// Read the Docker Compose file
			composeMap, err := dockeComposeMap(composeFile)
			if err != nil {
				return err
			}
			// patch network
			o11yNetworks := map[any]any{
				"o11y": map[string]any{
					"name":     "o11y",
					"external": true,
					"driver":   "bridge",
				},
			}

			nAny, ok := composeMap["networks"]
			if !ok {
				composeMap["networks"] = o11yNetworks
			} else {
				existingNet, ok := nAny.(map[any]any)
				if !ok {
					return fmt.Errorf("can't patch existing network config: %T", nAny)
				}
				maps.Copy(existingNet, o11yNetworks)
				composeMap["networks"] = existingNet
			}

			// patch service network
			services, ok := composeMap["services"].(map[any]any)
			if ok {
				for _, sAny := range services {
					service, ok := sAny.(map[any]any)
					if !ok {
						return fmt.Errorf("unexpected type for service")
					}
					serviceNetworks, ok := service["networks"].([]any)
					if ok {
						if !slices.Contains(serviceNetworks, "o11y") {
							service["networks"] = append(serviceNetworks, "o11y", "default")
						}
					} else {
						service["networks"] = []string{"o11y", "default"}
					}

					environments, ok := service["environment"].(map[any]any)
					if ok {
						_, ok = environments["OTEL_EXPORTER_OTLP_ENDPOINT"]
						if ok {
							environments["OTEL_EXPORTER_OTLP_ENDPOINT"] = "otel-collector"
							service["environment"] = environments
						}
					} else {
						service["environment"] = map[string]string{"OTEL_EXPORTER_OTLP_ENDPOINT": "otel-collector"}
					}

				}
			} else {
				return fmt.Errorf("error during injecting external network to service definition")
			}
			fmt.Printf("writhing changes to %s", composeFile)
			return writeDockerCompose(composeFile, composeMap)
		},
	}
}

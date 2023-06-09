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
		Name:    "patch",
		Usage:   "patch your Docker Compose file to the o11y network",
		Aliases: []string{"pc", "patch-compose"},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "path to to your Docker Compose file",
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
						if !slices.Contains(serviceNetworks, "default") {
							service["networks"] = append(serviceNetworks, "default")
						}
						if !slices.Contains(serviceNetworks, "o11y") {
							service["networks"] = append(serviceNetworks, "o11y")
						}
					} else {
						service["networks"] = []string{"default", "o11y"}
					}

					environments, ok := service["environment"].(map[any]any)
					const otelCollectorEndpoint = "otel-collector:4317"
					const zipkinEndpoint = "tempo:9411"
					if ok {
						environments["OTEL_EXPORTER_OTLP_ENDPOINT"] = otelCollectorEndpoint
						environments["OTEL_EXPORTER_ZIPKIN_ENDPOINT"] = zipkinEndpoint
						service["environment"] = environments
					} else {
						service["environment"] = map[string]string{
							"OTEL_EXPORTER_OTLP_ENDPOINT":   otelCollectorEndpoint,
							"OTEL_EXPORTER_ZIPKIN_ENDPOINT": zipkinEndpoint,
						}
					}
				}
			} else {
				return fmt.Errorf("error during injecting external network to service definition")
			}
			fmt.Printf("writhing changes to %s\n", composeFile)
			return writeDockerCompose(composeFile, composeMap)
		},
	}
}

package cli

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

// genOpenCommand generates the open command
func genOpenCommand() *cli.Command {
	return &cli.Command{
		Name:    "open",
		Usage:   "Open a web page for a specific service",
		Aliases: []string{"o"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "service",
				Aliases:  []string{"s"},
				Usage:    "specify the service to open: loki, tempo, prometheus, prometheus-direct or pyroscope-direct",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			service := c.String("service")
			url := ""

			switch service {
			case "loki":
				url = "http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D"
			case "tempo":
				url = "http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22tempo%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22tempo%22,%22uid%22:%22tempo%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D"
			case "prometheus":
				url = "http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22prometheus%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22prometheus%22,%22uid%22:%22prometheus%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D"
			case "prometheus-direct":
				url = "http://localhost:9090/"
			case "pyroscope-direct":
				url = "http://localhost:4040/"
			default:
				return fmt.Errorf("unsupported service")
			}

			err := browser.OpenURL(url)
			if err != nil {
				return err
			}

			fmt.Printf("\n 🌏 Opened %s in your default web browser\n\n", service)
			return nil
		},
	}
}

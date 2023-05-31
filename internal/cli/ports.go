package cli

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

// genPortsCommand generates the ports command
func genPortsCommand() *cli.Command {
	return &cli.Command{
		Name:    "ports",
		Usage:   "List the available ports used by run-o11y-run",
		Aliases: []string{"p"},
		Action: func(c *cli.Context) error {
			data := [][]string{
				{"3000/tcp", "Grafana"},
				{"3100/tcp", "Loki"},
				{"4040/tcp", "Pyroscope"},
				{"4317/tcp", "OTLP (gRPC)"},
				{"4318/tcp", "OTLP (HTTP)"},
				{"8094/tcp", "Syslog (RFC3164)"},
				{"9090/tcp", "Prometheus Direct"},
				{"9411/tcp", "Zipkin"},
				{"14268/tcp", "Jaeger"},
			}

			fmt.Println()
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Port", "Service"})
			table.SetAutoFormatHeaders(true)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold})
			table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
				tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor})
			table.AppendBulk(data)
			table.Render()
			fmt.Println()

			return nil
		},
	}
}

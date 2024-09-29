package cli

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

func initBanner(c *cli.Context) error {
	banner := `
_____ _____ _____     _____ ___   ___   __ __     _____ _____ _____
| __  |  |  |   | |___|     |_  | |_  | |  |  |___| __  |  |  |   | |
|    -|  |  | | | |___|  |  |_| |_ _| |_|_   _|___|    -|  |  | | | |
|__|__|_____|_|___|   |_____|_____|_____| |_|     |__|__|_____|_|___|
`
	fmt.Println(banner)

	return nil
}

func New(version, commit, date string) *cli.App {
	c := []color.Attribute{color.FgRed, color.FgGreen, color.FgYellow, color.FgMagenta, color.FgCyan, color.FgWhite, color.FgHiRed, color.FgHiGreen, color.FgHiYellow, color.FgHiBlue, color.FgHiMagenta, color.FgHiCyan, color.FgHiWhite}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c), func(i, j int) { c[i], c[j] = c[j], c[i] })

	colors := make(map[int]func(...any) string)
	for i, attr := range c {
		colors[i] = color.New(attr).SprintFunc()
	}

	name := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s", colors[0]("r"), colors[1]("u"), colors[2]("n"), colors[3]("-"), colors[4]("o"), colors[5]("1"), colors[6]("1"), colors[7]("y"), colors[8]("-"), colors[9]("r"), colors[10]("u"), colors[11]("n"))

	flags := getGlobalFlags()

	var v string
	if version == "" {
		v = "develop"
	} else {
		v = fmt.Sprintf("v%v-%v (%v)", version, commit, date)
	}

	app := &cli.App{
		Name:    name,
		Usage:   "A single-binary 🌯 wrapper around `docker compose` with embedded configurations to effortlessly run your local observability stack ",
		Version: v,
		Flags:   flags,
		Commands: []*cli.Command{
			genCleanCommand(),
			genOpenCommand(),
			genPortsCommand(),
			genStartCommand(),
			genStopCommand(),
			genPatchComposeCommand(),
		},
		Before: initBanner,
	}

	app.EnableBashCompletion = true

	return app
}

// checkDockerAvailability checks if Docker is available in the PATH
func checkDockerAvailability() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

// getHomeDir returns the user's home directory
func getHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("Failed to get home directory:", err)
		os.Exit(1)
	}
	return home
}

// removeDirectory removes a directory and all its contents
func removeDirectory(dir string) error {
	return os.RemoveAll(dir)
}

// runDockerCompose runs a docker compose command
func runDockerCompose(dir, subcommand string, flags ...string) error {
	args := []string{"compose", subcommand}
	args = append(args, flags...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("docker compose %s failed to start: %w", subcommand, err)
	}

	if slices.Contains(args, "--detach") {
		fmt.Println("🚀 Started in detached mode")
		return nil
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-sigChan:
		fmt.Printf("\nReceived interrupt signal, stopping Docker %s...\n", subcommand)
		cancel()

		// Wait for graceful shutdown with a timeout
		select {
		case <-time.After(10 * time.Second):
			fmt.Printf("Forcefully stopping Docker %s...\n", subcommand)
			_ = cmd.Process.Kill()
		case <-done:
			fmt.Printf("Docker %s stopped gracefully.\n", subcommand)
		}

		return nil
	case err := <-done:
		if err != nil {
			return fmt.Errorf("docker compose %s failed: %w", subcommand, err)
		}
		return nil
	}
}

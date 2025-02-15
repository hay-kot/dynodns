package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hay-kot/dynodns/app/commands"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

type DNSRecord struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
	TTL     int    `json:"ttl"`
}

func main() {
	ctrl := commands.New()

	app := &cli.App{
		Name:    "dynodns",
		Usage:   "client for setting up dynamic DNS",
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "log-level",
				Usage:    "log level (debug, info, warn, error, fatal, panic)",
				Value:    "info",
				EnvVars:  []string{"LOG_LEVEL"},
				Category: "Options",
			},
			&cli.StringFlag{
				Name:     "log-style",
				Usage:    "set the log style to console or json",
				Value:    "console",
				EnvVars:  []string{"LOG_STYLE", "LOG_FORMAT"},
				Category: "Options",
			},
			&cli.IntFlag{
				Name:        "interval",
				Usage:       "interval in seconds to check for ip changes",
				Value:       300,
				Destination: &ctrl.Interval,
				EnvVars:     []string{"INTERVAL"},
				Category:    "Options",
			},
			&cli.StringFlag{
				Name:        "ping-url",
				Usage:       "Healthcheck Ping URL",
				Destination: &ctrl.PingURL,
				EnvVars:     []string{"PING_URL"},
				Category:    "Options",
			},
			&cli.StringFlag{
				Name:        "porkbun.subdomain",
				Usage:       "porkbun subdomain to update",
				Value:       "dns",
				Destination: &ctrl.PorkBunSubDomain,
				EnvVars:     []string{"PORKBUN_SUBDOMAIN"},
				Category:    "Porkbun",
			},
			&cli.StringFlag{
				Name:        "porkbun.domain",
				Usage:       "porkbun domain to update",
				Required:    true,
				Destination: &ctrl.PorkBunDomain,
				EnvVars:     []string{"PORKBUN_DOMAIN"},
				Category:    "Porkbun",
			},
			&cli.StringFlag{
				Name:        "porkbun.endpoint",
				Usage:       "porkbun api endpoint",
				Value:       "https://api.porkbun.com/api/json/v3",
				Destination: &ctrl.PorkBunEndpoint,
				EnvVars:     []string{"PORKBUN_API_ENDPOINT"},
				Category:    "Porkbun",
			},
			&cli.StringFlag{
				Name:        "porkbun.key",
				Usage:       "porkbun api key",
				Required:    true,
				Destination: &ctrl.PorkBunKey,
				EnvVars:     []string{"PORKBUN_API_KEY"},
				Category:    "Porkbun",
			},
			&cli.StringFlag{
				Name:        "porkbun.secret",
				Usage:       "porkbun api secret",
				Required:    true,
				Destination: &ctrl.PorkBunSecret,
				EnvVars:     []string{"PORKBUN_API_SECRET"},
				Category:    "Porkbun",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.String("log-style") != "json" {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			}

			level, err := zerolog.ParseLevel(ctx.String("log-level"))
			if err != nil {
				return err
			}

			log.Logger = log.Level(level)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "runs the client",
				Action: func(cliCtx *cli.Context) error {
					ctx, cancel := context.WithCancel(context.Background())

					// Start the task in a separate goroutine
					finished, err := ctrl.Run(ctx)
					if err != nil {
						cancel()
						return err
					}

					// Listen for termination signals from the terminal
					sigCh := make(chan os.Signal, 1)
					signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
					go func() {
						<-sigCh
						log.Info().Msg("cancellation request received. terminating...")
						cancel()
					}()

					// Wait for cancellation
					<-ctx.Done()
					<-finished
					return nil
				},
			},
			{
				Name:  "test-ip",
				Usage: "test external IP finder",
				Action: func(ctx *cli.Context) error {
					ip, err := ctrl.GetPublicIP(context.Background())
					if err != nil {
						return err
					}

					log.Info().Str("ip", ip).Msg("found external ip")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run dynodns")
	}
}

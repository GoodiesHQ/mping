package main

import (
	"context"
	_ "embed"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goodieshq/mping/pkg/pinger"
	"github.com/goodieshq/mping/pkg/privcheck"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var Version string = "dev"

var app *cli.Command

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	}).Level(zerolog.ErrorLevel)

	app = &cli.Command{
		Name:        "mping",
		Usage:       "Multi-Target Ping Utility",
		Description: "mping is a command-line utility that allows users to ping multiple targets simultaneously.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"V"},
				Usage:   "Show the version number and exit",
			},
			&cli.BoolFlag{
				Name:    "ipv4",
				Aliases: []string{"4"},
				Usage:   "Use IPv4 for name resolution",
			},
			&cli.BoolFlag{
				Name:    "ipv6",
				Aliases: []string{"6"},
				Usage:   "Use IPv6 for name resolution",
			},
			&cli.Uint32Flag{
				Name:    "count",
				Aliases: []string{"c"},
				Usage:   "Number of echo requests to send to each target (0 = unlimited)",
				Value:   0,
			},
			&cli.Float64Flag{
				Name:    "interval",
				Aliases: []string{"i"},
				Usage:   "Interval (in seconds) between sending each packet (minimum 0.01)",
				Value:   1.0,
			},
			&cli.Float64Flag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "Timeout (in seconds) to wait for each reply (minimum 0.01)",
				Value:   1.0,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable more verbose logging for debug output",
			},
		},
		UsageText: "mping [options] target1[=Label1] target2[=Label2] ...",
		Action:    action,
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	args := make([]string, 0, len(os.Args))

	for _, arg := range os.Args {
		switch arg {
		case "-4":
			args = append(args, "--ipv4")
		case "-6":
			args = append(args, "--ipv6")
		default:
			args = append(args, arg)
		}
	}

	err := app.Run(ctx, args)
	if err != nil {
		if err == context.Canceled {
			return
		}
		log.Fatal().Err(err).Msg("Application error")
	}
}

func action(ctx context.Context, c *cli.Command) error {
	if c.Bool("version") {
		println(Version)
		return nil
	}

	// Check for administrative privileges
	if !privcheck.HasAdmin() {
		return cli.Exit("Administrative privileges are required to use mping.", 1)
	}

	// Set verbose logging if specified
	if c.Bool("verbose") {
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
		log.Debug().Msg("Verbose logging enabled")
	}

	// Check if timeout is less than intervalFloat
	intervalFloat := c.Float64("interval")
	timeoutFloat := c.Float64("timeout")
	ipv4 := c.Bool("ipv4")
	ipv6 := c.Bool("ipv6")

	// clamp timeout to interval, can't ping for longer than the interval between pings
	if timeoutFloat > intervalFloat {
		log.Warn().Msgf("Timeout (%.2f seconds) is greater than interval (%.2f seconds). Setting timeout to interval value.", timeoutFloat, intervalFloat)
		timeoutFloat = intervalFloat
	}

	if intervalFloat < 0.01 {
		return cli.Exit("Interval must be at least 0.01 seconds", 1)
	}

	if timeoutFloat < 0.01 {
		return cli.Exit("Timeout must be at least 0.01 seconds", 1)
	}

	interval := time.Microsecond * time.Duration(intervalFloat*1e6)
	timeout := time.Microsecond * time.Duration(timeoutFloat*1e6)

	if !ipv4 && !ipv6 {
		log.Debug().Msg("No IP version specified, defaulting to IPv4")
		ipv4 = true
	}

	if ipv4 && ipv6 {
		return cli.Exit("Cannot specify both --ipv4 and --ipv6 flags", 1)
	}

	var ipVersion pinger.IPVersion = pinger.IPV0

	if ipv4 {
		ipVersion = pinger.IPV4
	} else if ipv6 {
		ipVersion = pinger.IPV6
	}

	// Acquire the targets
	targets := parseTargets(c.Args().Slice()...)
	if len(targets) == 0 {
		return cli.Exit("No targets specified", 1)
	}

	resolved := pinger.ResolveTargets(ctx, targets, ipVersion)
	for _, rt := range resolved {
		if rt.IP == nil {
			log.Warn().Msgf("Could not resolve target: %s", rt.Label)
			continue
		}
		log.Debug().Msgf("Resolved Target: %s (%s)", rt.IP.String(), rt.Label)
	}

	count := c.Uint32("count")
	showRTT := true
	if showRTT {
		log.Debug().Msg("Latency display enabled")
	}

	log.Debug().Uint32("count", count).Msgf("Pinging %d targets", len(targets))
	pinger := pinger.NewPinger(pinger.PingerOptions{
		ResolvedTargets: resolved,
		Count:           count,
		Interval:        interval,
		Timeout:         timeout,
		ShowRTT:         showRTT,
	})
	return pinger.Start(ctx)
}

func parseTargets(args ...string) []*pinger.Target {
	var targets []*pinger.Target

	for i, arg := range args {
		target := pinger.ParseTarget(arg)
		if target == nil {
			log.Warn().Msgf("Skipping empty target #%d", i+1)
			continue
		}
		targets = append(targets, target)
	}

	return targets
}

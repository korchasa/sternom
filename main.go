package main

import (
	"context"
	"fmt"
	"github.com/korchasa/sternom/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	Version string
)

type Options struct {
	nomadAddr  string
	follow     bool
	onlyStdout bool
	onlyStderr bool
	nbytes     int64
	color      string
	version    bool
}

var opts = &Options{
	nomadAddr:  "NOMAD_ADDR",
	follow:     false,
	onlyStdout: false,
	onlyStderr: false,
	nbytes:     -1,
	color:      "auto",
}

func main() {
	log.SetFlags(0)

	cmd := &cobra.Command{}
	cmd.Use = "sternom job-or-alloc-prefix"
	cmd.Short = "Tail multiple jobs and allocations from Nomad"

	cmd.Flags().StringVarP(&opts.nomadAddr, "address", "a", opts.nomadAddr, "The address of the Nomad server. Overrides the NOMAD_ADDR environment variable if set.")
	cmd.Flags().BoolVarP(&opts.follow, "follow", "f", opts.follow, "Whether the logs should be followed")
	cmd.Flags().BoolVar(&opts.onlyStdout, "stdout", opts.onlyStdout, "Show only stdout log")
	cmd.Flags().BoolVar(&opts.onlyStderr, "stderr", opts.onlyStderr, "Show only stderr log")
	cmd.Flags().Int64VarP(&opts.nbytes, "tail", "t", opts.nbytes, "The number of bytes from the end of the logs to show. Defaults to -1, showing all logs.")
	cmd.Flags().StringVar(&opts.color, "color", opts.color, "Color output. Can be 'always', 'never', or 'auto'")
	cmd.Flags().BoolVarP(&opts.version, "version", "v", opts.version, "Print the version and exit")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if opts.version {
			fmt.Printf("sternom version %s\n", Version)
			return nil
		}

		narg := len(args)
		if narg != 1 {
			return cmd.Help()
		}
		config, err := parseConfig(args)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err = pkg.Run(ctx, config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return nil
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func parseConfig(args []string) (*pkg.Config, error) {
	prefix := args[0]

	opts.nomadAddr = os.Getenv("NOMAD_ADDR")

	cf := opts.color
	if cf == "always" {
		color.NoColor = false
	} else if cf == "never" {
		color.NoColor = true
	} else if cf != "auto" {
		return nil, errors.New("color should be one of 'always', 'never', or 'auto'")
	}

	showStdout, showStderr := true, true
	if opts.onlyStdout && opts.onlyStderr {
		return nil, errors.New("can't combine stdout and stderr flags")
	} else if opts.onlyStdout {
		showStderr = false
	} else if opts.onlyStderr {
		showStdout = false
	}

	return &pkg.Config{
		JobsOrAllocPrefix: prefix,
		NomadAddress:      opts.nomadAddr,
		Follow:            opts.follow,
		ShowStdout:        showStdout,
		ShowStderr:        showStderr,
		TailBytes:         opts.nbytes,
	}, nil
}

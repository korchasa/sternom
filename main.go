package main

import (
	"context"
	"fmt"
	"github.com/korchasa/sternom/pkg"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	Version string
	opts    = &pkg.Options{
		NomadAddr:  "NOMAD_ADDR",
		Follow:     false,
		TailBytes:  -1,
		New:        false,
		OnlyStdout: false,
		OnlyStderr: false,
		Color:      "auto",
	}
)

func main() {
	log.SetFlags(0)

	cmd := &cobra.Command{}
	cmd.Use = "sternom job-or-alloc-prefix"
	cmd.Short = "Tail multiple jobs and allocations from Nomad"

	cmd.Args = cobra.ExactArgs(1)
	cmd.Flags().StringVarP(&opts.NomadAddr, "address", "a", opts.NomadAddr, "The address of the Nomad server. Overrides the NOMAD_ADDR environment variable if set.")
	cmd.Flags().BoolVarP(&opts.Follow, "follow", "f", opts.Follow, "Whether the logs should be followed")
	cmd.Flags().Int64VarP(&opts.TailBytes, "tail", "t", opts.TailBytes, "The number of bytes from the end of the logs to show. Defaults to -1, showing all logs.")
	cmd.Flags().BoolVarP(&opts.New, "new", "n", opts.New, "Shorthand for --follow and --tail 0")
	cmd.Flags().BoolVar(&opts.OnlyStdout, "stdout", opts.OnlyStdout, "Show only stdout log")
	cmd.Flags().BoolVar(&opts.OnlyStderr, "stderr", opts.OnlyStderr, "Show only stderr log")
	cmd.Flags().StringVar(&opts.Color, "color", opts.Color, "Color output. Can be 'always', 'never', or 'auto'")
	cmd.Flags().BoolVarP(&opts.Version, "version", "v", opts.Version, "Print the version and exit")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if opts.Version {
			fmt.Printf("sternom version %s\n", Version)
			return nil
		}

		config, err := pkg.ParseCLIArguments(args[0], opts)
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"korchasa/sternom/pkg"
	"log"
	"os"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/fatih/color"
)

const version = "1.11.0"

type Options struct {
	nomadAddr  string
	timestamps bool
	tail       int64
	color      string
	version    bool
	template   string
	output     string
}

var opts = &Options{
	tail:     -1,
	color:    "auto",
	template: "",
	output:   "default",
}

func main() {
	cmd := &cobra.Command{}
	cmd.Use = "sternom job-or-alloc-prefix"
	cmd.Short = "Tail multiple jobs and allocations from Nomad"

	cmd.Flags().StringVarP(&opts.nomadAddr, "address", "a", os.Getenv("NOMAD_ADDR"), "The address of the Nomad server. Overrides the NOMAD_ADDR environment variable if set.")
	cmd.Flags().BoolVarP(&opts.timestamps, "timestamps", "p", opts.timestamps, "Print timestamps")
	cmd.Flags().Int64Var(&opts.tail, "tail", opts.tail, "The number of lines from the end of the logs to show. Defaults to -1, showing all logs.")
	cmd.Flags().StringVar(&opts.color, "color", opts.color, "Color output. Can be 'always', 'never', or 'auto'")
	cmd.Flags().BoolVarP(&opts.version, "version", "v", opts.version, "Print the version and exit")
	cmd.Flags().StringVar(&opts.template, "template", opts.template, "Template to use for log lines, leave empty to use --output flag")
	cmd.Flags().StringVarP(&opts.output, "output", "o", opts.output, "Specify predefined template. Currently support: [default, raw, json]")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if opts.version {
			fmt.Printf("sternom version %s\n", version)
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

		err = pkg.RunApp(ctx, config)
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

	var tailLines *int64
	if opts.tail != -1 {
		tailLines = &opts.tail
	}

	colorFlag := opts.color
	if colorFlag == "always" {
		color.NoColor = false
	} else if colorFlag == "never" {
		color.NoColor = true
	} else if colorFlag != "auto" {
		return nil, errors.New("color should be one of 'always', 'never', or 'auto'")
	}

	t := opts.template
	if t == "" {
		switch opts.output {
		case "default":
			if color.NoColor {
				t = "{{.PodName}} {{.ContainerName}} {{.Message}}"
			} else {
				t = "{{color .PodColor .PodName}} {{color .ContainerColor .ContainerName}} {{.Message}}"
			}
		case "raw":
			t = "{{.Message}}"
		case "json":
			t = "{{json .}}\n"
		}
	}

	funs := map[string]interface{}{
		"json": func(in interface{}) (string, error) {
			b, err := json.Marshal(in)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
		"color": func(color color.Color, text string) string {
			return color.SprintFunc()(text)
		},
	}
	tpl, err := template.New("log").Funcs(funs).Parse(t)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse template")
	}

	return &pkg.Config{
		JobsOrAllocPrefix: prefix,
		NomadAddress:      opts.nomadAddr,
		Timestamps:        opts.timestamps,
		TailLines:         tailLines,
		Template:          tpl,
	}, nil
}

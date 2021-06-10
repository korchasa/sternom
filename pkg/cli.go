package pkg

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"os"
)

type Options struct {
	NomadAddr  string
	Follow     bool
	TailBytes  int64
	New        bool
	OnlyStdout bool
	OnlyStderr bool
	TaskName   string
	Color      string
	Version    bool
}

func ParseCLIArguments(prefix string, opts *Options) (*Config, error) {
	opts.NomadAddr = os.Getenv("NOMAD_ADDR")

	cf := opts.Color
	if cf == "always" {
		color.NoColor = false
	} else if cf == "never" {
		color.NoColor = true
	} else if cf != "auto" {
		return nil, errors.New("Color should be one of 'always', 'never', or 'auto'")
	}

	showStdout, showStderr := true, true
	if opts.OnlyStdout && opts.OnlyStderr {
		return nil, errors.New("can't combine stdout and stderr flags")
	} else if opts.OnlyStdout {
		showStderr = false
	} else if opts.OnlyStderr {
		showStdout = false
	}

	if opts.New {
		opts.Follow = true
		opts.TailBytes = 0
	}

	return &Config{
		JobsOrAllocPrefix: prefix,
		NomadAddress:      opts.NomadAddr,
		Follow:            opts.Follow,
		ShowStdout:        showStdout,
		ShowStderr:        showStderr,
		TaskName:          opts.TaskName,
		TailBytes:         opts.TailBytes,
	}, nil
}

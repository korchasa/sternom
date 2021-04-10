package pkg

import (
	"text/template"
)

// Config contains the config (Where is my rock?)
type Config struct {
	JobsOrAllocPrefix string
	NomadAddress      string
	Timestamps        bool
	TailLines         *int64
	Template          *template.Template
}

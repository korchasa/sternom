package pkg

// Config contains the config (Where is my rock?)
type Config struct {
	JobsOrAllocPrefix string
	NomadAddress      string
	Timestamps        bool
	Follow            bool
	ShowStdout        bool
	ShowStderr        bool
	TailBytes         int64
}

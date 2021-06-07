# Sternom

## Installation

```bash
go get github.com/korchasa/sternom
sternom --help
```

## Usage

```
Tail multiple jobs and allocations from Nomad

Usage:
  sternom job-or-alloc-prefix [flags]

Flags:
  -a, --address string   The address of the Nomad server. Overrides the NOMAD_ADDR environment variable if set. (default "NOMAD_ADDR")
      --color string     Color output. Can be 'always', 'never', or 'auto' (default "auto")
  -f, --follow           Whether the logs should be followed
  -h, --help             help for sternom
  -n, --new              Shorthand for --follow and --tail 0
      --stderr           Show only stderr log
      --stdout           Show only stdout log
  -t, --tail int         The number of bytes from the end of the logs to show. Defaults to -1, showing all logs. (default -1)
  -v, --version          Print the version and exit
```

![old records](./docs/old.png)

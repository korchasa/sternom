# Sternom

## Installation

```bash
go get github.com/korchasa/sternom
sternom --help
```

### Docker

Images: https://github.com/korchasa/sternom/pkgs/container/sternom

## Usage

```
Tail multiple jobs and allocations from Nomad

Usage:
  sternom job-or-alloc-prefix [flags]

Flags:
  -a, --address string      The address of the Nomad server. Overrides the NOMAD_ADDR environment variable if set. (default "NOMAD_ADDR")
      --color string        Color output. Can be 'always', 'never', or 'auto' (default "auto")
  -e, --exclude -e a -e b   Exclude log records by pattern. Multiple filters: -e a -e b or `-e a,b`
  -i, --filter -i a -i b    Filter log records by pattern. Multiple filters: -i a -i b or `-i a,b`
  -f, --follow              Whether the logs should be followed
  -h, --help                help for sternom
  -n, --new                 Shorthand for --follow and --tail 0
      --stderr              Show only stderr log
      --stdout              Show only stdout log
  -t, --tail int            The number of bytes from the end of the logs to show. Defaults to -1, showing all logs. (default -1)
      --task string         Show logs only for one task
  -v, --version             Print the version and exit
```

![old records](./docs/old.png)

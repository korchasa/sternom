# Sternom

## Installation

```bash
go get github.com/korchasa/sternom
sternom --help
```

## Usage

Tail old records by job prefix:

```bash
sternom ovc
```

![old records](./docs/old.png)

Skip old records and follow logs by allocation id:

```bash
sternom 2b79b1e1 --follow --tail 0
```

![old records](./docs/follow.png)

Watch for 404 errors from `95.217.226.19`, but skip `Build/OPM1.171019.026` user-agent:

```bash
sternom mp3 --follow | grep "nginx" | grep " 404 " | grep -v "Build/OPM1.171019.026"| grep "95.217.226.19"
```

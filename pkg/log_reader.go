package pkg

import (
	"github.com/hashicorp/nomad/api"
	"strings"
)

func LogReader(prefix string, in <-chan *api.StreamFrame, out chan<- string) {
	for data := range in {
		if data == nil {
			continue
		}
		for _, s := range strings.Split(string(data.Data), "\n") {
			if len(s) > 0 {
				out <- prefix + s
			}
		}
	}
}

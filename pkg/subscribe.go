package pkg

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/nomad/api"
	"hash/fnv"
	"strings"
)

func subscribe(client *api.Client, subs []*Subscription) {
	for _, sub := range subs {
		fmt.Printf("Subscription %s\n", sub.String())
		al := &api.Allocation{ID: sub.Alloc, NodeID: sub.Node}
		dataCh, errorsCh := client.AllocFS().Logs(al, true, sub.Task, "stderr", "end", 0, nil, nil)
		go func(sub *Subscription) {
			jobColor, allocColor, taskColor := determineColors(sub.Job, sub.Alloc, sub.Task)
			for {
				select {
				case data := <- dataCh:
					for _, s := range strings.Split(string(data.Data), "\n") {
						_, _ = jobColor.Printf("%s", sub.Job)
						fmt.Printf(":")
						_, _ = allocColor.Printf("%s", sub.AllocShort)
						fmt.Printf("[")
						_, _ = taskColor.Printf("%s", sub.Task)
						fmt.Printf("] %s\n", s)
					}
				case err := <- errorsCh:
					fmt.Printf("Error from %s: %v", sub.Alloc, err)
				}
			}
		}(sub)
	}
}

var jobColors = []*color.Color{
	color.New(color.FgHiCyan),
	color.New(color.FgHiGreen),
	color.New(color.FgHiMagenta),
	color.New(color.FgHiYellow),
	color.New(color.FgHiBlue),
	color.New(color.FgHiRed),
}

var otherColors = []*color.Color{
	color.New(color.FgCyan),
	color.New(color.FgGreen),
	color.New(color.FgMagenta),
	color.New(color.FgYellow),
	color.New(color.FgBlue),
	color.New(color.FgRed),
}

func determineColors(job, alloc, task string) (jobColor, allocColor, taskColor *color.Color) {
	hash := fnv.New32()
	_, _ = hash.Write([]byte(job))
	jobColor = jobColors[hash.Sum32() % uint32(len(jobColors))]

	hash = fnv.New32()
	_, _ = hash.Write([]byte(alloc))
	allocColor = otherColors[hash.Sum32() % uint32(len(otherColors))]

	hash = fnv.New32()
	_, _ = hash.Write([]byte(task))
	taskColor = otherColors[hash.Sum32() % uint32(len(otherColors))]

	return
}

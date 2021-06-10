package pkg

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/nomad/api"
	"hash/fnv"
	"log"
	"sync"
)

func Subscriber(client *api.Client, conf *Config, subs <-chan Subscription, out chan<- string, wg *sync.WaitGroup) {
	for sub := range subs {
		jobColor, allocColor, taskColor := determineColors(sub.Job, sub.Alloc, sub.Task)
		name := fmt.Sprintf("%s:%s[%s]",
			jobColor.Sprintf("%s", sub.Job),
			allocColor.Sprintf("%s", sub.AllocShort),
			taskColor.Sprintf("%s", sub.Task))
		log.Printf("+ %s\n", name)

		from := "end"
		offset := conf.TailBytes
		if conf.TailBytes == 0 {
			offset = 1
		} else if conf.TailBytes == -1 {
			from = "start"
			offset = 0
		}

		al := &api.Allocation{ID: sub.Alloc, NodeID: sub.Node}

		if conf.ShowStdout {
			wg.Add(1)
			go func(sub Subscription) {
				defer wg.Done()
				stdoutCh, _ := client.AllocFS().Logs(al, conf.Follow, sub.Task, "stdout", from, offset, nil, nil)
				LogReader(fmt.Sprintf("%s  ", name), stdoutCh, out)
			}(sub)
		}

		if conf.ShowStderr {
			wg.Add(1)
			go func(sub Subscription) {
				defer wg.Done()
				stderrCh, _ := client.AllocFS().Logs(al, conf.Follow, sub.Task, "stderr", from, offset, nil, nil)
				LogReader(fmt.Sprintf("%s%s ", name, color.New(color.BgRed).Sprint("!")), stderrCh, out)
			}(sub)
		}
	}
}

var jobColors = []*color.Color{
	color.New(color.FgHiCyan),
	color.New(color.FgGreen),
	color.New(color.FgMagenta),
	color.New(color.FgYellow),
	color.New(color.FgBlue),
	color.New(color.FgRed),
}

var otherColors = []*color.Color{
	color.New(color.FgHiCyan),
	color.New(color.FgGreen),
	color.New(color.FgMagenta),
	color.New(color.FgYellow),
	color.New(color.FgBlue),
	color.New(color.FgRed),
}

func determineColors(job, alloc, task string) (jobColor, allocColor, taskColor *color.Color) {
	hash := fnv.New32()
	_, _ = hash.Write([]byte(job))
	jobColor = jobColors[hash.Sum32()%uint32(len(jobColors))]

	hash = fnv.New32()
	_, _ = hash.Write([]byte(alloc))
	allocColor = otherColors[hash.Sum32()%uint32(len(otherColors))]

	hash = fnv.New32()
	_, _ = hash.Write([]byte(task))
	taskColor = otherColors[hash.Sum32()%uint32(len(otherColors))]

	return
}

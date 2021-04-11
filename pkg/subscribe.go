package pkg

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
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
				LogReader(name+"  ", stdoutCh, out)
			}(sub)
		}

		if conf.ShowStderr {
			wg.Add(1)
			go func(sub Subscription) {
				defer wg.Done()
				stderrCh, _ := client.AllocFS().Logs(al, conf.Follow, sub.Task, "stderr", from, offset, nil, nil)
				LogReader(name+"! ", stderrCh, out)
			}(sub)
		}
	}
}

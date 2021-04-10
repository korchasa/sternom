package pkg

import (
	"context"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func RunApp(_ context.Context, conf *Config) error {
	client, err := api.NewClient(&api.Config{Address: conf.NomadAddress})

	if err != nil {
		log.Fatalf("Can't init client: %s", err)
	}

	var subs []*Subscription
	subs, err = findJobSubscription(client, conf.JobsOrAllocPrefix)
	if err != nil {
		log.Fatalf("Can't find Job subscriptions: %s", err)
	}
	log.Printf("%d Job subsriptions\n", len(subs))
	if len(subs) == 0 {
		subs, err = findAllocSubscription(client, conf.JobsOrAllocPrefix)
		if err != nil {
			log.Fatalf("Can't find allocations subscriptions: %s", err)
		}
		log.Printf("%d Alloc subsriptions\n", len(subs))
	}

	cancelCh := make(chan os.Signal)
	signal.Notify(cancelCh, os.Interrupt, syscall.SIGINT)
	go func() {
		<-cancelCh
		log.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()

	outputCh := make(chan string, 10)
	go func() {
		for {
			select {
			case str := <-outputCh:
				fmt.Println(str)
			}
		}
	}()

	wg := &sync.WaitGroup{}
	subscribe(client, conf, subs, outputCh, wg)
	wg.Wait()

	return nil
}

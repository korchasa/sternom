package pkg

import (
	"context"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RunApp(ctx context.Context, conf *Config) error  {
	client, err := api.NewClient(&api.Config{Address: conf.NomadAddress})

	if err != nil {
		log.Fatalf("Can't init client: %s", err)
	}

	var subs []*Subscription
	subs, err = findJobSubscription(client, conf.JobsOrAllocPrefix)
	if err != nil {
		log.Fatalf("Can't find Job subscriptions: %s", err)
	}
	fmt.Printf("%d Job subsriptions\n", len(subs))
	if len(subs) == 0 {
		subs, err = findAllocSubscription(client, conf.JobsOrAllocPrefix)
		if err != nil {
			log.Fatalf("Can't find allocations subscriptions: %s", err)
		}
		fmt.Printf("%d Alloc subsriptions\n", len(subs))
	}

	subscribe(client, subs)

	cancelCh := make(chan os.Signal)
	signal.Notify(cancelCh, os.Interrupt, syscall.SIGTERM)
	<-cancelCh
	return nil
}

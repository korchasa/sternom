package pkg

import (
	"context"
	"github.com/hashicorp/nomad/api"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Run(_ context.Context, conf *Config) error {
	client, err := api.NewClient(&api.Config{Address: conf.NomadAddress})
	if err != nil {
		log.Fatalf("Can't init client: %s", err)
	}

	cancelCh := make(chan os.Signal)
	signal.Notify(cancelCh, os.Interrupt, syscall.SIGINT)
	go func() {
		<-cancelCh
		log.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()

	outputCh := make(chan string, 10)
	go PrintLogRecord(outputCh)

	subsCh := make(chan Subscription)
	go SubscriptionFinder(client, subsCh, conf.JobsOrAllocPrefix)

	wg := &sync.WaitGroup{}
	Subscriber(client, conf, subsCh, outputCh, wg)
	wg.Wait()

	return nil
}

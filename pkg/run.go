package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/hashicorp/nomad/api"
)

func Run(_ context.Context, conf *Config) error {
	var nomadSkipVerify bool
	if v := os.Getenv("NOMAD_SKIP_VERIFY"); v != "" {
		var err error
		nomadSkipVerify, err = strconv.ParseBool(v)
		if err != nil {
			return err
		}
	}

	client, err := api.NewClient(&api.Config{
		Address:  conf.NomadAddress,
		SecretID: os.Getenv("NOMAD_TOKEN"),
		TLSConfig: &api.TLSConfig{
			Insecure: nomadSkipVerify,
		},
	})
	if err != nil {
		log.Fatalf("Can't init client: %s", err)
	}

	cancelCh := make(chan os.Signal)
	signal.Notify(cancelCh, os.Interrupt, syscall.SIGINT)
	go func() {
		<-cancelCh
		fmt.Fprintf(os.Stderr, "\r- Ctrl+C pressed in Terminal\n")
		os.Exit(0)
	}()

	outputCh := make(chan string, 10)
	go LogRecordsPrinter(outputCh, *conf.FilterStr, *conf.ExcludeStr)

	subsCh := make(chan Subscription)
	go SubscriptionFinder(client, subsCh, conf.JobsOrAllocPrefix, conf.TaskName)

	wg := &sync.WaitGroup{}
	Subscriber(client, conf, subsCh, outputCh, wg)
	wg.Wait()

	return nil
}

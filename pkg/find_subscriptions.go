package pkg

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"log"
	"strings"
	"time"
)

func SubscriptionFinder(client *api.Client, subsCh chan<- Subscription, prefix string) {
	var currentSubs []*Subscription
	for {
		newSubs, err := findSubscriptions(client, prefix)
		if err != nil {
			log.Fatalf("Can't find subscriptions: %s", err)
		}
		if len(newSubs) == 0 {
			log.Printf("No jobs or allocations found by prefix `%s`", prefix)
		}
		for _, ns := range newSubs {
			found := false
			for _, cs := range currentSubs {
				if ns.String() == cs.String() {
					found = true
				}
			}
			if !found {
				subsCh <- *ns
				currentSubs = append(currentSubs, ns)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func findSubscriptions(client *api.Client, prefix string) ([]*Subscription, error) {
	var subs []*Subscription
	subs, err := findJobSubscription(client, prefix)
	if err != nil {
		return nil, fmt.Errorf("error on job subsriptions search: %v", err)
	}
	if len(subs) == 0 {
		subs, err = findAllocSubscription(client, prefix)
		if err != nil {
			return nil, fmt.Errorf("error on alloc subsriptions search: %v", err)
		}
	}
	return subs, nil
}

func findJobSubscription(client *api.Client, prefix string) ([]*Subscription, error) {
	jobs, _, err := client.Jobs().PrefixList(prefix)
	if err != nil {
		return nil, fmt.Errorf("can't make a search for a Job: %v", err)
	}
	var subs []*Subscription
	for _, j := range jobs {
		list, _, err := client.Jobs().Allocations(j.ID, false, nil)
		if err != nil {
			return nil, fmt.Errorf("can't make a search for a Job `%s` allocations: %v", j.Name, err)
		}
		for _, al := range list {
			if al.ClientStatus != api.AllocClientStatusRunning {
				continue
			}
			for t := range al.TaskStates {
				subs = append(subs, NewSubscription(al.NodeID, j.Name, al.ID, t))
			}
		}
	}
	return subs, nil
}

func findAllocSubscription(client *api.Client, prefix string) ([]*Subscription, error) {
	var subs []*Subscription
	list, _, err := client.Allocations().List(nil)
	if err != nil {
		return nil, fmt.Errorf("can't make a search for a allocation by prefix `%s`: %v", prefix, err)
	}
	for _, al := range list {
		if !strings.HasPrefix(al.ID, prefix) {
			continue
		}
		if al.ClientStatus != api.AllocClientStatusRunning {
			continue
		}
		for t := range al.TaskStates {
			subs = append(subs, NewSubscription(al.NodeID, al.JobID, al.ID, t))
		}
	}
	return subs, nil
}

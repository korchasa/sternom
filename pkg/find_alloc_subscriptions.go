package pkg

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
)

func findAllocSubscription(client *api.Client, prefix string) ([]*Subscription, error) {
	var subs []*Subscription
	list, _, err := client.Allocations().PrefixList(prefix)
	if err != nil {
		return nil, fmt.Errorf("can't make a search for a allocation by prefix `%s`: %v", prefix, err)
	}
	for _, al := range list {
		for t := range al.TaskStates {
			subs = append(subs, NewSubscription(al.NodeID, al.JobID, al.ID, t))
		}
	}
	return subs, nil
}
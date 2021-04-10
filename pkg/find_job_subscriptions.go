package pkg

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
)

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

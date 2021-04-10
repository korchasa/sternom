package pkg

import (
	"fmt"
	"strings"
)

type Subscription struct {
	Node  string
	Job   string
	Alloc string
	AllocShort string
	Task  string
}

func NewSubscription(node, job, alloc, task string) *Subscription {
	ap := strings.Split(alloc, "-")
	return &Subscription{
		Node:       node,
		Job:        job,
		Alloc:      alloc,
		AllocShort: ap[0],
		Task:       task,
	}
}

func (s *Subscription) String() string {
	return fmt.Sprintf("%s:%s[%s]", s.Job, s.AllocShort, s.Task)
}

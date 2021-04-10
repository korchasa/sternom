package pkg

import (
	"fmt"
	"strings"
)

type AllocSpec struct {
	jobName string
	allocID string
	nodeID string
	tasks []string
}

func (a *AllocSpec) String() string {
	return fmt.Sprintf("%s:%s:%s on %s", a.jobName, a.allocID, strings.Join(a.tasks, "|"), a.nodeID)
}



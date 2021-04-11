package pkg

import (
	"github.com/fatih/color"
	"hash/fnv"
)

var jobColors = []*color.Color{
	color.New(color.FgHiCyan),
	color.New(color.FgHiGreen),
	color.New(color.FgHiMagenta),
	color.New(color.FgHiYellow),
	color.New(color.FgHiBlue),
	color.New(color.FgHiRed),
}

var otherColors = []*color.Color{
	color.New(color.FgCyan),
	color.New(color.FgGreen),
	color.New(color.FgMagenta),
	color.New(color.FgYellow),
	color.New(color.FgBlue),
	color.New(color.FgRed),
}

func determineColors(job, alloc, task string) (jobColor, allocColor, taskColor *color.Color) {
	hash := fnv.New32()
	_, _ = hash.Write([]byte(job))
	jobColor = jobColors[hash.Sum32()%uint32(len(jobColors))]

	hash = fnv.New32()
	_, _ = hash.Write([]byte(alloc))
	allocColor = otherColors[hash.Sum32()%uint32(len(otherColors))]

	hash = fnv.New32()
	_, _ = hash.Write([]byte(task))
	taskColor = otherColors[hash.Sum32()%uint32(len(otherColors))]

	return
}

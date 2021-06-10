package pkg

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func LogRecordsPrinter(outputCh <-chan string, filterStr []string, excludeStr []string) {
	for str := range outputCh {
		orig := str
		included := true
		if len(filterStr) != 0 {
			included = false
			for _, f := range filterStr {
				if strings.Contains(str, f) {
					str = strings.ReplaceAll(str, f, color.New(color.FgHiRed).Sprint(f))
					included = true
				}
			}
		}
		for _, e := range excludeStr {
			if strings.Contains(orig, e) {
				included = false
			}
		}
		if included {
			fmt.Println(str)
		}
	}
}

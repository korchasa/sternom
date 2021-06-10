package pkg

import (
	"fmt"
	"strings"
)

func LogRecordsPrinter(outputCh <-chan string, filterStr []string, excludeStr []string) {
	for str := range outputCh {
		if isIncluded(str, filterStr, excludeStr) {
			fmt.Println(str)
		}
	}
}

func isIncluded(str string, filterStr []string, excludeStr []string) bool {
	included := true
	if len(filterStr) != 0 {
		included = false
		for _, f := range filterStr {
			if strings.Contains(str, f) {
				included = true
			}
		}
	}
	for _, e := range excludeStr {
		if strings.Contains(str, e) {
			included = false
		}
	}
	return included
}

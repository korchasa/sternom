package pkg

import "fmt"

func PrintLogRecord(outputCh <-chan string) {
	for str := range outputCh {
		fmt.Println(str)
	}
}

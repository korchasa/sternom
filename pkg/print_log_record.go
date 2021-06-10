package pkg

import "fmt"

func LogRecordsPrinter(outputCh <-chan string) {
	for str := range outputCh {
		fmt.Println(str)
	}
}

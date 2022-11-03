package pkg

import (
	"bufio"
	"io"
	"sync"

	"github.com/hashicorp/nomad/api"
)

func LogReader(prefix string, in <-chan *api.StreamFrame, out chan<- string) {
	var wg sync.WaitGroup

	pr, pw := io.Pipe()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer pw.Close()

		for data := range in {
			if data == nil {
				continue
			}

			pw.Write(data.Data)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(pr)

		for scanner.Scan() {
			out <- prefix + scanner.Text()
		}
	}()

	wg.Wait()
}

package app

import (
	"fmt"
	"sync"

	"gitlab.ozon.dev/chppppr/homework/internal/workers"
)

func ShowResult(wg *sync.WaitGroup, wk *workers.Workers) {
	defer wg.Done()

	for res := range wk.Results {
		fmt.Printf("\0337")
		fmt.Printf("\r\n\n\033[4F")
		fmt.Printf("\033[K")
		if res.Err != nil {
			fmt.Printf("\rError for response: %v\n", res.Response)
			fmt.Print("\033[K")
			fmt.Println(res.Err)
		} else {
			fmt.Printf("\rOK Response: %v\n", res.Response)
			fmt.Print("\033[K\n")
		}
		fmt.Print("\0338")
	}
}

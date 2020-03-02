package main

import (
	"fmt"
)

func main() {
	chCancel := make(chan struct{})
	chData := make(chan int)
	go func(chCancel chan struct{}, chData chan int) {
		var val int
		for {
			select {
			case <-chCancel:
				close(chData)
				return // exit func()
			case chData <- val:
				val++
			}
		}
	}(chCancel, chData)

	for val := range chData {
		fmt.Println("get value:", val)
		if val > 3 {
			chCancel <- struct{}{}
		}
	}
	close(chCancel)
}

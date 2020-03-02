package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func startNode(nodeNumber int, wg *sync.WaitGroup, chQuota chan struct{}) {
	defer wg.Done()
	chQuota <- struct{}{}             // get slot
	time.Sleep(10 * time.Millisecond) // имитация работы воркера
	for i := 0; i < 4; i++ {
		fmt.Printf("node %v: %v\n", nodeNumber, i)
		if i%2 == 0 {
			<-chQuota             // release slot
			chQuota <- struct{}{} // get slot
		}
		runtime.Gosched() // переключение
	}
	<-chQuota
}

func main() {
	var wg sync.WaitGroup
	chQuota := make(chan struct{}, 1) // quotaLmit = 1 размер слота
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go startNode(i, &wg, chQuota)
	}
	wg.Wait()
	fmt.Println("succesfully complete")
}

package main

import (
	"fmt"
	"runtime"
	"time"
)

func startNode(nodeNumber int, ch <-chan string) {
	for i := range ch {
		time.Sleep(10 * time.Millisecond) // имитация работы воркера
		fmt.Printf("Node - %d get from chanel - %s value - %s\n", nodeNumber, "chData", i)
		runtime.Gosched() // переключение
	}
	fmt.Printf("Node %v successfully completed tasks\n", nodeNumber)
}

func main() {
	list := []string{"Jan", "Feb", "Mar", "Apr", "May", "June", "Jul", "Aug", "Sep", "Nov", "Dec"}
	chData := make(chan string, 2)
	for i := 0; i < 3; i++ {
		go startNode(i, chData) // пул x 3 workers
	}
	// отправляем в канал строку очередь
	for _, month := range list {
		chData <- month
	}

	close(chData)

	time.Sleep(100 * time.Millisecond)
}

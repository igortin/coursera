package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// сущность которую можно передать в каждую gouroutine и кинуть сигнал в канал ctx.Done через функцию finish
	ctx, finish := context.WithCancel(context.Background())
	result := make(chan int)

	for i := 0; i < 10; i++ {

		// передаем контекст, номер воркера и канал окуда основная гоуритмна получит результат
		go worker(ctx, i, result)
	}

	// получаем результат из небуфферизованного канала
	fmt.Println("first goroutine procceed: ", <-result)

	// отправляем сигнал в канал ctx.Done() на прервание других гоурутин
	finish()
	close(result)
}

func worker(ctx context.Context, workerNum int, out chan<- int) {
	// имитируем работу
	waitTime := time.Duration(rand.Intn(100) + 10) //* time.Millisecond
	fmt.Println("WorkerNum:", workerNum, "work time:", waitTime)
	//создаем таймер который отправит сигнал после паузы длительностью waitTime
	timer := time.NewTimer(waitTime)
	// multiplexer

	select {
	case <-ctx.Done():
		return // stop goroutine если получили сигнал из канала
	case <-timer.C: // как только сигнал, значение workerNum отправиться в канал out
		out <- workerNum
	}
}

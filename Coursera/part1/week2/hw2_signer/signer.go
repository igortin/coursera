package main

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

var count int

func main() {
	inputData := []int{1, 2, 3, 5, 777}
	count = len(inputData)
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, val := range inputData {
				out <- val
			}
			// close(out)
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
	}
	ExecutePipeline(hashSignJobs...)
}

func ExecutePipeline(s ...job) {
	chIn := make(chan interface{})
	chOut := make(chan interface{})
	for _, funcItem := range s {
		go funcItem(chIn, chOut)
	}

LOOP:
	for {
		value := <-chIn
		switch value.(type) {
		case struct{}:
			break LOOP
		default:
			chIn <- value
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func SingleHash(in, out chan interface{}) {
	// start := time.Now()
	var c int
	wg := &sync.WaitGroup{}
	for {
		value := <-out
		switch value.(type) {
		case string:
			out <- value
			time.Sleep(time.Microsecond * 10)
			continue
		case int:
			wg.Add(1)
			c++
			data, _ := value.(int)
			dataStr := strconv.Itoa(data)

			chMd5 := make(chan string)
			defer close(chMd5)
			go func() {
				md5Get(dataStr, chMd5)
			}()
			go func() {
				defer wg.Done()
				chCrc321 := make(chan string)
				defer close(chCrc321)
				go func() {
					crc32Get(dataStr, chCrc321)
				}()

				chCrc322 := make(chan string)
				defer close(chCrc322)
				go func() {
					val := <-chMd5
					crc32Get(val, chCrc322)
				}()
				v1 := <-chCrc321
				v2 := <-chCrc322
				in <- v1 + "~" + v2
			}()
			wg.Wait()
		}
		if count == c {
			break
		}
	}
	// fmt.Println("Singlehash", time.Since(start))
	//fmt.Println("SingleHash completed")
}

func crc32Get(val string, ch chan string) {
	ch <- DataSignerCrc32(val)
}

func md5Get(data string, ch chan string) {
	ch <- DataSignerMd5(data)
}

func MultiHash(in, out chan interface{}) {
	// start_ := time.Now()
	var c int
	wgEx := &sync.WaitGroup{}
	for {
		value := <-in
		wgEx.Add(1)
		c++
		go func() {
			defer wgEx.Done()
			wgE := &sync.WaitGroup{}
			wgE.Add(1)
			go func() {
				defer wgE.Done()
				m := make(map[int]string, 5)
				mutex := &sync.Mutex{}
				wg := &sync.WaitGroup{}
				data, _ := value.(string)
				for i := 0; i < 6; i++ {
					wg.Add(1)
					go crc32GetTh(wg, data, i, m, mutex)
					runtime.Gosched()
				}
				wg.Wait()

				out <- m[0] + m[1] + m[2] + m[3] + m[4] + m[5]

				runtime.Gosched()
			}()
			wgE.Wait()
		}()
		if count == c {
			break
		}
	}
	wgEx.Wait()
	//fmt.Println("Multhash", time.Since(start_))
	//fmt.Println("MultiHash complete success")
}

func crc32GetTh(wg *sync.WaitGroup, data string, index int, m map[int]string, mutex *sync.Mutex) {
	mutex.Lock()
	m[index] = DataSignerCrc32(strconv.Itoa(index) + data)
	mutex.Unlock()
	wg.Done()
}

func CombineResults(in, out chan interface{}) {
	var c int
	defer close(out)
	s := make([]string, 0)
	for {
		val := <-out
		switch val.(type) {
		case string:
			c++
			data := val.(string)
			s = append(s, data)
		case int:
			out <- val
			time.Sleep(time.Microsecond * 10)
		}
		if count == c {
			break
		}
	}
	sort.Strings(s)
	var r string
	for _, v := range s {
		r = r + v + "_"
	}
	fmt.Println(r[:len(r)-1])
	// fmt.Println("777")
	in <- struct{}{}
	// fmt.Println("CombineResults successfull completed")
}

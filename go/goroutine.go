//package main
//
//import (
//	"fmt"
//	"os"
//	"os/signal"
//	"sync"
//	"syscall"
//	"time"
//)
//
//func main() {
//
//	//fmt.Println("return", test())
//
//	//goroutineTest()
//
//	running := true
//	f := func() {
//		for running {
//			fmt.Println("sub proc running...")
//			time.Sleep(1 * time.Second)
//		}
//		fmt.Println("sub proc exit")
//	}
//	go f()
//	go f()
//	go f()
//	time.Sleep(2 * time.Second)
//	running = false
//	time.Sleep(3 * time.Second)
//	fmt.Println("main proc exit")
//	ch := make(chan int)
//	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
//	go sum(s[:len(s)/2], ch)
//	go sum(s[len(s)/2:], ch)
//	x, y := <-ch, <-ch
//	fmt.Println(x, y, x+y)
//}
//
//func sum(s []int, c chan int) {
//	sum := 0
//	for _, v := range s {
//		sum += v
//	}
//	c <- sum
//}
//
//func goroutineTest() {
//	stop := make(chan bool)
//	var wg sync.WaitGroup
//	for i := 0; i < 3; i++ {
//		wg.Add(1)
//		go func(stop <-chan bool) {
//			defer wg.Done()
//			consumer(stop)
//		}(stop)
//	}
//	waitForSignal()
//	close(stop)
//	fmt.Println("stopping all jobs")
//	wg.Wait()
//}
//
//func waitForSignal() {
//	sigs := make(chan os.Signal)
//	signal.Notify(sigs, os.Interrupt)
//	signal.Notify(sigs, syscall.SIGTERM)
//	<-sigs
//}
//
//func consumer(stop <-chan bool) {
//	for {
//		select {
//		case <-stop:
//			fmt.Println("exit sub goroutine")
//			return
//		default:
//			fmt.Println("running")
//			time.Sleep(500 * time.Millisecond)
//		}
//	}
//}
//
//func signalHandler(stop <-chan bool) {
//	<-stop
//}
//
//func test() (i int) {
//	defer func() {
//		i += 1
//		fmt.Println("defer1", i)
//	}()
//	defer func() {
//		i += 1
//		fmt.Println("defer2", i)
//	}()
//	return i
//}
//
////func main() {
////	fmt.Println("return", test())
////}

package main

import "fmt"

func main() {
	chanNum := 26
	chanQueue := make([]chan struct{}, chanNum)
	var result = 0
	exitChan := make(chan struct{})
	for i := 0; i < chanNum; i++ {
		chanQueue[i] = make(chan struct{})
		if i == chanNum-1 {
			go func(i int) {
				chanQueue[i] <- struct{}{}
			}(i)
		}
	}
	for i := 0; i < chanNum; i++ {
		var lastChan, curChan chan struct{}
		if i == 0 {
			lastChan = chanQueue[chanNum-1]
		} else {
			lastChan = chanQueue[i-1]
		}
		curChan = chanQueue[i]
		go func(i byte, lastChan, curChan chan struct{}) {
			for {
				if result > 26 {
					exitChan <- struct{}{}
				}
				<-lastChan
				fmt.Printf("%c\n", i)
				result++
				curChan <- struct{}{}
			}
		}('A'+byte(i), lastChan, curChan)
	}
	<-exitChan
	fmt.Println("done")
}

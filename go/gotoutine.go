package main

import (
	"fmt"
	"time"
)

func main() {
	running := true
	f := func() {
		for running {
			fmt.Println("sub proc running...")
			time.Sleep(1 * time.Second)
		}
		fmt.Println("sub proc exit")
	}
	go f()
	go f()
	go f()
	time.Sleep(2 * time.Second)
	running = false
	time.Sleep(3 * time.Second)
	fmt.Println("main proc exit")
	ch := make(chan int)
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	go sum(s[:len(s)/2], ch)
	go sum(s[len(s)/2:], ch)
	x, y := <-ch, <-ch
	fmt.Println(x, y, x+y)
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

//package main

// func main() {
// 	goroutineTest()
// }

// func goroutineTest() {
// 	stop := make(chan bool)
// 	var wg sync.WaitGroup
// 	for i := 0; i < 3; i++ {
// 		wg.Add(1)
// 		go func(stop <-chan bool) {
// 			defer wg.Done()
// 			consumer(stop)
// 		}(stop)
// 	}
// 	waitForSignal()
// 	close(stop)
// 	fmt.Println("stopping all jobs")
// 	wg.Wait()
// }

// func waitForSignal() {
// 	sigs := make(chan os.Signal)
// 	signal.Notify(sigs, os.Interrupt)
// 	signal.Notify(sigs, syscall.SIGTERM)
// 	<-sigs
// }

// func consumer(stop <-chan bool) {
// 	for {
// 		select {
// 		case <-stop:
// 			fmt.Println("exit sub goroutine")
// 			return
// 		default:
// 			fmt.Println("running")
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 	}
// }

// func signalHandler(stop <-chan bool) {
// 	<-stop
// }

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
//func main() {
//	fmt.Println("return", test())
//}

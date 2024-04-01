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

//package main
//
//import "fmt"
//
//func main() {
//	chanNum := 26
//	chanQueue := make([]chan struct{}, chanNum)
//	var result = 0
//	exitChan := make(chan struct{})
//	for i := 0; i < chanNum; i++ {
//		chanQueue[i] = make(chan struct{})
//		if i == chanNum-1 {
//			go func(i int) {
//				chanQueue[i] <- struct{}{}
//			}(i)
//		}
//	}
//	for i := 0; i < chanNum; i++ {
//		var lastChan, curChan chan struct{}
//		if i == 0 {
//			lastChan = chanQueue[chanNum-1]
//		} else {
//			lastChan = chanQueue[i-1]
//		}
//		curChan = chanQueue[i]
//		go func(i byte, lastChan, curChan chan struct{}) {
//			for {
//				if result > 26 {
//					exitChan <- struct{}{}
//				}
//				<-lastChan
//				fmt.Printf("%c\n", i)
//				result++
//				curChan <- struct{}{}
//			}
//		}('A'+byte(i), lastChan, curChan)
//	}
//	<-exitChan
//	fmt.Println("done")
//}

//package main
//
//import (
//	"fmt"
//	"sync"
//)
//
//const (
//	MAX     = 100000
//	GoCount = 5
//)
//
//func main() {
//	solution(MAX, GoCount)
//}
//
//func solution(max, goCount int) *[]int {
//	lock := sync.Mutex{}
//	wg := sync.WaitGroup{}
//	result := make([]int, 0, MAX)
//	count := 1
//	wg.Add(GoCount)
//	for i := 0; i < goCount; i++ {
//		go func(i int) {
//			for {
//				lock.Lock()
//				now := count
//				lock.Unlock()
//				if now > max {
//					wg.Done()
//					return
//				}
//				if now%goCount == i {
//					fmt.Println(now)
//					result = append(result, now)
//					count++
//				}
//			}
//		}(i)
//	}
//	wg.Wait()
//	return &result
//}

//package main
//
//import (
//	"fmt"
//	"time"
//)
//
//const (
//	N      = 3
//	MAXNUM = 100
//)
//
//func main() {
//	var channels []chan struct{}
//	for i := 0; i < N; i++ {
//		channels = append(channels, make(chan struct{}, 1))
//	}
//	count := 0
//	for i := 0; i < N; i++ {
//		num := i
//		go func() {
//			for {
//				select {
//				case <-channels[num]:
//					if count > MAXNUM {
//						return
//					}
//					fmt.Printf("%c", rune('A'+num))
//					count++
//					time.Sleep(time.Millisecond * 300)
//					channels[(num+1)%N] <- struct{}{}
//				}
//			}
//		}()
//	}
//	channels[0] <- struct{}{}
//	select {}
//}

//

package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"sort"
)

func main() {
	var n, m, q int
	fmt.Scanf("%d %d %d", &n, &m, &q)

	queue := list.New()                   // 创建一个新的队列
	buyers := make(map[string]int)        // 创建一个map来记录每个人购买的抽赏数量
	scanner := bufio.NewScanner(os.Stdin) // 创建scanner读取输入
	for i := 0; i < q; i++ {
		scanner.Scan()
		line := scanner.Text()
		var op int
		fmt.Sscanf(line, "%d", &op)
		switch op {
		case 1:
			var name string
			fmt.Sscanf(line, "%d %s", &op, &name)
			queue.PushBack(name) // 将名字加入队列末尾
			if _, exists := buyers[name]; !exists {
				buyers[name] = 0 // 初始化购买数量为0
			}
		case 2:
			var name string
			fmt.Sscanf(line, "%d %s", &op, &name)
			for e := queue.Front(); e != nil; e = e.Next() {
				if e.Value == name {
					queue.Remove(e) // 移除队列中的名字
					break
				}
			}
		case 3:
			var x int
			fmt.Sscanf(line, "%d %d", &op, &x)
			if queue.Len() > 0 {
				name := queue.Remove(queue.Front()).(string) // 移除并返回队列最前面的名字
				buyers[name] += x
				n -= x
				if n <= m { // 检查是否触发特殊规则
					if queue.Len() > 0 {
						name := queue.Remove(queue.Front()).(string)
						buyers[name] += n
					}
					n = 0
					queue.Init() // 清空队列
				}
			}
		case 4:
			fmt.Println(queue.Len()) // 输出队列中的人数
		}
	}

	// 输出购买详情
	var names []string
	for name := range buyers {
		names = append(names, name)
	}
	sort.Strings(names) // 按字典序排序名字
	for _, name := range names {
		fmt.Printf("%s %d\n", name, buyers[name])
	}
}

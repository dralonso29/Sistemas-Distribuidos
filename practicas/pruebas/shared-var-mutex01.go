package main

import (
	"fmt"
	"sync"
	"strconv"
	"time"
)

type blocking struct {
	mylock chan struct{}
	i int
}

var b blocking //global var

func doStuff(i int, n *sync.WaitGroup)  {
	defer n.Done()

}

//!+main
func main() {
	fmt.Println("Comienza el programa")
	var n sync.WaitGroup
	b.ch1 = make(chan string, 5)
	// b.ch2 = make(chan strings, 5)
	for i := 0; i < 4; i++ {
		n.Add(1)
		str := "GORUTINA "+strconv.Itoa(i)
		b.ch1 <- str
		// }else{
		// 	b.ch2 <- str
		// }
		go doStuff(i, &n)
		time.Sleep(1*time.Second)
	}

	n.Wait()
	fmt.Println("Termina el programa")
}
//!-main

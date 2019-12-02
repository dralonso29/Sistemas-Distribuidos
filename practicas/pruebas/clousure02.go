package main

import (
	"fmt"
	"sync"
)

type blocking struct {
	num int
	blc chan struct{}
}

var b blocking
//!+main
func main() {
	fmt.Println("Comienza el programa")
	var n sync.WaitGroup
	for i := 0; i < 10; i++ {
		n.Add(1)
		go func(i int) {
			defer n.Done()
			fmt.Println("i = ", i)
		}(i)
	}
	n.Wait()
	fmt.Println("Termina el programa")
}
//!-main

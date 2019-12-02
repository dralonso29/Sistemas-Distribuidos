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
	n.Add(1)
	go func() {
		defer n.Done()
		fmt.Println("HOLA")
	}()
	n.Wait()
	fmt.Println("Termina el programa")
}
//!-main

package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

const NMATCHES = 4		// total number of matches
const MTIME	= 20000		// time of a match in ms

//!+myUnlock
func myUnlock(ch chan struct {})  {
	for i := 0; i < NMATCHES; i++ {
		fmt.Println("Unlock: ",i)
		<- ch
	}
}
//!-myUnlock

//!+myLock
func myLock(ch chan struct {})  {
	for i := 0; i < NMATCHES; i++ {
		fmt.Println("Lock: ",i)
		ch <- struct{}{}
	}
}
//!-myLock

//!+match
func match(start chan struct{}, id int, n *sync.WaitGroup)  {
	n.Done()
	fmt.Println("El partido ", id, " esperando")
	start <- struct{}{}
	fmt.Println("Partido ", id, ": ", time.Now().String())
}
//!-match

//!+matchesGenerator
func matchesGenerator()  {
	var n sync.WaitGroup
	n.Wait()
	start := make(chan struct{}, NMATCHES)
	myLock(start)
	for i := 0; i < NMATCHES; i++ {
		n.Add(1)
		go match(start, i, &n)
	}
	// hacer algun algoritmo para poder hacer que todos los partidos empiecen a la vez
	myUnlock(start)
}
//!-matchesGenerator

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	matchesGenerator()
	fmt.Println("Termina el programa")
}
//!-main

package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

const NMATCHES = 1		// total number of matches
const MTIME	= 20000		// time of a match in ms
const CDTIME = 750		// countdown time in ms
const ENDMATCH = 5		// playtime of a match is 20 secs (10 chances * 2 seconds between chances = 20 seconds)
const TCHANCE = 2000	// time between chances in ms

// type mdata struct {
//
// }

//!+myUnlock
func myUnlock(ch chan struct {})  {
	for i := 0; i < NMATCHES; i++ {
		//fmt.Println("Unlock: ",i)
		<- ch
	}
}
//!-myUnlock

//!+myLock
func myLock(ch chan struct {})  {
	for i := 0; i < NMATCHES; i++ {
		//fmt.Println("Lock: ",i)
		ch <- struct{}{}
	}
}
//!-myLock

//!+playmatch
func playmatch()  {
	team := ""
	occ := ""
	prob1 := 50
	prob2 := 70
	n1 := rand.Intn(100) + 1
	if n1 <= prob1 {
		team = "local"
	}else{
		team = "visitor"
	}
	n2 := rand.Intn(100) + 1
	if n2 <= prob2 {
		occ = "fail"
	}else{
		occ = "goal"
	}
	fmt.Println("Ocasion del ", team, " ha sido ", occ, "n1 = ", n1, ", n2 = ", n2)
}
//!-playmatch

//!+match
func match(start chan struct{}, id int, n *sync.WaitGroup)  {
	defer n.Done()
	start <- struct{}{}
	for i := 1; i <= ENDMATCH; i++ {
		playmatch()		// generate a goal
		time.Sleep(TCHANCE * time.Millisecond)
	}
	// fmt.Println("Partido ", id, ": ", time.Now().String())
}
//!-match

//!+startCountdown
func startCountdown()  {
	fmt.Println("The ", NMATCHES, " matches start in...")
	for i := 3; i >= 1 ; i-- {
		fmt.Printf("\r%d",i)
		time.Sleep(CDTIME * time.Millisecond)
	}
	fmt.Printf("\rGO!\n")
}
//!-startCountdown

//!+matchesGenerator
func matchesGenerator()  {
	var n sync.WaitGroup
	defer n.Wait()
	start := make(chan struct{}, NMATCHES)
	myLock(start)
	for i := 0; i < NMATCHES; i++ {
		n.Add(1)
		go match(start, i, &n)
	}
	startCountdown()
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

package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)


const NCOCHES = 4 //numero total de coches
const MINSEC = 1
const MAXSEC = 3

type car struct {
	id string
	way string
	state string
}

type bridge struct {
	eCars int
	wCars int
	totalCars int
	//arrayCars car[NCOCHES]
}

func randomTime(n int)  int { // returns int between 1 and n (both included)

	return rand.Intn(n) + MINSEC
}

func sideSelector(prob int)  string{
	n := rand.Intn(100) + 1
	if n <= prob {
		return "east"
	}
	return "west"
}

func carRunning(id int, side string, n *sync.WaitGroup)  {
	defer n.Done()
	fmt.Println("Coche ",id," en el lado ", side)
}

func carGenerator()  {
	prob := rand.Intn(100) + 1
	var n sync.WaitGroup
	for i := 0; i < NCOCHES; i++ {
		side := sideSelector(prob)
		n.Add(1)
		go carRunning(i+1, side, &n) // goroutines
		time.Sleep(time.Duration(randomTime(MAXSEC)) * time.Second) //we'll wait random time between MINSEC and MAXSEC
	}
	n.Wait()
}

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //importante que sino en todas las ejecuciones nos va a salir el mismo numero
	fmt.Println("Comienza el programa")
	carGenerator()
	fmt.Println("Termina el programa")
}
//!-main

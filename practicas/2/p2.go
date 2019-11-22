package main

import (
	"fmt"
	"time"
	"math/rand"
)


const NCOCHES = 10 //numero total de coches
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

func carRunning(id int, side string)  {
	fmt.Println("Coche ",id," en el lado ", side)
}

func carGenerator()  {
	prob := rand.Intn(100) + 1
	for i := 0; i < NCOCHES; i++ {
		side := sideSelector(prob)
		carRunning(i+1, side) // goroutines
		time.Sleep(time.Duration(randomTime(MAXSEC)) * time.Second) //we'll wait random time between MINSEC and MAXSEC
	}
}

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //importante que sino en todas las ejecuciones nos va a salir el mismo numero
	fmt.Println("Comienza el programa")
	//aqui iria la gorutina que debe esperar a todos los coches
	//y aqui se deberia hacer el carGenerator
	carGenerator()
}
//!-main

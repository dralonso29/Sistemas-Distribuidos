package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)


const NCOCHES = 4 		//numero total de coches
const MAXSEC = 3000 	//ms

type car struct {
	id int				// car's name (just an identifier)
	side string			// east or west
	state string		// arriving, waiting, crossing, finished
}

type bridge struct {
	eastch chan <- car
	westch chan <- car
	way string			// from east to west or from west to east
}

func sideSelector(prob int)  string{
	n := rand.Intn(100) + 1
	if n <= prob {
		return "east"
	}
	return "west"
}

func arrive2bridge(c *car)  {
	fmt.Println("Coche ",c.id," llegando al puente")
	time.Sleep(time.Duration(rand.Intn(MAXSEC)) * time.Millisecond)
}

func carRunning(c *car, n *sync.WaitGroup)  {
	defer n.Done()
	fmt.Println("Coche ",c.id," creado en el lado ", c.side)
	c.state = "waiting" // primero habria que mirar si se puede entrar, en caso contrario se pondria el estado waiting
	//cuando van llegando al puente hay que ir metiendo a los coches en el buffered channel que corresponda
	arrive2bridge(c)
	fmt.Println("Coche ",c.id," ha llegado al puente")
}

func carGenerator()  {
	//br := new(bridge)
	prob := rand.Intn(100) + 1
	var n sync.WaitGroup
	for i := 0; i < NCOCHES; i++ {
		side := sideSelector(prob)
		n.Add(1)
		c := new(car)
		c.id = i+1
		c.side = side
		go carRunning(c, &n) // goroutines
		time.Sleep(time.Duration(rand.Intn(MAXSEC)) * time.Millisecond) //we'll wait random time between zero and MAXSEC
	}
	n.Wait()
}

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	carGenerator()
	fmt.Println("Termina el programa")
}
//!-main

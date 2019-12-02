package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

const NCARS = 3			//numero total de coches
const MAXSEC = 3000 	//ms

type car struct {
	id int				// car's name (just an identifier)
	side string			// east or west
	state string		// arriving, waiting, crossing, finished
}

type bridge struct {
	sync.Mutex
	eastch chan <- car
	westch chan <- car
	way string			// from east to west or from west to east
	crossing chan struct{}
}

var brg bridge			// global var. Used by all cars

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

func setWay(side string)  string{
	if side == "east" {
		return "west"
	}
	return "east"
}

func carRunning(c *car, n *sync.WaitGroup)  {
	defer n.Done()
	fmt.Println("Coche ",c.id," creado en el lado ", c.side)
	c.state = "waiting" // primero habria que mirar si se puede entrar, en caso contrario se pondria el estado waiting
	//cuando van llegando al puente hay que ir metiendo a los coches en el buffered channel que corresponda
	arrive2bridge(c)
	fmt.Println("Coche ",c.id," ha llegado al puente")

	if len(brg.way) == 0 {
		go func() {
			brg.crossing <- struct{}{}
			fmt.Println("El puente esta vacio")
			brg.way = setWay(c.side)
			fmt.Println("El coche ",c.id,"ha establecido el sentido hacia ", brg.way)
			time.Sleep(4000 * time.Millisecond)
		}()
		time.Sleep(4000 * time.Millisecond)
		<- brg.crossing
	}
}

func initializeBridge()  {
	brg.way = ""
	brg.eastch = make(chan car, NCARS)
	brg.westch = make(chan car, NCARS)
	brg.crossing = make(chan struct{})
}

func carGenerator()  {
	initializeBridge()
	prob := rand.Intn(100) + 1
	var n sync.WaitGroup
	for i := 0; i < NCARS; i++ {
		n.Add(1)
		c := new(car)
		c.id = i+1
		c.side = sideSelector(prob)
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

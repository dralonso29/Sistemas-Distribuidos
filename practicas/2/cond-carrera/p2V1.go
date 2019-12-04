package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

const NCARS = 6		// numero total de coches
const MAXSEC = 3000 	// ms
const CROSSTIME = 6000	// ms
const CARSONBRG = 5 	// cars limit on bridge
const NTOKEN = 1

type car struct {
	id int				// car's name (just an identifier)
	side string			// east or west
}

type bridge struct {
	sync.Mutex
	eastch chan car
	westch chan car
	token chan struct{}
	way string			// from east to west or from west to east
	eastcars int
	westcars int
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

func lockSideBridge(brgside string)  {
	for i := 0; i < CARSONBRG; i++ {
		//fmt.Println("Lock: ", i)
		if brgside == "east" {
			brg.eastch <- car{0, ""}
			continue
		}
		brg.westch <- car{0, ""}
	}
}

func unlockSideBridge(brgside string)  {
	for i := 0; i < CARSONBRG; i++ {
		//fmt.Println("Unlock: ", i)
		if brgside == "east" {
			<- brg.eastch
			continue
		}
		<- brg.westch
	}
}

func carRunning(c *car, n *sync.WaitGroup)  {
	defer n.Done()
	fmt.Println("Coche ",c.id," creado en el lado ", c.side)
	arrive2bridge(c)
	fmt.Println("Coche ",c.id," ha llegado al puente")

	//<- brg.token
	if c.side == "east"{
		brg.eastch <- car{c.id, c.side}
	}else{
		brg.westch <- car{c.id, c.side}
	}
	brg.Lock()
	if len(brg.way) == 0 {
		fmt.Println("El puente esta vacio")
		brg.way = setWay(c.side)
		fmt.Println("El coche ",c.id,"ha establecido el sentido hacia ", brg.way)
		lockSideBridge(brg.way)
	}
	//brg.token <- struct{}{}
	brg.Unlock()
	fmt.Println("Coche ",c.id," cruzando el puente")
	time.Sleep(CROSSTIME * time.Millisecond)
	brg.Lock()
	fmt.Println("Coche ",c.id, " ya ha cruzado")
	if c.side == "east" {
		<- brg.eastch
		//fmt.Println("Coche ", c.id,": numero coches este: ", brg.eastcars, " , long eastch: ", len(brg.eastch))
		if len(brg.eastch) == 0 {
			//fmt.Println("Coche ",c.id, " es el ultimo del lado ",c.side, ", desbloquea al otro lado")
			unlockSideBridge(brg.way)
			brg.way = ""
		}
		brg.Unlock()
	}else{
		<- brg.westch
		//fmt.Println("Coche ", c.id,": numero coches oeste: ", brg.westcars, " , long westch: ", len(brg.westch), "carsLeft = ", brg.carsLeft)
		if len(brg.westch) == 0 {
			//fmt.Println("Coche ",c.id, " es el ultimo del lado ",c.side, ", desbloquea al otro lado")
			unlockSideBridge(brg.way)
			brg.way = ""
		}
		brg.Unlock()
	}
}

func initializeBridge()  {
	brg.way = ""
	brg.eastch = make(chan car, CARSONBRG)
	brg.westch = make(chan car, CARSONBRG)
	//brg.token = make(chan struct{}, NTOKEN)
	brg.eastcars = 0
	brg.westcars = 0
}

func countSideCars(side string)  {
	defer brg.Unlock()
	brg.Lock()
	if side == "east" {
		brg.eastcars++
		return
	}
	brg.westcars++
}

func carGenerator()  {
	var n sync.WaitGroup
	defer n.Wait()
	initializeBridge()
	//brg.token <- struct{}{}
	prob := rand.Intn(100) + 1
	for i := 0; i < NCARS; i++ {
		n.Add(1)
		c := new(car)
		c.id = i+1
		c.side = sideSelector(prob)
		countSideCars(c.side)
		go carRunning(c, &n) // goroutines
		time.Sleep(time.Duration(rand.Intn(MAXSEC)) * time.Millisecond) //we'll wait random time between zero and MAXSEC
	}
}

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	carGenerator()
	//<- brg.token 
	fmt.Println("Se han creado ",brg.eastcars," coches east y ",brg.westcars," coches west")
	fmt.Println("Termina el programa")
}
//!-main

package main

import (
	"fmt"
	"time"
)

func cliente(i int, chRecepcionista chan int, salaEspera chan chan int) {
	c := make(chan int)
	select {
	case chRecepcionista <- 1:
		x := <-chRecepcionista
		if x == 1 {
			salaEspera <- c
			fmt.Println("Cliente", i, " :me siento en la sala de espera")
			<-c
			fmt.Println("Cliente", i, " :me empiezan a cortar el pelo")
			<-c
			fmt.Println("Cliente", i, " :terminan de cortarme el pelo")
			chRecepcionista <- 0
		} else {
			fmt.Println("Cliente", i, "me voy de la barberia, esta llena")
		}

	}
}

func barbero(i int, salaEspera chan chan int) {
	for {
		select {
		case k := <-salaEspera:
			fmt.Println("Barbero ", i, " : Comienzo a cortar el pelo")
			k <- 2
			time.Sleep(1 * time.Second)
			k <- 2
			fmt.Println("Barbero ", i, " : Termino de cortar el pelo")
		}
	}
}

func recepcionista(chRecepcionista chan int, salaEspera chan chan int) {
	i := 0
	for {
		select {
		case x := <-chRecepcionista:
			if x == 1 {
				if i < 5 {
					i++
					chRecepcionista <- 1
				} else {
					chRecepcionista <- 0
				}
			} else {
				i--
			}
		}

	}
}

func main() {
	salaEspera := make(chan chan int)
	chRecepcionista := make(chan int)

	go func() {
		barbero(1, salaEspera)

	}()

	go func() {
		barbero(2, salaEspera)
	}()

	go func() {
		recepcionista(chRecepcionista, salaEspera)
	}()

	i := 0
	for i < 200 {
		go func() {
			cliente(i, chRecepcionista, salaEspera)
		}()
		time.Sleep(200 * time.Millisecond)
		i++
	}

	time.Sleep(4 * time.Second)
}

package main

import (
	"fmt"
	// "time"
)
const NCARS = 5

//!+main
func main() {
	fmt.Println("Comienza el programa")
	bufchan := make(chan struct{}, 5)
	for i := 0; i < 3; i++ {
		bufchan <- struct{}{}
	}
	close(bufchan) // antes de iterar hay que cerrar el channel porque sino se queda esperando a que alguien le envie cosas
	fmt.Println("Longitud del canal: ", len(bufchan))
	for i := range bufchan {
		fmt.Println("Iteracion ",i)
	}
	fmt.Println("Termina el programa")
}
//!-main

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"fmt"
)

type client struct {
	channel chan<- string // an outgoing message channel
	id string
	privConn bool
}

func findClient(clients map[client]bool, name string)  client{ //returns a client
	for cli := range clients {
		if cli.id ==  name{
			return cli
		}
	}
	return client{nil, "", false}
}

//!+main
func main() {
	clients := make(map[client]bool)
	ch := make(chan string)
	clients[client{ch, "pepe", false}] = true
	//no podria buscar en el mapa solo por el identificador "pepe". Tengo que pasarle todo
	fmt.Println(clients[client{ch, "pepe", false}])

	cli := findClient(clients, "pepe")
	if len(cli.id) > 0 {
		fmt.Println("Buscamos a pepe -> id = ", cli.id)
	}
	cli2 := findClient(clients, "juan")
	if len(cli2.id) > 0 {
		fmt.Println("Buscamos a juan -> id = ", cli2.id)
	}

}

//!-main

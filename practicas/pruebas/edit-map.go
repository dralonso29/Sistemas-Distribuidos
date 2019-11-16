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
	name string
	connected bool
}

func editMap(clients map[client]bool, cli client) {
	clients[cli] = false
	cli.connected = true
	fmt.Println("editMap: Modificamos a ", cli.name, " poniendo a false su value")
	fmt.Println("editMap: Modificamos a ", cli.name, " poniendo a true su campo connected: ", cli.connected)
}

//!+main
func main() {
	clients := make(map[client]bool)
	cli := client{"pepe", false}
	clients[cli] = true
	editMap(clients, cli)
	fmt.Println("main: el value de pepe es ",clients[cli])
	fmt.Println("main: el campo connected de pepe es ",cli.connected)
}

//!-main

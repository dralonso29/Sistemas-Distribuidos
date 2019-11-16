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

//!+main
func main() {
	clients := make(map[client]bool)
	cli := client{"pepe", false}
	clients[cli] = true

}

//!-main

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 219.
//!+

// Clock1 is a TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
	"flag"
	"fmt"
	"strconv"
)

func main() {
	//Se nos devuelve un puntero a la direccion donde se almacena el int
	port := flag.Int("port", 8000, "an int");
	flag.Parse();
	// Si no ponemos el '*', imprimimos la direccion de memoria donde se almacena el int.
	allocated := "localhost:"+strconv.Itoa(*port);
	listener, err := net.Listen("tcp", allocated)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		handleConn(conn) // handle one connection at a time
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

//!-

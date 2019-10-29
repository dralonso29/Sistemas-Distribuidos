// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 221.
//!+

// Netcat1 is a read-only TCP client.
package main

import (
	"io"
	"log"
	"net"
	"fmt"
	"os"
	"strings"
	"bufio"
)

func main() {
	argc := len(os.Args)
	if argc < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <zone1> <zone2> ... <zoneN>\n", os.Args[0])
		os.Exit(1)
	}
	for _, argv := range os.Args[1:] {
		splited := strings.Split(argv, "=");
		if len(splited) != 2{
			fmt.Fprintf(os.Stderr, "Error: invalid arg: %s\n", splited)
			os.Exit(1)
		}
		city := splited[0]
		host := splited[1]
		fmt.Println(city + ": "+host)
		go showTime(city, host);
	}
	for {
		//infinite loop
	}
}

func showTime(city string, host string)  {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(city, os.Stdout, conn)
}

func mustCopy(city string, dst io.Writer, src io.Reader) {
	buf := bufio.NewScanner(src)
	for buf.Scan() {
		// no podemos usar el \r porque no sabemos en que orden van a salir las horas y estariamos
		// reemplazando horas que no deberian ser reemplazadas
		fmt.Fprintf(dst, "%s: %s\n",city,buf.Text())
	}
	if buf.Err() != nil {
		fmt.Fprintf(os.Stderr, "Unespected error: %s\n", buf.Err())
		os.Exit(1)
	}
}

//!-

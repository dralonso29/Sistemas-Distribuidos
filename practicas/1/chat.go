// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

//!+broadcaster
type client struct {
	channel chan<- string // an outgoing message channel
	id string
}

type climsg struct {
	msg string
	who string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan climsg) // all incoming client messages
	broadcast = make(chan string)
)

func announceClients(clients map[client]bool)  {
	broadcast <- "List of clients:"
	for cli := range clients {
		go func(c client) {
			broadcast <- c.id
		}(cli)
	}
}

func sendMsg(clients map[client]bool, cm climsg)  {
	for cli := range clients {
		if cli.id != cm.who{
			go func(c client) {
				c.channel <- cm.msg
			}(cli)
		}
	}
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			go sendMsg(clients, msg)

		case msg := <- broadcast:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.channel <- msg
			}

		case cli := <-entering:
			clients[cli] = true
			// si no ponemos a announceClients en una gorutina
			// se queda bloqueado esperando a que se saque algo del canal
			// El problema es que la funcion que saca cosas de los canales es broadcaster
			go announceClients(clients)

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.channel)
			// aqui pasa lo mismo que en el case anterior
			go announceClients(clients)
		}
	}
}

//!-broadcaster

//!+handleConn
func introduceName(ch chan <- string, out chan <- string, conn net.Conn)  {
	ch <- "Introduce un nombre de usuario: "
	input := bufio.NewScanner(conn)
	for input.Scan() {
		out <- input.Text()
		// aqui se podria enviar por un canal a broadcaster en vez de por out
		// el nombre para ver si esta repetido
		break
	}
	// NOTE: ignoring potential errors from input.Err()
}


func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	namech := make(chan string)
	go clientWriter(conn, ch)

	// who := conn.RemoteAddr().String()
	go introduceName(ch, namech, conn)
	who := <- namech
	ch <- "You are " + who
	broadcast <- who + " has arrived"
	entering <- client{ch, who}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg := who + ": " + input.Text()
		messages <- climsg{msg, who}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- client{ch, who}
	broadcast <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main

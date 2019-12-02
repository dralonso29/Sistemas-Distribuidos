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
	"strings"
)

//!+broadcaster
type client struct {
	channel chan<- string // an outgoing message channel
	id string
	privConn bool
}

type climsg struct {
	msg string
	who string
}

type checkclient struct {
	name string
	repeated bool
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan climsg) // multicast channel. send msgs to all clients except to sender
	broadcast = make(chan string) // all incoming client messages
	validclient = make(chan checkclient) // channel to know if a username exists or not
)

const PRIVATE = "!private"

func announceClients(clients map[client]bool)  {
	broadcast <- "List of clients:"
	for cli := range clients {
		go func(c client) {
			broadcast <- c.id
		}(cli)
	}
}

func sendMsg(clients map[client]bool, senderclient climsg)  {
	for cli := range clients {
		if cli.id != senderclient.who{
			go func(c client) {
				c.channel <- senderclient.msg
			}(cli)
		}
	}
}

func findClient(clients map[client]bool, name string)  client{ //returns a client
	for cli := range clients {
		if cli.id ==  name{
			return cli
		}
	}
	return client{nil, "", false}
}

func privateRequest(sender climsg, token string, clients map[client]bool)  bool{
	splited := strings.Fields(sender.msg)
	if len(splited)==3 && (splited[1]==token) {
		name := splited[2]
		cli := findClient(clients, name)
		reqCli := findClient(clients, sender.who)

		// To do a private communication we need to guarantee that exists client and
		// is not in another private communication
		if cli.id != "" && reqCli.id != "" && clients[cli] && clients[reqCli] {
			return true
		}
	}
	return false
}

func getCliNameFromMsg(msg string)  string{
	return strings.Fields(msg)[2]
}

func broadcaster() {
	clients := make(map[client]bool) // bool indicates if client is free to use a priv conn
	for {
		select {
		case msg := <-messages:
			fmt.Println(privateRequest(msg, PRIVATE, clients))
			cli := findClient(clients, msg.who)

			if privateRequest(msg, PRIVATE, clients) {
				clients[cli] = false
				name := getCliNameFromMsg(msg.msg)
				cli2 := findClient(clients,name)
				clients[cli2] = false
			}

			if clients[cli] {
				go sendMsg(clients, msg)
			}
			// if privateRequest(msg, PRIVATE) {
			// 	//crear canal privado en una gorutina
			// }else{
			// 	go sendMsg(clients, msg)
			// }


		case msg := <- broadcast:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.channel <- msg
			}

		case validcli := <- validclient:
			for cli := range clients {
				if cli.id == validcli.name {
					validcli.repeated = true
					break
				}
			}
			validclient <- validcli

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
	ch <- "Enter a username: "
	input := bufio.NewScanner(conn)
	for input.Scan() {
		name := input.Text()
		if len(name) <= 0 {
			ch <- "User "+name+" too short. Try again: "
			continue
		}
		isrepeated := false // we asume that client is not repeated
		validclient <- checkclient{name, isrepeated}
		cli := <- validclient
		if !cli.repeated {
			out <- input.Text()
			break
		}
		ch <- "User "+name+" already used. Try again: "
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
	entering <- client{ch, who, false}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		text := input.Text()
		if len(text) > 0 { // to avoid send empty messages
			msg := who + ": " + text
			messages <- climsg{msg, who}
		}
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- client{ch, who, false}
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

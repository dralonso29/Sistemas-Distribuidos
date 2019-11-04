
package main

import (
	"fmt"
)
// Si intentamos enviar un nil por un canal de strings vemos que peta
// alonsod@dasus:~/gows/practicas/pruebas$ go run send-nil.go
// # command-line-arguments
// ./send-nil.go:13:11: cannot use nil as type string in send

var channel = make(chan string)

func main() {
	go func() {
		channel <- nil //si ponemos el string "hola", vemos que funciona correctamente
	}()
	fmt.Println(<-channel)
}

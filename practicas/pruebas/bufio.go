package main

import (
	"fmt"
	"os"
	"bufio"
)
// En este caso, solamente se ven los numeros pares del bucle for,
// porque en las iteraciones pares se escribe lo que hay en el canal, y
// en las impares se leen los valores pares que estan en el canal
func main() {
    fmt.Print("Introduce un nombre: ")
	input := bufio.NewScanner(os.Stdin)
	for input.Scan(){
		fmt.Println("Has metido:"+ input.Text())
		break
	}

}

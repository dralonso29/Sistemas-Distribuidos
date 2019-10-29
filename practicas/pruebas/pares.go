
package main

import (
	"fmt"
)
// En este caso, solamente se ven los numeros pares del bucle for,
// porque en las iteraciones pares se escribe lo que hay en el canal, y
// en las impares se leen los valores pares que estan en el canal
func main() {
    ch := make(chan int, 1) // probar a poner distintos valores de tamaño de canal
    for i := 0; i < 10; i++ {
        select {
        case x := <-ch:
            fmt.Println(x) // "0" "2" "4" "6" "8"
        case ch <- i:
        }
    }
}

// En el caso de que el buffer sea mayor de 1, vemos que cada ejecucion es diferente
// Esto se debe a que en cada ejecucion, el select coge de manera aleatoria uno de los
// case, por eso es no determinista. Vemos que algunas veces se imprime el 1, otras no, etc
// alonsod@dasus:~/gows/practicas/pruebas$ go run pares.go
// 0
// 2
// 3
// 5
// alonsod@dasus:~/gows/practicas/pruebas$ go run pares.go
// 0
// 1
// 4
// 5
// 8

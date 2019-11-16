package main

import (
	"fmt"
)
// Si queremos saber si un elemento esta dentro de un mapa, tenemos que mirar
// el valor que nos devuelve el diccionario (la variable que he llamado exists), porque si
// hiciesemos usr := mymap["pepe"], nos devolveria el valor booleano asociado. En este caso
//he puesto que sea false aposta para que se vea bien el problema.

//!+main
func main() {
	mymap := make(map[string]bool)
	mymap["pepe"] = false // añadimos solamente a pepe, con value false

	_, exists1 := mymap["juan"]
	fmt.Println("Esta juan en el mapa? -->",exists1)
	_, exists2 := mymap["pepe"]
	fmt.Println("Esta pepe en el mapa? -->",exists2)
}

//!-main

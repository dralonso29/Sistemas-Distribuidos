package main

import (
	"fmt"
	"time"
	"math/rand"
)

const MAX = 4

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	p := [MAX]string{"a","b","c","d"}
	ind := rand.Intn(MAX)
	fmt.Println("Indice: ",ind)
	l := len(p)
	for i := ind; i < l+ind; i++ {
		fmt.Println(p[i%l], ", i%l: ",i%l)
	}
	fmt.Println("Termina el programa")
}
//!-main

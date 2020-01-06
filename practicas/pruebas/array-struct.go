package main

import (
	"fmt"
	"time"
	"math/rand"
)

type mydata struct {
	local int
	visitor int
	invalid bool
}

const MAX = 4

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	a := []mydata{}
	for i := 0; i < MAX; i++ {
		n := mydata{local: 0, visitor: 0, invalid: false}
		a = append(a,n)
	}
	fmt.Println(a)
	fmt.Println(a[0].invalid)
	// p := [MAX]string{"a","b","c","d"}
	// ind := rand.Intn(MAX)
	// fmt.Println("Indice: ",ind)
	// l := len(p)
	// for i := ind; i < l+ind; i++ {
	// 	fmt.Println(p[i%l], ", i%l: ",i%l)
	// }
	fmt.Println("Termina el programa")
}
//!-main

package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

// Let's try if 2 gorutines can start at the same time

const NGORUTINES = 2

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
    
	for i := 0; i < NGORUTINES; i++ {
        go func(i int) {
            fmt.Println("G",i, " creada")
        }(i)
    }
	fmt.Println("Termina el programa")
}
//!-main

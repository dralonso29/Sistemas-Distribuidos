package main

import (
	"fmt"
    "sync"
	"time"
)
const N = 5 // numero de gorutinas a lanzar
const TIME1 = 4	// tiempo de espera para liberar el buffered channel

func lockAll(wg *sync.WaitGroup, bufchan chan struct {})  {
    defer wg.Done()
    fmt.Println("Bloqueamos todo")
    fmt.Println("Longitud del canal: ", len(bufchan))
    for i := 0; i < N; i++ {
		bufchan <- struct{}{}
        fmt.Println("Longitud del canal: ", len(bufchan))
	}
}

func unlockGR(wg *sync.WaitGroup, bufchan chan struct {})  {
    defer wg.Done()
    <- bufchan
    fmt.Println("Me han desbloqueado")
}

//!+main
func main() {
	fmt.Println("Comienza el programa")
	bufchan := make(chan struct{}, N)
    var wg sync.WaitGroup
    wg.Add(1)
    go lockAll(&wg, bufchan)
    time.Sleep(TIME1 * time.Second)
    for i := 0; i < N; i++ {
        wg.Add(1)
        go unlockGR(&wg, bufchan)
    }
    wg.Wait()

	close(bufchan)
	fmt.Println("Longitud del canal despues de leer todo: ", len(bufchan))
	fmt.Println("Termina el programa")
}
//!-main

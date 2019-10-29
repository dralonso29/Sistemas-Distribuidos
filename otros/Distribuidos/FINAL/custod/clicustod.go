package main

import (
	"log"
	"fmt"
	"net/rpc"
	"time"
)

type Arg struct{
	Num int
	St  string
}

func main() {
	clicustod, err := rpc.DialHTTP("tcp", "localhost"+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &Arg{3, "Hola"}
	var reply int
	err = clicustod.Call("Srv.Hola", args, &reply)
	if err != nil {
		log.Fatal("argumentos error:", err)
	}
	fmt.Printf("Conectando: ", args.Num, args.St, reply)
	time.Sleep(1 * time.Hour)
}

package main

import (
	"fmt"
	"net"
	"time"
	"net/http"
	"net/rpc"
	"log"
)

type Args struct {
	Num int
	St  string
}

//record con el estado del sevidor
type Srv struct{
}

var contador= 0

func (sr *Srv) Hola(n *Args, reply *int) error {
	fmt.Println("contador", contador)
	contador++
	fmt.Println("N:", n.Num, "String", n.St)
	return nil
}

func main() {
	arith := new(Srv)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	time.Sleep(3 * time.Hour)
}

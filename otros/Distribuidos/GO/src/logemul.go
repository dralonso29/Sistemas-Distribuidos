//Jorge Vela Peña
//Grado en Ingeniería Telemática
package main

import(
	"fmt"
	"logicclock"
	"logiclog"
	"os"
)

type sendMessage struct{
	msg string
	name string
	lg []logicclock.Log
}

func main(){
	clockPC1 := logicclock.NewClock()
	clockPC2 := logicclock.NewClock()

    f1, _ := os.Create("/tmp/dat")
    //check(err)

	//channelPC1 := make(chan sendMessage)

    xlog:=sendMessage{"Mensaje","PC1",clockPC1.GetClock()}
	clockPC2.Message(xlog.msg,xlog.name, &xlog.lg) //PC1 to PC2
	l := logiclog.NewLog(clockPC2.GetClock(), xlog.msg, xlog.name)
	l.CreateLog(f1)
	
    xlog2:=sendMessage{"Mensaje2","PC2",clockPC2.GetClock()}
	clockPC1.Message(xlog2.msg,xlog2.name, &xlog2.lg) //PC1 to PC2
	l1 := logiclog.NewLog(clockPC1.GetClock(), xlog2.msg, xlog2.name)
	l1.CreateLog(f1)  

	fmt.Println("Hola")
}
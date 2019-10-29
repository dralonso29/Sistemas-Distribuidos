//Jorge Vela Peña
//Grado en Ingeniería Telemática


package main

import(
	"fmt"
	"logicclock"
	"os"
)

type sendMessage struct{
	msg string
	name string 
}

func main(){
	clockPaco := logicclock.NewClock()
	clockJuan := logicclock.NewClock()

    f1, err := os.Create("/tmp/dat")

	clockJuan.Message("Msg1","Paco", &clockPaco.(logicSystem),f1) //Paco a Juan




	fmt.Println("Hola")
}
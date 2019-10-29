//Jorge Vela Peña	
//Grado en Ingeniería Telemática
package main

import(
	"fmt"
	"time"
	"sync"
)

const (
	ItsFull=0
	AccessTheRoom=1
	CheckWaitingRoom=2
	CutHair=3
	FinCutHair=4
	LeaveWaitingRoom=5
	IsThereClients=6
	ThereIsClients=7
	ThereIsNotClients=8
)

func receptionistFunc(receptionist chan int){
		ClientsWaitRoom :=0
		for{
			select{
				case x := <- receptionist:
					if x == CheckWaitingRoom {
						if ClientsWaitRoom < 5 {
							ClientsWaitRoom++
							receptionist<-AccessTheRoom
						}else{
							receptionist<-ItsFull
						}
					}else if x == IsThereClients{
						if ClientsWaitRoom != 0 {
							receptionist<-ThereIsClients
						}else{
							receptionist<-ThereIsNotClients
						}
					}else if x == LeaveWaitingRoom {
						ClientsWaitRoom--
					}
			}
		}
}

func client (NumClient int, receptionist chan int, WaitingRoom chan chan int ){
		ChanClient := make(chan int)
		select {
			case receptionist <- CheckWaitingRoom:
			action := <-receptionist
			if action == AccessTheRoom{
				fmt.Println("Cliente ",NumClient,": me siento en la sala de espera")
				WaitingRoom <- ChanClient
				x := <- ChanClient
				if x == CutHair {
					receptionist <- LeaveWaitingRoom
					fmt.Println("Cliente ",NumClient,": me corto el pelo")
					y := <- ChanClient
					if y == FinCutHair {
						fmt.Println("Cliente ",NumClient,": termino de cortarme el pelo")
					}
				}
			}else{
				fmt.Println("Cliente", NumClient ,": me voy de la barberia, esta llena")
			}
		}
}

func barber(i int, WaitingRoom chan chan int, receptionist chan int){
	for {
		select {
			case k := <-WaitingRoom:	
				fmt.Println("Barbero ", i, " : Comienzo a cortar el pelo")
				k<-CutHair
				time.Sleep(1 * time.Second)
				k<-FinCutHair
				fmt.Println("Barbero ", i, " : Termino de cortar el pelo")

				receptionist <-IsThereClients
				x :=<-receptionist
				if x == ThereIsNotClients{
					fmt.Println("Me duermo esperando clientes")
				}
		}
	}
}

func main() {
	WaitingRoom := make(chan chan int)
	receptionist := make(chan int)
	var wg sync.WaitGroup

	go func(){ 
		receptionistFunc(receptionist)
	}()

	numBarber:=0
	for numBarber < 2 {
		go func(){ 
			barber(numBarber, WaitingRoom, receptionist)
		}()
		time.Sleep(200 * time.Millisecond)
		numBarber++
	}

	wg.Add(40)
	i:=0
	for i < 40 {
		go func() {
			client(i, receptionist, WaitingRoom)
			wg.Done()
		}()
		time.Sleep(200 * time.Millisecond)
		i++
	}

	wg.Wait()
}

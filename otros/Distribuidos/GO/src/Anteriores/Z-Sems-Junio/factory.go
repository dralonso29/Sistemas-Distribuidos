//Jorge Vela Peña
//Grado en Ingeniería Telemática

package main

import(
	"sem"
	"sync"
	"fmt"
	"time"
)

type ProductProductionLine struct {
	sem *sem.Sem
	mutexProd sync.Mutex
	buffer []int
}

func createProduct(prod *ProductProductionLine, i int){
	prod.mutexProd.Lock()
	prod.sem.Up()
	prod.buffer = append(prod.buffer, i)
	prod.mutexProd.Unlock()
}


func ObtainProduct(prod *ProductProductionLine) int{
	prod.mutexProd.Lock()
	prod.sem.Down()
	fmt.Println("ENTRA")
	ValueId := prod.buffer[0]
	prod.buffer = prod.buffer[1:len(prod.buffer)]
	prod.mutexProd.Unlock()
	return ValueId
}

func robot(NumRobot int, IdScreen int,IdCase int, IdCab [5]int,IdMotherBoard int){
	fmt.Println("robot", NumRobot, "cables", IdCab[0], IdCab[1], IdCab[2], IdCab[3], IdCab[4], "pantalla", IdScreen, "carcasa", IdCase, "placa", IdMotherBoard, "Comenzando")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("robot", NumRobot, "cables", IdCab[0], IdCab[1], IdCab[2], IdCab[3], IdCab[4], "pantalla", IdScreen, "carcasa", IdCase, "placa", IdMotherBoard, "Terminado")
}

func main(){
	fmt.Println("HOLA")
	var wg sync.WaitGroup
		
	var syncScreen, syncCase, syncCable, syncMotherboard sync.Mutex
	var buffScreen, buffCase, buffCable, buffMotherboard []int
	semScreen := sem.NewSem(0)
	semCase:= sem.NewSem(0)
	semCable := sem.NewSem(0)
	semMotherboard := sem.NewSem(0)

	Screen :=ProductProductionLine{semScreen,syncScreen,buffScreen }
	Case :=ProductProductionLine{semCase,syncCase,buffCase }
	Cable :=ProductProductionLine{semCable,syncCable,buffCable }
	Motherboard := ProductProductionLine{semMotherboard,syncMotherboard,buffMotherboard }	


	wg.Add(4)
	go func(){
		i := 0
		for i < 15 {
			createProduct(&Screen, i)
			i++
		}
		wg.Done()
	}()	

	go func() {
		i := 0
		for i < 15 {
			createProduct(&Case, i)
			i++
		}
		wg.Done()

	}()

	go func() {
		i := 0
		for i < 15 {
			createProduct(&Cable, i)
			i++
		}
		wg.Done()
	}()

	go func() {
		i := 0
		for i < 100 {
			createProduct(&Motherboard, i)
			i++
		}
		wg.Done()
	}()


	//go func() {
		IdScreen := ObtainProduct(&Screen)
		fmt.Println(IdScreen)
		/*IdCase := ObtainProduct(&Case)
		IdMotherBoard := ObtainProduct(&Motherboard)
		var IdCable [5]int
		l := 0
		for l < 5 {
			IdCable[l] = ObtainProduct(&Cable)
			l++
		}
		robot(0, IdScreen, IdCase, IdCable,IdMotherBoard)*/
		//wg.Done()
	//}()*/
	wg.Wait()
}


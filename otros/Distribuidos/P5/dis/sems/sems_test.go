//Eva Gordo Calleja
//Grado telemática

package sems

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	fmt.Println("---------------------")
	fmt.Println("-----PRIMER TEST-----")
	fmt.Println("---------------------")
	s := NewSem(1)
	numero := 0
	for i := 0; i < 4; i++ {
		go func() {
			s.Down()
			numero++
			fmt.Println("NUMERO", numero)
			s.Up()
		}()
	}
	time.Sleep(2 * time.Second)
}

func Test2(t *testing.T) {
	fmt.Println("----------------------")
	fmt.Println("-----SEGUNDO TEST-----")
	fmt.Println("----------------------")

	s := NewSem(0)
	for i := 0; i < 100; i++ {
		go func() {
			s.Up()
			fmt.Println("UP", s.numero)
			s.Down()
			fmt.Println("DOWN", s.numero)
		}()
	}
	time.Sleep(2 * time.Second)

	if s.numero != 0 {
		t.Error("Test failed.")
	}
}

func Test3(t *testing.T) {
	fmt.Println("----------------------")
	fmt.Println("-----TERCER TEST------")
	fmt.Println("----------------------")
	s := NewSem(0)
	for i := 0; i < 6; i++ {
		go func() {
			s.Down()
			fmt.Println("DOWN", s.numero)
		}()
	}
	s.Up()
	fmt.Println("UP", s.numero)
	s.Up()
	fmt.Println("UP", s.numero)

	time.Sleep(2 * time.Second)
}

func Test4(t *testing.T) {
	fmt.Println("---------------------")
	fmt.Println("-----CUARTO TEST-----")
	fmt.Println("---------------------")
	s := NewSem(10)
	for i := 0; i < 50; i++ {
		go func(n int) {
			s.Down()
			fmt.Println("DOWN", n)
		}(i)
	}
	time.Sleep(2 * time.Second)
}

func Test5(t *testing.T) {
	fmt.Println("---------------------")
	fmt.Println("-----QUINTO TEST-----")
	fmt.Println("---------------------")

	mutexesc := NewSem(1)
	mutexnl := NewSem(1)
	torn := NewSem(1)
	var nl int
	var array []int
	s1 := rand.NewSource(42)
	contador := rand.New(s1)

	for i := 0; i < 10; i++ {
		go func() {
			torn.Down()
			mutexesc.Down()
			torn.Up()
			array = append(array, contador.Intn(100))
			time.Sleep(1 * time.Millisecond)
			mutexesc.Up()
		}()
	}
	for i := 0; i < 15; i++ {
		go func() {
			torn.Down()
			torn.Up()
			mutexnl.Down()
			nl++
			if nl == 1 {
				mutexesc.Down()
			}
			mutexnl.Up()
			fmt.Println(array)
			time.Sleep(2 * time.Millisecond)
			mutexnl.Down()
			nl--
			if nl == 0 {
				mutexesc.Up()
			}
			mutexnl.Up()
		}()
	}
	time.Sleep(2 * time.Second)
}

func Test6(t *testing.T) {
	fmt.Println("---------------------")
	fmt.Println("-----SEXTO TEST------")
	fmt.Println("---------------------")

	N := 10
	ticketsem := NewSem(0)
	holesem := NewSem(N)
	buf := make([]int, N)
	ticket := 0

	go func() {
		i := 0
		for r := 0; r < 10; r++ {
			ticket++
			fmt.Println("Produzco: ", ticket)
			holesem.Down()
			buf[i] = ticket
			i = (i + 1) % N
			ticketsem.Up()
		}
	}()
	go func() {
		j := 0
		for r := 0; r < 10; r++ {
			ticketsem.Down()
			ticket := buf[j]
			j = (j + 1) % N
			holesem.Up()
			fmt.Println("Consumo: ", ticket)
			for i, v := range buf {
				if v == ticket {
					buf[i] = buf[len(buf)-1]
				}
			}
		}
	}()
	time.Sleep(2 * time.Second)
}


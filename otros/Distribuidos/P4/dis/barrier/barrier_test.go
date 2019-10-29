//Eva Gordo Calleja
//Grado telemática

package barrier

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	fmt.Println("---------------------")
	fmt.Println("-----PRIMER TEST-----")
	fmt.Println("---------------------")
	var b = NewBarrier(3)
	start := time.Now()
	fmt.Println("Hora inicial: ", start)
	time.Sleep(3 * time.Second)
	go func() {
		b.Wait()
		go1 := time.Now()
		fmt.Println(go1, "sale 1")
	}()

	go func() {
		b.Wait()
		go2 := time.Now()
		fmt.Println(go2, "sale 2")
	}()

	go func() {
		b.Wait()
		go3 := time.Now()
		fmt.Println(go3, "sale 3")
	}()

	go func() {
		b.Wait()
		go4 := time.Now()
		fmt.Println(go4, "sale 4")
	}()

	go func() {
		b.Wait()
		go5 := time.Now()
		fmt.Println(go5, "sale 5")
	}()
	go func() {
		b.Wait()
		go6 := time.Now()
		fmt.Println(go6, "sale 6")
	}()
	go func() {
		b.Wait()
		go7 := time.Now()
		fmt.Println(go7, "sale 7")
	}()

	end := time.Now()
	fmt.Println("Hora final: ", end)

	time.Sleep(3 * time.Second)
	diff := end.Sub(start)

	fmt.Println("Han pasado: ", diff)

}

func Test2(t *testing.T) {
	fmt.Println("---------------------")
	fmt.Println("-----SEGUNDO TEST----")
	fmt.Println("---------------------")
	var b = NewBarrier(3)
	for i := 0; i <= 7; i++ {
		time.Sleep(2 * time.Second)
		go func() {
			fmt.Println(b.contador, "CONTADOR antes wait")
			b.Wait()
			fmt.Println(b.contador, "-----CONTADOR despues wait")
		}()
	}
	time.Sleep(2 * time.Second)
}

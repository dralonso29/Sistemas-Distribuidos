//Jorge Vela Peña
//Grado en Ingeniería Telemática

package rendez

import (
	"sync"
	"testing"
)

func TestSec(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			x := Rendezvous(1, 5)
            wg.Done()
			if x!=5 {
				t.Error("Test failed. Value",x," is not correct")
			}

		}()
	}

	wg.Add(1)
	go func() {
		x:=Rendezvous(4, "g")
		wg.Done()
		if x!="h" {
			t.Error("Test failed. Value",x," is not correct")
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			x := Rendezvous(2, 2)
	        wg.Done()
			if x!=2 {
			    wg.Done()
				t.Error("Test failed. Value",x," is not correct")
			}

		}()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			x:= Rendezvous(3, "s")
            wg.Done()
			if x!="s"{
				t.Error("Test failed. Value",x," is not correct")
			}
		}()
	}

	wg.Add(1)
	go func() {
		x:=Rendezvous(4, "h")
		wg.Done()
		if x!="g" {
			t.Error("Test failed. Value",x," is not correct")
		}
	}()
	wg.Wait()

}







//Jorge Vela Peña
//Grado en Ingeniería Telemática

package rendez

import (
	"sync"
	"testing"
)

func TestFirst(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		s := Rendezvous(1, 2)
		wg.Done()
		if s != 1 {
			t.Error("Test failed")
		}
	}()

	wg.Add(1)
	go func() {
		s := Rendezvous(1, 1)
		wg.Done()
		if s != 2 {
			t.Error("Test failed")
		}
	}()

	wg.Add(1)
	go func() {
		s := Rendezvous(2, 4)
		wg.Done()
		if s != 5 {
			t.Error("Test failed")
		}
	}()

	wg.Add(1)
	go func() {
		s := Rendezvous(3, 1)
		wg.Done()
		if s != 7 {
			t.Error("Test failed")
		}
	}()

	wg.Add(1)
	go func() {
		s := Rendezvous(4, 2)
		wg.Done()
		if s != 6 {
			t.Error("Test failed")
		}
	}()

	wg.Add(1)
	go func() {
		s := Rendezvous(3, 7)
		wg.Done()
		if s != 1 {
			t.Error("Test failed")
		}
	}()

	wg.Add(1)
	go func() {
		s := Rendezvous(2, 5)
		wg.Done()
		if s != 4 {
			t.Error("Test failed")
		}
	}()

	wg.Add(1)
	go func() {
		s := Rendezvous(4, 6)
		wg.Done()
		if s != 2 {
			t.Error("Test failed")
		}
	}()
	wg.Wait()
}

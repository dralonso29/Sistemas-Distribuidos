//Jorge Vela Peña
//Grado en Ingeniería Telemática

package sem

import(
	"sync"
	"testing"
)

func TestFirst(t *testing.T) {
	var wg sync.WaitGroup
	semTest := NewSem(0)

	numDown :=0
	for numDown < 4 {
		wg.Add(1)
		go func() {
			semTest.Down()
			wg.Done()
		}()
		numDown++
	}

	wg.Add(1)
	go func() {
		numUp :=0
		for numUp < 4 {
			semTest.Up()
			numUp++
		}
		wg.Done()
	}()

	wg.Wait()

	if semTest.numeroMax != 0{
		t.Error("Test failed.")
	}
}

func TestSecond(t *testing.T) {
	numUps :=200000
	numWaits := numUps + 1
	var wg sync.WaitGroup
	semTest := NewSem(0)
	
	wg.Add(numWaits)
	go func() {
		numUp :=0
		for numUp < numUps {
			semTest.Up()
			numUp++
		}
		wg.Done()
	}()

	numDown :=0
	for numDown < numUps {
		numDown++
		go func() {
			semTest.Down()
			wg.Done()
		}()
	}
	wg.Wait()

	if semTest.numeroMax != 0{
		t.Error("Test failed.")
	}
}



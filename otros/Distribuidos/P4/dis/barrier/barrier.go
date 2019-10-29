//Eva Gordo Calleja
//Grado telemática

package barrier

import (
	"sync"
)

type Barrier struct {
	numero   int
	contador int
	muxN     *sync.Mutex
	waitGr   *sync.WaitGroup
}

func NewBarrier(n int) *Barrier {
	muxN := &sync.Mutex{}
	waitGr := &sync.WaitGroup{}
	waitGr.Add(1)
	contador := 0
	return &Barrier{n, contador, muxN, waitGr}
}

func (b *Barrier) Wait() {
	b.muxN.Lock()
	b.contador++
	newWaitGr := b.waitGr
	if b.contador == b.numero {
		b.waitGr.Done()
		b.contador = 0
		b.waitGr = &sync.WaitGroup{}
		b.waitGr.Add(1)
	}
	b.muxN.Unlock()
	newWaitGr.Wait()
}

//Eva Gordo Calleja
//Grado telemática

package sems

import (
	"container/list"
	"sync"
)

type Sem struct {
	numero int
	muxN   *sync.Mutex
	list   *list.List
}

func NewSem(n int) *Sem {
	if n < 0{
		panic("No puede inicializar a negativo")
	}
	muxN := &sync.Mutex{}
	lista := list.New()
	return &Sem{n, muxN, lista}
}

func (s *Sem) Down() {
	s.muxN.Lock()
	if s.numero == 0 {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		s.list.PushBack(wg)
		s.muxN.Unlock()
		wg.Wait()
		return
	}
	s.numero--
	s.muxN.Unlock()
}


func (s *Sem) Up() {
	s.muxN.Lock()
	if s.list.Len() > 0 {
		wg := &sync.WaitGroup{}
		wait := s.list.Front()
		wg = wait.Value.(*sync.WaitGroup)
		s.list.Remove(wait)
		wg.Done()
	}
	s.numero++
	s.muxN.Unlock()
}

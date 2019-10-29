//Jorge Vela Peña
//Grado en Ingeniería Telemática.

package sem

import(
	"sync"
)

type UpDowner interface{
	Up()
	Down()
}

type Sem struct {
	c *sync.Cond
	numeroMax int
}


func NewSem(ntok int) *Sem{
	var m sync.Mutex
	c := sync.NewCond(&m)
	return &Sem{
		c,
		ntok,
	}
}

func (s *Sem) Up(){
	s.c.L.Lock()
	s.numeroMax++
	s.c.Signal()
	s.c.L.Unlock()
} 


func (s *Sem) Down(){
	s.c.L.Lock()
	if s.numeroMax == 0 {
		s.c.Wait()
	}
	s.numeroMax--
	s.c.L.Unlock()
}

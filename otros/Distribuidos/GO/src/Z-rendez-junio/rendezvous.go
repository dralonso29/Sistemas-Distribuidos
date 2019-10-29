//Jorge Vela Peña
//Grado en Ingeniería Telemática

package rendez


import (
	"sync"
)

type countWords struct {
	wg  *sync.WaitGroup
	val interface{}
}

var syncMap struct {
	sync.Mutex
	s map[int]*countWords
}

func Rendezvous(tag int, val interface{}) interface{} {
	var wg sync.WaitGroup
	var finalP *countWords

	syncMap.Lock()
	if(syncMap.s == nil){
		syncMap.s = make(map[int]*countWords)
	}

	if(syncMap.s[tag] == nil){
		syncMap.s[tag] = &countWords{&wg, val}
		wg.Add(1)
		finalP = syncMap.s[tag]
		syncMap.Unlock()
		wg.Wait();
		return finalP.val
	}
	sp := *syncMap.s[tag]
	finalP = &sp
	syncMap.s[tag].val=val
	syncMap.s[tag].wg.Done()
	delete(syncMap.s, tag)
	syncMap.Unlock()
	return finalP.val 	
}




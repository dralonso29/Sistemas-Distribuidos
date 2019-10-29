package rendez

import (
	"fmt"
	"sync"
)

type countWords struct {
    //W string
    //wg sync.WaitGroup
    val  interface{}
}

var m map[int]countWords 


func printMap(){
	fmt.Println(m)
}


func Rendezvous(tag int, val interface{}) interface{}{
	//var wg sync.WaitGroup
	if (m == nil){
		m = make(map[int]countWords)
	}
	m[tag] = countWords{val}		
	return nil
}





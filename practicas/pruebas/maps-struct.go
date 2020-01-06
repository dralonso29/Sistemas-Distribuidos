package main

import (
	"fmt"
	"time"
	"math/rand"
)

const NMATCHES = 1		// total number of matches
const MTIME	= 20000		// time of a match in ms
const CDTIME = 750		// countdown time in ms
const ENDMATCH = 5		// playtime of a match is 20 secs (10 chances * 2 seconds between chances = 20 seconds)
const TCHANCE = 2000	// time between chances in ms

type mdata struct {		// all match's info about all other matches and itself
	localgoals int		// register of local goals
	visitorgoals int	// register of visitor goals
	invalid chan int 	// in this case int is match identifier
	mustfollow bool		// if a match must be followed or not
}

var invalid chan int

func fillMap(m map[int]mdata)  {
	fmt.Println("antes de modificar el valor data[0].mustfollow en fillMap: ", m[0].mustfollow)
	m[0] = mdata{1,2,invalid, true}
}

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	invalid = make(chan int, NMATCHES)
	data := make(map[int]mdata)
	data[0] = mdata{1,2,invalid, false}
	data[1] = mdata{0,3,invalid, true}
	fmt.Println("Antes En el main data[0].mustfollow = ",data[0].mustfollow)
	fillMap(data)
	fmt.Println("Despues En el main data[0].mustfollow = ",data[0].mustfollow)
	fmt.Println("Termina el programa")
}
//!-main

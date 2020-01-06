package main

import (
	"fmt"
	"time"
	"math/rand"
)

const NMATCHES = 3		// total number of matches
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

type mcomm struct {
	id int
	msg string
}

var matches []chan mcomm

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	matches := make([]chan mcomm, NMATCHES)
	for i := 0; i < NMATCHES; i++ {
		//n.Add(1)
		ch := make(chan mcomm, 0)
		matches[i] = ch
	}
	fmt.Println(matches)
	fmt.Println("Termina el programa")
}
//!-main

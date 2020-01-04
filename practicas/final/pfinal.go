package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

const NMATCHES = 2				// total number of matches
const MTIME	= 20000				// time of a match in ms
const CDTIME = 750				// countdown time in ms
const CHANCES = 1				// playtime of a match is 20 secs (10 chances * 2 seconds between chances = 20 seconds)
const TCHANCE = 2000			// time between chances in ms

type mcomm struct {
	matchid int
	msg string
}

type mdata struct {				// all match's info about all other matches and itself
	sync.Mutex
	localscores int				// register of local goals
	visitorscores int			// register of visitor goals
	mch chan mcomm 				// match's channel
	followedby [NMATCHES]bool	// array of matches that follows me
	scoreinvalid [NMATCHES]bool	// if a match has an invalid scoreboard
}

var matches [NMATCHES]chan mcomm

//!+myUnlock
func myUnlock(ch chan struct {}, n int)  {
	for i := 0; i < n; i++ {
		//fmt.Println("Unlock: ",i)
		<- ch
	}
}
//!-myUnlock

//!+myLock
func myLock(ch chan struct {}, n int)  {
	for i := 0; i < n; i++ {
		//fmt.Println("Lock: ",i)
		ch <- struct{}{}
	}
}
//!-myLock

//!+chooseTeam
func chooseTeam()  string{
	prob := 50
	n := rand.Intn(100)
	if n <= prob {
		return "local"
	}
	return "visitor"
}
//!-chooseTeam

//!+isGoal
func isGoal()  bool{
	prob := 0
	n := rand.Intn(100)
	return n > prob
}
//!-isGoal

//!+matchHandler
func matchHandler(id int, data *mdata)  {
	for {
		select {
		case mid := <- matches[id]:
			// buscar en el array que se le ha pasado a matchHandler e invalidar el resultado
			// Cada vez que se invalida un resultado del array habria que usar un mutex(?)
			fmt.Println("E",id," Handler: msg de ",mid.matchid, " es ",mid.msg)
			data.Lock()
			data.scoreinvalid[mid.matchid] = true
			data.Unlock()
		}
	}
}
//!-matchHandler

//!+fillMatchData
func fillMatchData(id int)  *mdata{
	data := new(mdata)
	data.localscores = 0
	data.visitorscores = 0
	data.mch = matches[id]
	data.scoreinvalid = [NMATCHES]bool{false, false}//, false, false}
	if id == 0 {
		data.followedby = [NMATCHES]bool{false, true} // test
		// data.followedby = [NMATCHES]bool{false, true, false, false} //good one
	}else if id == 1 {
		data.followedby = [NMATCHES]bool{true, false} // test
		// data.followedby = [NMATCHES]bool{true, false, true, false} // good one
	}//else if id == 2 {
	// 	data.followedby = [NMATCHES]bool{true, true, false, false}
	// }else{
	// 	data.followedby = [NMATCHES]bool{false, false, true, false}
	// }
	return data
}
//!-fillMatchData

//!+sendMessage
func sendMessage(id int, msg string)  {
	for i := 0; i < NMATCHES; i++ {
		if id == i {
			continue
		}
		matches[i] <- mcomm{id, msg}
	}
}
//!-sendMessage

//!+match
func match(start chan struct{}, id int, n *sync.WaitGroup)  {
	defer n.Done()
	start <- struct{}{}
	data := fillMatchData(id)
	go matchHandler(id, data)
	//fmt.Println("Match", id, ": ", data)
	for i := 1; i <= CHANCES; i++ {
		team := chooseTeam()
		isgoal := isGoal()
		fmt.Println("E",id, ": Ocasion del ", team, ", goal = ", isgoal)
		if isgoal {
			fmt.Println("E",id," invalida al resto de partidos")
			data.Lock()
			if team == "local" {
				data.localscores++
			}else{
				data.visitorscores++
			}
			data.Unlock()
			go sendMessage(id, "invalid") 	// bucle for que envia por un canal buffered llamado invalidate que reciba el id del partido que tiene
											// un resultado invalido
		}
		time.Sleep(TCHANCE * time.Millisecond)
	}
	data.Lock()
	fmt.Println("E",id,": data: ", data)
	data.Unlock()
	// fmt.Println("Partido ", id, ": ", time.Now().String())
}
//!-match

//!+startCountdown
func startCountdown()  {
	fmt.Println(NMATCHES, " matches start in...")
	for i := 3; i >= 1 ; i-- {
		fmt.Printf("\r%d",i)
		time.Sleep(CDTIME * time.Millisecond)
	}
	fmt.Printf("\rGO!\n")
}
//!-startCountdown

//!+matchesGenerator
func matchesGenerator()  {
	var n sync.WaitGroup
	defer n.Wait()
	start := make(chan struct{}, NMATCHES)
	myLock(start, NMATCHES)
	for i := 0; i < NMATCHES; i++ {
		n.Add(1)
		matches[i] = make(chan mcomm, 0)
		go match(start, i, &n)
	}
	startCountdown()
	//fmt.Println("Channels: ",matches)
	myUnlock(start, NMATCHES)
	// fmt.Println(matches[0])
	// fmt.Println(matches[1])
	// fmt.Println(matches[2])
	// fmt.Println(matches[3])

}
//!-matchesGenerator

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	matchesGenerator()
	fmt.Println("Termina el programa")
}
//!-main

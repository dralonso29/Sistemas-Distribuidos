package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

const NMATCHES = 2				// total number of matches
const NCENTRAL = 1				// total central systems
const MTIME	= 20000				// time of a match in ms
const STIME = 5000				// time between checking stadiums in ms
const CDTIME = 750				// countdown time in ms
const CHANCES = 1				// playtime of a match is 20 secs (10 chances * 2 seconds between chances = 20 seconds)
const TCHANCE = 2000			// time between chances in ms
const PMISS = 70				// prob of goal's failure in %

type mcomm struct {				// protocol between stadiums and central
	matchid int					// statium id (centralid = NMATCHES)
	msg string					// type of message (invalid, central, update)
}

type omatch struct {			// info from other matches
	olocalscore int				// other local score
	ovisitorscore int			// other visitor score
	invalid bool				// if we have an invalid copy of match's score
	isfollowed bool				// if this match is followed by me
}

type mdata struct {				// all match's info about all other matches and itself
	sync.Mutex
	localscore int				// register of local goals
	visitorscore int			// register of visitor goals
	mch chan mcomm 				// match's channel
	otherscores []omatch		// array of other matches' data. Position of array corresponds with matches' ids
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
	n := rand.Intn(100)
	return n > PMISS
}
//!-isGoal

//!+matchHandler
func matchHandler(myid int, data *mdata)  {
	for {
		select {
		case id := <- matches[myid]:
			fmt.Println("E",myid," Handler: msg de ",id.matchid, " es ",id.msg)
			data.Lock()
			if id.msg == "invalid" {
				if data.otherscores[id.matchid].isfollowed {
					data.otherscores[id.matchid].invalid = true
				}
			}

			data.Unlock()
		}
	}
}
//!-matchHandler

//!+initOtherScores
func initOtherScores(data *mdata)  {
	data.otherscores = []omatch{}
	for i := 0; i < NMATCHES; i++ {
		n := omatch{olocalscore: 0, ovisitorscore: 0, invalid: false, isfollowed: false}
		data.otherscores = append(data.otherscores, n)
	}
}
//!-initOtherScores

//!+fillMatchData
func fillMatchData(id int)  *mdata{
	data := new(mdata)
	data.localscore = 0
	data.visitorscore = 0
	data.mch = matches[id]
	initOtherScores(data)
	if id == 0 {
		data.otherscores[1].isfollowed = true		//TEST
		//data.otherscores[2].isfollowed = true 	// OK
	}else if id == 1 {
		data.otherscores[0].isfollowed = true		// TEST
		// data.otherscores[0].isfollowed = true 	// OK
		// data.otherscores[2].isfollowed = true	//OK
	}
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
				data.localscore++
			}else{
				data.visitorscore++
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

//!+chooseStadiumsOrder
func chooseStadiumsOrder()  [NMATCHES]int{
	p := [NMATCHES]int{0,1}	// TEST
	a := [NMATCHES]int{}
	// p := [NMATCHES]int{0,1,2,3} // OK
	ind := rand.Intn(NMATCHES)
	l := len(p)
	for i := ind; i < l+ind; i++ {
		//fmt.Println(p[i%l], ", i%l: ",i%l)
		a[i-ind] = p[i%l]
	}
	return a
}
//!-chooseStadiumsOrder

//!+centralSystem
func centralSystem(start chan struct{}, id int, n *sync.WaitGroup)  {
	defer n.Done()
	start <- struct{}{}
	fmt.Println("Lanzamos al sistema central")
	a := [NMATCHES]int{}
	for i := 0; i < CHANCES; i++ {
		a = chooseStadiumsOrder()
	}
	fmt.Println(a)
}
//!-centralSystem

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
	l := NMATCHES+NCENTRAL
	start := make(chan struct{}, l)
	myLock(start, l)
	for i := 0; i < NMATCHES; i++ {
		n.Add(1)
		matches[i] = make(chan mcomm, 0)
		go match(start, i, &n)
	}
	sysid := NMATCHES
	n.Add(1)
	go centralSystem(start, sysid, &n)
	startCountdown()
	myUnlock(start, l)
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

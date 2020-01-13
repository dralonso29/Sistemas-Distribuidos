package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
)

const NMATCHES = 2				// total number of matches
const NCENTRAL = 1				// total central systems
const NTOT = NMATCHES + NCENTRAL// suma de partidos y el sistema central
const TGAME	= 10				// time of a match in s
const STIME = 5					// time between checking stadiums in s
const CDTIME = 750				// countdown time in ms
const CHANCES = 10				// playtime of a match is 20 secs (10 chances * 2 seconds between chances = 20 seconds)
const TCHANCE = 2				// time between chances in s
const PMISS = 0				// prob of goal's failure in % ######### CAMBIAR ESTO
const SECOND = 1				// just one second
const CSCHANS = 2				// number of central systems's channels (invalid, getinfo)

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

type csdata struct {			// info from central system
	sync.Mutex
	csch chan mcomm 			// central system's channel
	matches [NMATCHES]mdata		// array of matches' data. Position of array corresponds with matches' ids
	invalid [NMATCHES]bool		// to know if central's system copy of match's data is invalid or not
	writeback int				// number of writeback operations in central system
}

type controller struct {
	sync.Mutex
	chance int
	matches int
}

type mcomm struct {				// protocol between stadiums and central
	id int						// statium id (centralid = NMATCHES)
	msg string					// type of message (invalid, now, central, update)
	round int					// cental system must know if is last round to close all channels
	matchdata *mdata			// header that contains a match info
}

var mchans [NMATCHES]chan mcomm
var cschans [CSCHANS]chan mcomm
var control controller

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
	fmt.Printf("%T\n",data)
	data.localscore = 0
	data.visitorscore = 0
	data.mch = mchans[id]
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
func sendMessage(id int, msg string, round int)  {
	fmt.Println("E",id," Envia mensaje de invalidar, ronda", round)
	for i := 0; i < NMATCHES; i++ {
		if id == i {
			continue
		}
		mchans[i] <- mcomm{id, msg, round, &mdata{}}
	}
	cschans[0] <- mcomm{id, msg, round, &mdata{}}
}
//!-sendMessage

//!+isSecondCase
func isSecondCase(i int)  bool{
	return ((i % STIME == 0) && (i/STIME)%2 != 0) // returns true if is five multiple and (i/5)%2 is odd (i = 5,15,25,35,45...)
}
//!-isSecondCase

//!+mustWarnCentralSystem
func mustWarnCentralSystem(i int)  bool{
	return control.chance%STIME == 0 || isSecondCase(i)
}
//!-mustWarnCentralSystem

//!+isOdd
func isOdd(i int)  bool{
	return i%2 != 0
}
//!-isOdd

//!+callCentralSystem
func callCentralSystem(id int, i int)  {
	control.matches++
	fmt.Println("E",id," i:",i,//"mustWarnCentralSystem: ",mustWarnCentralSystem(),", control.matches: ", control.matches,
					", control.chance: ", control.chance)

	if (mustWarnCentralSystem(i)) && control.matches == NMATCHES {
		// fmt.Println("Dentro de mustWarnCentralSystem: Avisamos al SC")
		cschans[0] <- mcomm{-1, "now", i, &mdata{}}
	}

	if control.matches == NMATCHES {
		if !isSecondCase(i) {
			control.chance++
		}
		control.matches = 0
	}
}
//!-callCentralSystem

//!+matchHandler
func matchHandler(myid int, data *mdata, cansend chan struct{})  {
	for {
		select {
		case info := <- mchans[myid]:
			fmt.Println("E",myid," Handler: msg de ",info.id, " es ",info.msg)
			time.Sleep(time.Second)
			if info.msg == "invalid" {
				fmt.Println("E",myid," le llega un invalid de E",info.id)
				if data.otherscores[info.id].isfollowed {
					data.Lock()
					data.otherscores[info.id].invalid = true
					data.Unlock()
					mchans[info.id] <- mcomm{myid, "updateme", -1, &mdata{}}
				}
			}else if info.msg == "getinfo" {
				cschans[1] <- mcomm{myid, "sendinfo", -1, data}
			}else if info.msg == "updateme" {
				fmt.Println("E",myid," le hace un updateme a E",info.id)
				mchans[info.id] <- mcomm{myid, "updated", -1, &mdata{localscore:data.localscore, visitorscore:data.visitorscore}}
			}else if info.msg == "updated" {
				fmt.Println("A E",myid," le llega un updated de E",info.id)
				data.Lock()
				data.otherscores[info.id].olocalscore = info.matchdata.localscore
				data.otherscores[info.id].ovisitorscore = info.matchdata.visitorscore
				data.otherscores[info.id].invalid = false
				data.Unlock()
				<- cansend
			}
		}
	}
}
//!-matchHandler

//!+match
func match(start chan struct{}, id int, n *sync.WaitGroup, cansend chan struct{}, close chan struct{})  {
	defer n.Done()
	var n2 sync.WaitGroup
	start <- struct{}{}
	data := fillMatchData(id)
	n2.Add(1)
	go matchHandler(id, data, cansend)
	//fmt.Println("Match", id, ": ", data)
	chance := 1
	for i := 1; i <= TGAME; i++ {
		//fmt.Println("E",id, ", i: ",i)
		if isOdd(i) {
			if isSecondCase(i) {
				//fmt.Println("E",id, ", i: ",i," avisaria al sistema central (dentro if)")
				control.Lock()
				callCentralSystem(id,i)
				control.Unlock()
			}
			time.Sleep(SECOND *time.Second)
			continue
		}
		team := chooseTeam()
		isgoal := isGoal()
		fmt.Println("E",id, ": Ocasion ",chance," del ", team, ", goal = ", isgoal)
		chance++
		if isgoal {
			//fmt.Println("E",id," invalida al resto de partidos")
			data.Lock()
			if team == "local" {
				data.localscore++
			}else{
				data.visitorscore++
			}
			fmt.Println("E",id, ": resultado ", data.localscore,"-",data.visitorscore)
			data.Unlock()
			cansend <- struct{}{}
			sendMessage(id, "invalid",i)
		}
		//fmt.Println("E",id, ", i: ",i," avisaria al sistema central")
		control.Lock()
		callCentralSystem(id, i)
		control.Unlock()
		// time.Sleep(TCHANCE * time.Second)
		time.Sleep(SECOND*time.Second)
	}
	data.Lock()
	fmt.Println("E",id,": data: ", data)
	data.Unlock()
	n2.Wait()
	close <- struct{}{}
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

func fillCentralSystemData(id int)  *csdata {
	csd := new(csdata)
	for i := 0; i < NMATCHES; i++ {
		// csd.csch = mchans[id]
		csd.matches[i].localscore = 0
		csd.matches[i].visitorscore = 0
		csd.invalid[i] = false
	}
	csd.writeback = 0
	return csd
}

//!+centralSystem
func centralSystem(start chan struct{}, id int, n *sync.WaitGroup)  {
	defer n.Done()
	start <- struct{}{}
	fmt.Println("Lanzamos al sistema central")
	csd := fillCentralSystemData(id)
	for {
		select {
		case info := <- cschans[0]:
			//fmt.Println("E",myid," Handler: msg de ",id.matchid, " es ",id.msg)
			if info.msg == "invalid" {
				csd.Lock()
				csd.invalid[info.id] = true
				//fmt.Println("SISTEMA CENTRAL: invalidado E" ,info.id,": ",csd.invalid[info.id])
				csd.Unlock()
			}else{
				fmt.Println("SC: ",info)
				storder := chooseStadiumsOrder()
				fmt.Println("El orden elegido de consulta es ", storder)
				for i := 0; i < len(storder); i++ {
					matchid := storder[i]
					mchans[matchid] <- mcomm{id, "getinfo", -1, &mdata{}}
					x := <-cschans[1]
					fmt.Println("SC: E",storder[i],": localscore = ",x.matchdata.localscore,", visitorscore = ",x.matchdata.visitorscore)
					csd.Lock()
					//fmt.Println("SC: E",storder[i],": localscore = ",x.matchdata.localscore,", visitorscore = ",x.matchdata.visitorscore)
					if csd.invalid[matchid] {
						csd.writeback++
						x.matchdata.Lock()
						csd.matches[matchid].localscore = x.matchdata.localscore
						csd.matches[matchid].visitorscore = x.matchdata.visitorscore
						//csd.matches[matchid].otherscores = x.matchdata.otherscores
						x.matchdata.Unlock()
						csd.invalid[matchid] = false
					}
					csd.Unlock()
				}

				if info.round == TGAME {
					fmt.Println("SC : data: ", csd)
					fmt.Println("SC : writeback = ", csd.writeback)
					return
				}
			}
		}
	}
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

//!+initController
func initController()  {
	control.chance = 1
	control.matches = 0
}
//!-initController

//!+closeAllChan
func closeAllChan(cl chan struct{})  {
	c := 0
	for  {
		<- cl
		c++
		if c == NMATCHES {
			for i := 0; i < NMATCHES; i++ {
				close(mchans[i])
			}
			break
		}
	}


}
//!-closeAllChan

//!+matchesGenerator
func matchesGenerator()  {
	var n sync.WaitGroup
	defer n.Wait()
	initController()
	start := make(chan struct{}, NTOT)
	cansend := make(chan struct{}, 1)
	close := make(chan struct{}, 0)
	myLock(start, NTOT)
	for i := 0; i < NMATCHES; i++ {
		n.Add(1)
		mchans[i] = make(chan mcomm, 1)
		go match(start, i, &n, cansend, close)
	}

	n.Add(1)
	csid := NMATCHES
	cschans[0] = make(chan mcomm, 0)
	cschans[1] = make(chan mcomm, 0)
	go centralSystem(start, csid, &n)
	// startCountdown()
	myUnlock(start, NTOT)
	n.Add(1)
	go closeAllChan(close)
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

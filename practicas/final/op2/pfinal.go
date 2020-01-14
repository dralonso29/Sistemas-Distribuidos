package main

import (
	"fmt"
	"time"
	"math/rand"
	"sync"
	"os"
	"io/ioutil"
	"strconv"
	"strings"
)

const NMATCHES = 4				// total number of matches
const NCENTRAL = 1				// total central systems
const NTOT = NMATCHES + NCENTRAL// suma de partidos y el sistema central
const TGAME	= 20				// time of a match in s
const STIME = 5					// time between checking stadiums in s
const CDTIME = 750				// countdown time in ms
const CHANCES = 10				// playtime of a match is 20 secs (10 chances * 2 seconds between chances = 20 seconds)
const TCHANCE = 2				// time between chances in s
const PMISS = 70				// prob of goal's failure in % ######### CAMBIAR ESTO
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

type stadistics struct {
	sync.Mutex
	ILCD int					// invalidar otra linea de una cache diferente
	ALCD int					// actualizar esa linea en otra cache diferente
	AMC int						// actualizar la memoria principal con el nuevo valor de una linea
}

var mchans [NMATCHES]chan mcomm
var cschans [CSCHANS]chan mcomm
var control controller
var stats stadistics

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
	data.localscore = 0
	data.visitorscore = 0
	data.mch = mchans[id]
	initOtherScores(data)
	if id == 0 {
		// data.otherscores[1].isfollowed = true		//TEST
		data.otherscores[2].isfollowed = true 	// OK
	}else if id == 1 {
		// data.otherscores[0].isfollowed = true		// TEST
		data.otherscores[0].isfollowed = true 	// OK
		data.otherscores[2].isfollowed = true	//OK
	}else if id == 2 {
		// data.otherscores[0].isfollowed = true		// TEST
		data.otherscores[3].isfollowed = true	//OK
	}else if id == 3 {
		data.otherscores[1].isfollowed = true 	// OK
	}
	return data
}
//!-fillMatchData

//!+sendMessage
func sendMessage(id int, msg string, round int)  {
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

//!+callCentralSystem
func callCentralSystem(id int, i int)  {
	control.matches++
	// fmt.Println("E",id," i:",i,//"mustWarnCentralSystem: ",mustWarnCentralSystem(),", control.matches: ", control.matches,
	// 				", control.chance: ", control.chance)

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

//!+isOdd
func isOdd(i int)  bool{
	return i%2 != 0
}
//!-isOdd

//!+showInfo
func showInfo(matchinfo *mdata, csinfo *csdata, ismatchinfo bool, id int)  {
	if ismatchinfo {
		fmt.Println("======================================================")
		fmt.Println("Stadium ",id,": ",matchinfo.localscore,"-",matchinfo.visitorscore)
		fmt.Println("Stadiums followed:")
		for i := 0; i < NMATCHES; i++ {
			if matchinfo.otherscores[i].isfollowed {
				fmt.Println("	Stadium ",i,": ",matchinfo.otherscores[i].olocalscore,"-",matchinfo.otherscores[i].ovisitorscore,
							"--> ", matchinfo.otherscores[i].invalid)
			}
		}
		fmt.Println("======================================================")
	}else{
		fmt.Println("======================================================")
		fmt.Println("Central system info")
		for i := 0; i < NMATCHES; i++ {
			fmt.Println("Stadium ", i,": ", csinfo.matches[i].localscore,"-",csinfo.matches[i].visitorscore)
			fmt.Println("Stadiums followed:")
			for j := 0; j < NMATCHES; j++ {
				if csinfo.matches[i].otherscores[j].isfollowed {
					fmt.Println("	Stadium ",j,": ",csinfo.matches[i].otherscores[j].olocalscore,"-",csinfo.matches[i].otherscores[j].ovisitorscore,
								"--> ", csinfo.matches[i].otherscores[j].invalid)
				}

			}
			fmt.Println("------------------------------------------------------")
		}
		fmt.Println("======================================================")
	}
}
//!-showInfo

//!+matchHandler
func matchHandler(myid int, data *mdata)  {
	for {
		select {
		case info := <- mchans[myid]:
			//fmt.Println("E",myid," Handler: msg de ",id.matchid, " es ",id.msg)
			data.Lock()
			if info.msg == "invalid" {
				if data.otherscores[info.id].isfollowed {
					stats.Lock()
					stats.ILCD++
					stats.Unlock()
					data.otherscores[info.id].invalid = true
				}
			}else if info.msg == "getinfo" {
				cschans[1] <- mcomm{myid, "sendinfo", -1, data}
			}else if info.msg == "updatecs" {
				cschans[1] <- mcomm{myid, "sendinfo", -1, &mdata{localscore:data.localscore,
					 											visitorscore: data.visitorscore}}
			}
			data.Unlock()
		}
	}
}
//!-matchHandler

//!+match
func match(start chan struct{}, id int, n *sync.WaitGroup)  {
	defer n.Done()
	start <- struct{}{}
	data := fillMatchData(id)
	go matchHandler(id, data)
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
			sendMessage(id, "invalid",i) //le quito el go delante de esto
		}
		//fmt.Println("E",id, ", i: ",i," avisaria al sistema central")
		control.Lock()
		callCentralSystem(id, i)
		control.Unlock()
		// time.Sleep(TCHANCE * time.Second)
		time.Sleep(SECOND*time.Second)
	}
	data.Lock()
	//fmt.Println("E",id,": data: ", data)
	showInfo(data, &csdata{}, true, id)
	data.Unlock()
	// fmt.Println("Partido ", id, ": ", time.Now().String())
}
//!-match

//!+chooseStadiumsOrder
func chooseStadiumsOrder()  [NMATCHES]int{
	// p := [NMATCHES]int{0,1}	// TEST
	a := [NMATCHES]int{}
	p := [NMATCHES]int{0,1,2,3} // OK
	ind := rand.Intn(NMATCHES)
	l := len(p)
	for i := ind; i < l+ind; i++ {
		//fmt.Println(p[i%l], ", i%l: ",i%l)
		a[i-ind] = p[i%l]
	}
	return a
}
//!-chooseStadiumsOrder

//!+fillCentralSystemData
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
//!-fillCentralSystemData

//!+checkInfoOtherScores
func checkInfoOtherScores(csd *csdata, data mcomm, matchid int)  {
	for i := 0; i < NMATCHES; i++ {
		if data.matchdata.otherscores[i].isfollowed && data.matchdata.otherscores[i].invalid {
			// fmt.Println("Habria que pedir la info actualizada")
			mchans[i] <- mcomm{-1, "updatecs", -1, &mdata{}}
			resp := <- cschans[1]
			csd.matches[i].localscore = resp.matchdata.localscore
			csd.matches[i].visitorscore = resp.matchdata.visitorscore
			csd.writeback++
			stats.Lock()
			stats.ALCD++
			stats.Unlock()
			data.matchdata.otherscores[i].olocalscore = resp.matchdata.localscore
			data.matchdata.otherscores[i].ovisitorscore = resp.matchdata.visitorscore
			data.matchdata.otherscores[i].invalid = false
		}
	}
}
//!-checkInfoOtherScores

//!+centralSystem
func centralSystem(start chan struct{}, id int, n *sync.WaitGroup) {
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
				// fmt.Println("SC: ",info)
				storder := chooseStadiumsOrder()
				fmt.Println("El orden elegido de consulta es ", storder)
				for i := 0; i < len(storder); i++ {
					matchid := storder[i]
					mchans[matchid] <- mcomm{id, "getinfo", -1, &mdata{}}
					x := <-cschans[1]
					fmt.Println("SC: E",storder[i],": localscore = ",x.matchdata.localscore,", visitorscore = ",x.matchdata.visitorscore)
					csd.Lock()
					//fmt.Println("SC: E",storder[i],": localscore = ",x.matchdata.localscore,", visitorscore = ",x.matchdata.visitorscore)
					x.matchdata.Lock()
					stats.Lock()
					stats.AMC++
					stats.Unlock()
					csd.matches[matchid].localscore = x.matchdata.localscore
					csd.matches[matchid].visitorscore = x.matchdata.visitorscore
					checkInfoOtherScores(csd, x, matchid)
					csd.matches[matchid].otherscores = x.matchdata.otherscores
					x.matchdata.Unlock()
					csd.invalid[matchid] = false
					csd.Unlock()
				}

				if info.round == TGAME {
					for i := 0; i < NMATCHES; i++ {
						close(mchans[i])
					}
					// fmt.Println("SC : data: ", csd)
					showInfo(&mdata{}, csd, false, id)
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

//!+matchesGenerator
func matchesGenerator()  {
	var n sync.WaitGroup
	defer n.Wait()
	initController()
	start := make(chan struct{}, NTOT)
	myLock(start, NTOT)
	for i := 0; i < NMATCHES; i++ {
		n.Add(1)
		mchans[i] = make(chan mcomm, 0)
		go match(start, i, &n)
	}

	n.Add(1)
	csid := NMATCHES
	cschans[0] = make(chan mcomm, 0)
	cschans[1] = make(chan mcomm, 0)
	go centralSystem(start, csid, &n)
	startCountdown()
	myUnlock(start, NTOT)
}
//!-matchesGenerator

//!+estadistics
func estadistics() {
		if _, err := os.Stat("./Estadisticas.txt"); os.IsNotExist(err) {
			os.Create("Estadisticas.txt")
			b:= []byte("N=0\nILCD=0\nALCD=0\nAMP=0\n")
			err = ioutil.WriteFile("Estadisticas.txt", b, 0644)
	 	 	 if err != nil {
	 	 			 panic(err)
	 	 	 }
		}

		b, err := ioutil.ReadFile("Estadisticas.txt")
 	 	if err != nil {
 			 panic(err)
 	 	}

		dat:=string(b)

		res := strings.Split(dat, "\n")
		N,_:= strconv.Atoi(strings.Split(res[0],"=")[1])
		ILCD,_:= strconv.Atoi(strings.Split(res[1],"=")[1])
		ALCD,_:= strconv.Atoi(strings.Split(res[2],"=")[1])
		AMC,_:= strconv.Atoi(strings.Split(res[3],"=")[1])
		N++
		ILCD=ILCD+stats.ILCD
		ALCD=ALCD+stats.ALCD
		AMC=AMC+stats.AMC
		b= []byte("N="+ strconv.Itoa(N)+"\nILCD="+ strconv.Itoa(ILCD)+"\nALCD="+strconv.Itoa(ALCD)+"\nAMP="+strconv.Itoa(AMC)+"\n")
		err = ioutil.WriteFile("Estadisticas.txt", b, 0644)
		if err != nil {
			 panic(err)
		}
		fmt.Println("----STATS---------")
		fmt.Println("En invalidar otra línea de una caché diferente tardo ",strconv.Itoa(ILCD/2/N), "u.t en promedio")
		fmt.Println("En actualizar esa línea en otra caché diferente tardo ",strconv.Itoa(ALCD/N), "u.t en promedio")
		fmt.Println("en actualizar la memoria principal con el nuevo valor de una línea tardo ",strconv.Itoa(AMC*2/N), "u.t en promedio")
		fmt.Println("DATOS TOMADOS CON ", N, "EJECUCIONES DEL PROGRAMA")
		fmt.Println("-------------------")
}
//!-estadistics

//!+main
func main() {
	rand.Seed(time.Now().UnixNano()) //need to call this function to get diferent pseudo random number every execution
	fmt.Println("Comienza el programa")
	matchesGenerator()
	estadistics()
	fmt.Println("Termina el programa")
}
//!-main

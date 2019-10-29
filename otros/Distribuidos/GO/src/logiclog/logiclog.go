//Jorge Vela Peña
//Grado en Ingeniería Telemática

package logiclog

import(
	"logicclock"
	"fmt"
	"strconv"
	"os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func WriteInFile(a string, f *os.File){
    n2, err := f.Write([]byte(a+ "\n"))
    check(err)
    fmt.Printf("wrote %d bytes\n", n2)
}


type Logiclog struct{
	clock []logicclock.Log
	msg string
	lastUser string
}


func NewLog(clock []logicclock.Log, msg string, lastUser string) *Logiclog{
	return &Logiclog{clock, msg, lastUser}
}


func (Logiclog Logiclog)CreateLog(f *os.File){
	//strLog := "(" + Logiclog.clock[0].user + Logiclog.lastUser + ")" + Logiclog.msg
	stringAGuardar:=""

	stringAGuardar = stringAGuardar+"(["

	z:=0
	for(z<len(Logiclog.clock)){
		stringAGuardar=stringAGuardar + "(" + Logiclog.clock[z].User+ ","+  strconv.Itoa(Logiclog.clock[z].Num) + ")"
		fmt.Println(Logiclog.clock[z].User, Logiclog.clock[z].Num)
		if(z!=len(Logiclog.clock)-1){
			stringAGuardar=stringAGuardar + ","
		}
		z++
	}
	stringAGuardar = stringAGuardar+"]"+Logiclog.lastUser+")"+ Logiclog.msg
	WriteInFile(stringAGuardar, f )
}
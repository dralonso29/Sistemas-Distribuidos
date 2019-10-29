//Jorge Vela Peña
//Grado en Ingeniería Telemática

package logicclock

import(
		//"fmt"
		//"strconv"
	    //"os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}
/*
func PruebaEscribir(a string, f *os.File){
    n2, err := f.Write([]byte(a+ "\n"))
    check(err)
    fmt.Printf("wrote %d bytes\n", n2)
}
*/

/*type logicclock interface{
	Message()
	//loge log.logicclock
}*/

type Log struct{
	Num int
	User string
}

//type clock struct{
//	Array []log
//	lastUser string
//}

type LogicSystem struct{
	Array []Log
	//reloj clock
	//mensaje string
}

func (N *LogicSystem)GetClock() []Log{
	return N.Array
}

func NewClock() *LogicSystem{
	var x []Log
	x=append(x,Log{0,""})
	//var xclock clock = clock{x,""}
	return &LogicSystem{x}
}

func (N *LogicSystem)Message(mLog string, User string, otherUserClock *[]Log){ //f *os.File
	j:=0;
	xotherUserClock := *otherUserClock
	for(j<len(*otherUserClock)){
		if len(*otherUserClock)==1 && xotherUserClock[0].Num == 0  && xotherUserClock[0].User == ""{ 
			break
		}
		i:=0
		for (i<len(N.Array)){
			if N.Array[i].User == xotherUserClock[j].User {
				if  xotherUserClock[j].Num > N.Array[i].Num  {
					N.Array[i].Num=xotherUserClock[j].Num
				}
				i=len(N.Array)+1
				break;
			} 
			i++
		}
		if(i!=len(N.Array)+1){
			if N.Array[0].Num == 0 && N.Array[0].User == "" {
				N.Array[0].Num = xotherUserClock[j].Num
				N.Array[0].User=xotherUserClock[j].User
			}else{
				N.Array = append(N.Array, Log{xotherUserClock[j].Num,xotherUserClock[j].User }) 
			}
		}
		j++
	}


	if N.Array[0].Num == 0 && N.Array[0].User == "" {
		N.Array[0].Num = 1
		N.Array[0].User=User
		//N.lastUser = User
		//N.mensaje=mLog
	}else{
		i:=0
		for (i<len(N.Array)){
			if(N.Array[i].User == User){
				N.Array[i].Num++
				//N.lastUser = User
				//N.mensaje=mLog
				i=len(N.Array)+1
				break;
			}
			i++
		}
		if(i!=len(N.Array)+1){
			N.Array = append(N.Array, Log{1,User}) 
			//N.lastUser = User
			//N.mensaje=mLog	
		}
	}
	/*
	stringAGuardar:=""

	stringAGuardar = stringAGuardar+"["

	z:=0
	for(z<len(N.Array)){
		stringAGuardar=stringAGuardar + "(" + N.Array[z].User+ ","+  strconv.Itoa(N.Array[z].Num) + ")"
		fmt.Println(N.Array[z].User, N.Array[z].Num)
		if(z!=len(N.Array)-1){
			stringAGuardar=stringAGuardar + ","
		}
		z++
	}
	stringAGuardar = stringAGuardar+"]"//+N.lastUser+")"+ N.mensaje


	//fmt.Println(N.lastUser)
	//fmt.Println(N.mensaje ,"\n")

	PruebaEscribir(stringAGuardar, f )*/
}
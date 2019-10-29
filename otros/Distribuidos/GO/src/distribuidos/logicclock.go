//Jorge Vela Peña
//Grado en Ingeniería Telemática

package logicclock

import(
		"fmt"
		"strconv"
	    "os"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func PruebaEscribir(a string, f *os.File){
    n2, err := f.Write([]byte(a+ "\n"))
    check(err)
    fmt.Printf("wrote %d bytes\n", n2)
}

type logicOrder interface{
	Message()
}

type log struct{
	num int
	user string
}

type clock struct{
	array []log
	lastUser string
}

type logicSystem struct{
	reloj clock
	mensaje string
}

func (N *logicSystem)GetClock() []log{
	return N.reloj.array
}

func NewClock() *logicSystem{
	var x []log
	x=append(x,log{0,""})
	var xclock clock = clock{x,""}
	return &logicSystem{xclock, ""}
}

func (N *logicSystem)Message(mLog string, user string, otherUserClock *[]log, f *os.File){
	j:=0;
	xotherUserClock := *otherUserClock
	for(j<len(*otherUserClock)){
		if len(*otherUserClock)==1 && xotherUserClock[0].num == 0  && xotherUserClock[0].user == ""{ 
			break
		}
		i:=0
		for (i<len(N.reloj.array)){
			if N.reloj.array[i].user == xotherUserClock[j].user {
				if  xotherUserClock[j].num > N.reloj.array[i].num  {
					N.reloj.array[i].num=xotherUserClock[j].num
				}
				i=len(N.reloj.array)+1
				break;
			} 
			i++
		}
		if(i!=len(N.reloj.array)+1){
			if N.reloj.array[0].num == 0 && N.reloj.array[0].user == "" {
				N.reloj.array[0].num = xotherUserClock[j].num
				N.reloj.array[0].user=xotherUserClock[j].user
			}else{
				N.reloj.array = append(N.reloj.array, log{xotherUserClock[j].num,xotherUserClock[j].user }) 
			}
		}
		j++
	}


	if N.reloj.array[0].num == 0 && N.reloj.array[0].user == "" {
		N.reloj.array[0].num = 1
		N.reloj.array[0].user=user
		N.reloj.lastUser = user
		N.mensaje=mLog
	}else{
		i:=0
		for (i<len(N.reloj.array)){
			if(N.reloj.array[i].user == user){
				N.reloj.array[i].num++
				N.reloj.lastUser = user
				N.mensaje=mLog
				i=len(N.reloj.array)+1
				break;
			}
			i++
		}
		if(i!=len(N.reloj.array)+1){
			N.reloj.array = append(N.reloj.array, log{1,user}) 
			N.reloj.lastUser = user
			N.mensaje=mLog	
		}
	}

	stringAGuardar:=""

	stringAGuardar = stringAGuardar+"(["

	z:=0
	for(z<len(N.reloj.array)){
		stringAGuardar=stringAGuardar + "(" + N.reloj.array[z].user+ ","+  strconv.Itoa(N.reloj.array[z].num) + ")"
		fmt.Println(N.reloj.array[z].user, N.reloj.array[z].num)
		if(z!=len(N.reloj.array)-1){
			stringAGuardar=stringAGuardar + ","
		}
		z++
	}
	stringAGuardar = stringAGuardar+"]"+N.reloj.lastUser+")"+ N.mensaje


	fmt.Println(N.reloj.lastUser)
	fmt.Println(N.mensaje ,"\n")

	PruebaEscribir(stringAGuardar, f )
}
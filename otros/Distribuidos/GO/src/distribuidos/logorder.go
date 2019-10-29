//Jorge Vela Peña
//Grado en Ingeniería Telemática

package main

import(
	"strconv"
	"fmt"
    "strings"
    "os"

    "container/list" 
)


type logicOrder interface{
	ReadFile()
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type clockStruct interface{
	Message()
}


type logOrder struct {
	s map[string]int
	name string
    msg string
}


var l = list.New()

func compare(lo logOrder, array []string){
    lugarInsertar:=0
    elMismoMensaje :=0 
    e1:=l.Front()
	for e := l.Front(); e != nil; e = e.Next() {
		//fmt.Println(e.Value.(logOrder).s)
        lugarInsertar = 0
        elMismoMensaje=0
		i :=0
		for i<len(array){
            if(e.Value.(logOrder).s[array[i]] ==0 ){
                //fmt.Println("NO EXISTE EL PRIMERO")
            }else if(e.Value.(logOrder).s[array[i]] == lo.s[array[i]]){
                elMismoMensaje++
                //fmt.Println("ES IGUAL")
            }else if(e.Value.(logOrder).s[array[i]] < lo.s[array[i]]){
                if(lugarInsertar ==0 || lugarInsertar==1){
                    lugarInsertar = 1
                }
				//fmt.Println("ES MAYOR")
			}else{
                if(lugarInsertar ==0 || lugarInsertar==2){
                    lugarInsertar = 2
                }
                e1 =e 
				//fmt.Println("ES MENOR")
                break;
			}
			i++
            if(elMismoMensaje==len(array) && e.Value.(logOrder).name == lo.name){ //
                lugarInsertar = 3
                break
            }
		}

        if(lugarInsertar==3){
            break
        }
        //fmt.Println(e.Value.(logOrder).name , lo.name)

        if (lugarInsertar ==0 && len(e.Value.(logOrder).s) > len(lo.s)){
            lugarInsertar=2
            break
        }
        if(lugarInsertar==2){
                break
        }
	}
    if(lugarInsertar == 1 || lugarInsertar==0){
        l.PushBack(lo)
    }else if(lugarInsertar ==2){
        l.InsertBefore(lo, e1)
    }
}

func imprimir(){
    for e := l.Front(); e != nil; e = e.Next() {
        fmt.Println(e.Value)
    }
}


func obtenerUsuariosNum(result []string){
    i:=0
    for i < len(result)-1 {	
		var arrayX []string	
        var lo logOrder
		lo.s = make(map[string]int)
        resultCut := result[i]

        array := strings.Split(resultCut, "]")
        array1 := strings.Split(array[0], "[")

        us :=0
        users := strings.Split(array1[1],"),(")
        for us < len(users) {
            if(us==0){
                x := strings.Split(users[us],"(")
                users[us]=x[1]  
            }
            if(us==len(users)-1){
                x := strings.Split(users[us],")")
                users[us]=x[0]  
            }
            //fmt.Println(users[us])
			userx := strings.Split(users[us],",")
			i, _ := strconv.Atoi(userx[1])
			lo.s[userx[0]]=i
			arrayX= append(arrayX, userx[0])
            us++    
        }

        array2 := strings.Split(array[1], ")")

		lo.name=array2[0]
		lo.msg=array2[1]

		if(l.Front() ==nil){
			fmt.Println("entra")
			l.PushBack(lo)
		}else{
			compare(lo, arrayX)
		}
        //fmt.Println("\n")
        i++
    }
}

func ReadFile(file *os.File) []string {
    strTotal :=""
    i:=0
    for(i<2){
        b1 := make([]byte, 60)
        _, err := file.Read(b1)
        if(err != nil){
            if(err.Error() == "EOF"){
                break
            }
        }
        strTotal = strTotal+ string(b1)
        check(err)
    }
    result := strings.Split(strTotal, "\n")
    return result
}

func main(){
    args := os.Args[1:]
    if(len(args) ==0){
        fmt.Println("No hay ficheros para leer")
    }else{
        i:=0
        for i< len(args){
            f, err := os.Open(args[i])
            check(err)
            result := ReadFile(f)
            obtenerUsuariosNum(result)
            i++
        }
        imprimir()
    }
}
package logicclock


import(
	"fmt"
	"testing"
	"os"
)


func TestOne(t *testing.T){
	fmt.Println("H")
	clockPaco :=NewClock()
	clockJuan :=NewClock()
	clockJose :=NewClock()
	clockMario :=NewClock()

    f1, err := os.Create("/tmp/dat")
    f2, err := os.Create("/tmp/dat2")
    f3, err := os.Create("/tmp/dat3")
    f4, err := os.Create("/tmp/dat4")

    check(err)

	clockJuan.Message("Msg1","Paco", &clockPaco.reloj.array,f2) //Paco a Juan
	clockPaco.Message("Msg2","Juan", &clockJuan.reloj.array,f1) //Juan a Paco
	clockJuan.Message("Msg3","Paco", &clockPaco.reloj.array,f2) //Paco a Juan
	clockPaco.Message("Msg4","Juan", &clockJuan.reloj.array,f1) //Juan a Paco
	clockJuan.Message("Msg5","Paco", &clockPaco.reloj.array,f2) //Paco a Juan
	clockPaco.Message("Msg6","Juan", &clockJuan.reloj.array,f1) //Juan a Paco
	clockJuan.Message("Msg7","Paco", &clockPaco.reloj.array,f2) //Paco a Juan
	clockPaco.Message("Msg8","Juan", &clockJuan.reloj.array,f1) //Juan a Paco

	clockJose.Message("Msg9","Mario",&clockMario.reloj.array,f3)
	clockMario.Message("Msg10","Jose",&clockJose.reloj.array,f4)
	clockJose.Message("Msg11","Mario",&clockMario.reloj.array,f3)
	clockMario.Message("Msg12","Jose",&clockJose.reloj.array,f4)
	clockJose.Message("Msg13","Mario",&clockMario.reloj.array,f3)
	clockMario.Message("Msg14","Jose",&clockJose.reloj.array,f4)

	clockPaco.Message("Msg15","Mario",&clockMario.reloj.array,f1)

	clockJuan.Message("Msg16","Jose",&clockJose.reloj.array,f2);
	clockJuan.Message("Msg17","Paco",&clockPaco.reloj.array,f2);

    /*
	clockJuan.Message("Msg1","Paco", &clockPaco.reloj.array,f) //Paco a Juan
	clockPaco.Message("Msg2","Juan", &clockJuan.reloj.array,f) //Juan a Paco
	clockJuan.Message("Msg3","Paco", &clockPaco.reloj.array,f) //Jose a Paco
	clockJose.Message("Msg4","Juan", &clockJuan.reloj.array,f) //Juan a Jose
	clockJuan.Message("Msg5","Jose", &clockJose.reloj.array,f) //Jose a Juan
	clockJuan.Message("Msg6","Paco", &clockPaco.reloj.array,f) //Jose a Paco
	*/


}


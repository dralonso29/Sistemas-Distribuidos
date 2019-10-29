package main


import (
    "fmt"
    "io/ioutil"
    "bufio"
    "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}


type countWords struct {
    W string
    C int
}

   
var arrayWords = []countWords{}

   
 func ExampleScanner_words(input string) {
			wordRepeat := false
    		scanner := bufio.NewScanner(strings.NewReader(input))
    		scanner.Split(bufio.ScanWords)   
    		
    		for scanner.Scan() {
				for i := range(arrayWords) {
					if(scanner.Text() == arrayWords[i].W){
						arrayWords[i].C = arrayWords[i].C + 1
						wordRepeat = true
						break;
					}
				}
				if(wordRepeat == false){
					arrayWords = append(arrayWords, countWords{scanner.Text(),1})
				}
				wordRepeat = false
    		}
}
    
func main(){    

	dat, err := ioutil.ReadFile("/tmp/dat")
	check(err)
	ExampleScanner_words(string(dat))



    for i := range(arrayWords) {
		place := arrayWords[i]
		fmt.Println(place)
    }
}

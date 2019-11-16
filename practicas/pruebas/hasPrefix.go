package main

import (
	"fmt"
	"strings"
)

//!+main
func main() {
	token := "!private pepe"
	fmt.Println(strings.HasPrefix(token, "!private"))
	fmt.Println(strings.HasPrefix(token, "!priv"))
	splited := strings.Fields(token)
	fmt.Println((len(splited)==2)&&(splited[0]=="!private"))
	splited = strings.Fields("")
	fmt.Println(len(splited))
}

//!-main

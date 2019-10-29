//Jorge Vela Peña
//Grado en Ingeniería Telemática

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

var m map[string]int
var mk []string

func readWord(input string) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		i := m[scanner.Text()]
		m[scanner.Text()] = i + 1
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	m = make(map[string]int)

	for i := range argsWithoutProg {
		dat, err := ioutil.ReadFile(os.Args[i+1])
		check(err)
		readWord(string(dat))
	}

	mk := make([]string, len(m))
	i := 0
	for k, _ := range m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)
	for j := range mk {
		fmt.Println(mk[j], m[mk[j]])
	}
}

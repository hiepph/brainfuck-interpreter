package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	in := "++>+++++"
	s := bufio.NewScanner(strings.NewReader(in))
	s.Split(bufio.ScanRunes)

	for s.Scan() {
		fmt.Println(s.Text())
	}
}

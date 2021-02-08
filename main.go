package main

import (
	"fmt"
	"os"
)

var cmdMap map[string]func(...string) error

func init() {
	cmdMap = map[string]func(...string) error{}
	cmdMap["launch"] = launch
	cmdMap["init"] = makeSeed
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("too few")
		return
	}
	fmt.Println(cmdMap[os.Args[1]](os.Args[2:]...))
}

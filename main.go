package main

import (
	"fmt"
	"os"
)

func main() {
	var cmdMap map[string]func(...string) error
	cmdMap = map[string]func(...string) error{}
	cmdMap["launch"] = launch
	cmdMap["init"] = makeSeed
	cmdMap["list"] = list
	cmdMap["remove"] = remove
	cmdMap["install"] = install
	cmdMap["insert"] = insert

	if len(os.Args) < 3 {
		fmt.Println("too few")
		return
	}
	v, ok := cmdMap[os.Args[1]]
	if !ok {
		fmt.Println("command not found")
		return
	}
	fmt.Println(v(os.Args[2:]...))
}

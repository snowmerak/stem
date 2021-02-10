package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/snowmerak/mocha/err"
)

//remove ... remove names matched regexp
func remove(args ...string) error {
	fmt.Println("get list")
	data, e := getList(args...)
	if e != nil {
		return err.Wrap(e, "remove-GetList")
	}
	if len(data) < 1 {
		return err.New("remove", "data from getList is zero")
	}
	fmt.Println("remove vm instances")
	{
		cmd := exec.Command("multipass", "delete")
		for _, v := range data {
			cmd.Args = append(cmd.Args, v.Name)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if e := cmd.Run(); e != nil {
			return err.Wrap(e, "remove-MultipassDelete")
		}
	}
	fmt.Println("purge vm instances")
	{
		cmd := exec.Command("multipass", "purge")
		for _, v := range data {
			cmd.Args = append(cmd.Args, v.Name)
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if e := cmd.Run(); e != nil {
			return err.Wrap(e, "remove-MultipassPurge")
		}
	}
	return nil
}

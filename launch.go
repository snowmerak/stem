package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/snowmerak/mocha/err"
)

// template ...
// define name, cpu number, ram and disk capacity(mb), using
// what kind of script language.
// launch <dir>
type template struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	CPU   int    `json:"cpu"`
	RAM   int    `json:"ram"`
	Disk  int    `json:"disk"`
}

const templateFileName = "template.json"

//makeSeed ... make seed file
func makeSeed(args ...string) error {
	dir := args[0]
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		if e := os.Mkdir(dir, os.ModePerm); e != nil {
			return err.Wrap(e, "maketemplate-Mkdir")
		}
	}

	filename := filepath.Join(dir, templateFileName)

	data := template{}
	data.Name = "new_instant"
	data.CPU = 1
	data.Disk = 1024
	data.RAM = 512

	f, e := os.Create(filename)
	if e != nil {
		return err.Wrap(e, "makeSeed-CreateFile")
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	if e := encoder.Encode(data); e != nil {
		return err.Wrap(e, "makeSeed-EncodeJSON")
	}

	return nil
}

func launch(args ...string) error {
	data := template{}
	dir := args[0]
	filename := filepath.Join(dir, templateFileName)

	// if not exist template.json, create new seed.json with standard data
	fmt.Println("check template file")
	if _, e := os.Stat(filename); os.IsNotExist(e) {
		fmt.Println("create template file")
		data.Name = "new_instant"
		data.CPU = 1
		data.Disk = 1024
		data.RAM = 512

		f, e := os.Create(filename)
		if e != nil {
			return err.Wrap(e, "launch-CreateFile")
		}

		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")

		if e := encoder.Encode(data); e != nil {
			return err.Wrap(e, "launch-EncodeJSON")
		}

		fmt.Println("template file is created")
		return nil
	}

	// open and read from template.json
	fmt.Println("read template file")
	{
		f, e := os.Open(filename)
		if e != nil {
			return err.Wrap(e, "launch-OpenFile")
		}

		decoder := json.NewDecoder(f)
		if e := decoder.Decode(&data); e != nil {
			return err.Wrap(e, "launch-DecodeJSON")
		}
	}

	vmName := data.Name

	// create VM instance
	fmt.Println("launch VM instance")
	{
		n := 1
		buf := bytes.NewBuffer([]byte{})
		for {
			cmd := exec.Command("multipass", "launch", "-c", fmt.Sprint(data.CPU), "-m", fmt.Sprintf("%dmb", data.RAM), "-d", fmt.Sprintf("%dmb", data.Disk), "-n", vmName, data.Image)
			cmd.Stderr = buf
			fmt.Printf("try launching %v\n", vmName)
			e := cmd.Run()
			buffer := buf.String()
			if e == nil {
				break
			}
			if e != nil && !strings.Contains(buffer, "already exists") {
				return err.Wrap(errors.New(buffer), "launch-CmdLaunch")
			}
			fmt.Printf("%v is already exist\n", vmName)
			vmName = fmt.Sprintf("%s-%d", data.Name, n)
			cmd.Args[len(cmd.Args)-1] = vmName
			n++
			buf.Reset()
		}
	}

	// apt update && upgrade
	fmt.Println("apt update && upgrade")
	{
		cmd := exec.Command("multipass", "exec", vmName, "--", "sudo", "apt", "update", "-y")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if e := cmd.Run(); e != nil {
			return err.Wrap(e, "launch-AptUpdate")
		}
		cmd = exec.Command("multipass", "exec", vmName, "--", "sudo", "apt", "upgrade", "-y")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if e := cmd.Run(); e != nil {
			return err.Wrap(e, "launch-AptUpgrade")
		}
	}

	fmt.Printf("%v is launched\n", vmName)

	return nil
}

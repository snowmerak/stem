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

// initData ...
// define name, cpu number, ram and disk capacity(mb), using
// what kind of script language.
type initData struct {
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	CPU         int      `json:"cpu"`
	RAM         int      `json:"ram"`
	Disk        int      `json:"disk"`
	UsingPython bool     `json:"using_python"`
	UsingRuby   bool     `json:"using_ruby"`
	UsingJulia  bool     `json:"using_julia"`
	Scripts     []string `json:"scripts"`
}

//makeSeed ... make seed file
func makeSeed(args ...string) error {
	dir := args[0]
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		if e := os.Mkdir(dir, os.ModePerm); e != nil {
			return err.Wrap(e, "makeSeed-Mkdir")
		}
	}

	filename := filepath.Join(dir, "seed.json")

	data := initData{}
	data.Name = "new_instant"
	data.CPU = 1
	data.Disk = 1024
	data.RAM = 512
	data.UsingJulia = false
	data.UsingPython = false
	data.UsingRuby = false
	data.Scripts = []string{}

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
	data := initData{}
	dir := args[0]
	filename := filepath.Join(dir, "seed.json")

	// if not exist seed.json, create new seed.json with standard data
	fmt.Println("check seed file")
	if _, e := os.Stat(filename); os.IsNotExist(e) {
		fmt.Println("create seed file")
		data.Name = "new_instant"
		data.CPU = 1
		data.Disk = 1024
		data.RAM = 512
		data.UsingJulia = false
		data.UsingPython = false
		data.UsingRuby = false
		data.Scripts = []string{}

		f, e := os.Create(filename)
		if e != nil {
			return err.Wrap(e, "launch-CreateFile")
		}

		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")

		if e := encoder.Encode(data); e != nil {
			return err.Wrap(e, "launch-EncodeJSON")
		}
	}

	// open and read from seed.json
	fmt.Println("read seed file")
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

	// install script language
	fmt.Println("install script languages")
	{
		if data.UsingJulia {
			cmd := exec.Command("multipass", "exec", vmName, "--", "sudo", "snap", "install", "julia", "--classic")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if e := cmd.Run(); e != nil {
				return err.Wrap(e, "launch-InstallJulia")
			}
		}
		if data.UsingRuby {
			cmd := exec.Command("multipass", "exec", vmName, "--", "sudo", "snap", "install", "ruby", "--classic")
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if e := cmd.Run(); e != nil {
				return err.Wrap(e, "launch-InstallRuby")
			}
		}
	}

	// transfer script file and run
	fmt.Println("install scripts")
	{
		for _, v := range data.Scripts {
			scriptPath := filepath.Join(dir, v)
			scriptFile := filepath.Base(scriptPath)
			{
				cmd := exec.Command("multipass", "transfer", scriptPath, fmt.Sprintf("%s:%s", vmName, scriptFile))
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if e := cmd.Run(); e != nil {
					return err.Wrap(e, "launch-TransferScript")
				}
			}
			{
				cmd := exec.Command("multipass", "exec", vmName, "--")
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				switch filepath.Ext(scriptPath) {
				case ".rb":
					cmd.Args = append(cmd.Args, "sudo", "snap", "run", "ruby", scriptFile)
				case ".jl":
					cmd.Args = append(cmd.Args, "sudo", "snap", "run", "julia", scriptFile)
				case ".py":
					cmd.Args = append(cmd.Args, "sudo", "python3", scriptFile)
				}
				if e := cmd.Run(); e != nil {
					return err.Wrap(e, "launch-RunScript")
				}
			}
		}
	}

	return nil
}

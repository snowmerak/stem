package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/snowmerak/mocha/err"
)

type seed struct {
	UsingPython bool     `json:"using_python"`
	UsingRuby   bool     `json:"using_ruby"`
	UsingJulia  bool     `json:"using_julia"`
	Scripts     []string `json:"scripts"`
}

const seedFileName = "seed.json"

//install stem install <vmName> <dir>
func install(args ...string) error {
	if len(args) < 2 {
		return err.New("install", "args is too feww")
	}
	data := seed{}
	dir := args[1]
	filename := filepath.Join(dir, seedFileName)

	// if not exist seed.json, create new seed.json with standard data
	fmt.Println("check seed file")
	if _, e := os.Stat(filename); os.IsNotExist(e) {
		fmt.Println("create seed file")
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

		fmt.Println("template file is created")
		return nil
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

	vmName := args[0]

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
			fmt.Printf("copy %v\n", scriptFile)
			{
				cmd := exec.Command("multipass", "transfer", scriptPath, fmt.Sprintf("%s:%s", vmName, scriptFile))
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if e := cmd.Run(); e != nil {
					return err.Wrap(e, "launch-TransferScript")
				}
			}
			fmt.Printf("execute %v\n", scriptFile)
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
			fmt.Printf("remove %v\n", scriptFile)
			{
				cmd := exec.Command("multipass", "exec", vmName, "sudo", "rm", scriptFile)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if e := cmd.Run(); e != nil {
					return err.Wrap(e, "launch-RemoveScript")
				}
			}
		}
	}

	fmt.Printf("launched %v\n", vmName)

	return nil
}

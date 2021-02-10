package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/snowmerak/mocha/err"
)

const dataFolderName = "fertilizer"

//insert ... stem insert <dir> <vmName>
func insert(args ...string) error {
	if len(args) < 2 {
		return err.New("insert", "args is too few")
	}
	dir := args[1]
	vmName := args[0]
	rootDirName := filepath.Base(dataFolderName)

	{
		f, e := os.Stat(dir)
		if os.IsNotExist(e) {
			os.Mkdir(filepath.Join(dir, dataFolderName), 0777)
			fmt.Println(fmt.Sprintf("%v is maed\n", dataFolderName))
			return nil
		}
		if e != nil {
			return err.Wrap(e, "insert=StateDir")
		}
		if !f.IsDir() {
			return err.Wrap(e, "insert=IsNotDirectory")
		}
	}

	fmt.Println("mount folder")
	{
		cmd := exec.Command("multipass", "mount", dir, fmt.Sprintf("%v:/media/%v", vmName, rootDirName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if e := cmd.Run(); e != nil {
			return err.Wrap(e, "insert-MultipassMount")
		}
	}

	fmt.Println("copy folder")
	{
		cmd := exec.Command("multipass", "exec", vmName, "--", "sudo", "cp", "-r", fmt.Sprintf("/media/%v", rootDirName), filepath.Base(args[1]))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if e := cmd.Run(); e != nil {
			return err.Wrap(e, "insert-CopyFolder")
		}
	}

	fmt.Println("unmount folder")
	{
		cmd := exec.Command("multipass", "unmount", fmt.Sprintf("%v:/media/%v", vmName, rootDirName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if e := cmd.Run(); e != nil {
			return err.Wrap(e, "insert-MultipassUnmount")
		}
	}

	fmt.Println("end inserting")

	return nil
}

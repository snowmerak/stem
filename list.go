package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/snowmerak/mocha/err"
)

//Info ... short information of instance
type Info struct {
	Ipv4    []string `json:"ipv4"`
	Name    string   `json:"name"`
	Release string   `json:"release"`
	State   string   `json:"state"`
}

//List ... I didn't want this
type List struct {
	List []Info `json:"list"`
}

//getList ... get list of regexp
func getList(args ...string) ([]Info, error) {
	data := List{}
	result := []Info{}
	{
		cmd := exec.Command("multipass", "list", "--format", "json")
		rs, e := cmd.Output()
		if e != nil {
			return nil, err.Wrap(e, "getList-GetList")
		}

		decoder := json.NewDecoder(bytes.NewReader(rs))
		if e := decoder.Decode(&data); e != nil {
			return nil, err.Wrap(e, "getList-Decode")
		}
	}

	{
		regex, e := regexp.Compile(args[0])
		if e != nil {
			return nil, err.Wrap(e, "getList-CompileRegex")
		}
		for _, v := range data.List {
			if regex.MatchString(v.Name) {
				result = append(result, v)
			}
		}
	}

	return result, nil
}

//list ... get list of regexp and print
func list(args ...string) error {
	data, e := getList(args...)
	if e != nil {
		return err.Wrap(e, "list-GetList")
	}
	for _, v := range data {
		fmt.Printf("%v: %v [%v]\n", v.Name, v.State, strings.Join(v.Ipv4, ", "))
	}
	return nil
}

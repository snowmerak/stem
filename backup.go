package main

import (
	"sync"

	"github.com/snowmerak/mocha/err"
)

//clone ... get tar
func clone(args ...string) error {
	ls, e := getList(args...)
	if e != nil {
		return err.Wrap(e, "clone-GetList")
	}

	wg := sync.WaitGroup{}
	ch := make(chan error)
	for _, v := range ls {
		wg.Add(1)
		go func(name string) {

		}(v.Name)
	}

	return nil
}

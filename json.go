package main

import (
	"github.com/bfun/cjsonsource"
	"sync"
)

func ParseAllJsonC(wg *sync.WaitGroup) {
	defer wg.Done()
	JSONMAP = cjsonsource.ParseJsonSourceJson()
}

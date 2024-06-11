package main

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

func get_svcname_by_procode() (services map[string]map[string]string) {
	services = make(map[string]map[string]string)
	fileName := path.Join(getRootDir(), "etc/enum/get_svcname_by_procode.txt")
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "enumData1") || strings.HasPrefix(line, "*") {
			continue
		}
		kv := strings.Split(line, "\t")
		if len(kv) != 2 {
			panic(line)
		}
		svc := kv[1]
		if strings.Contains(kv[0], ".") {
			ds := strings.Split(kv[0], ".")
			dta, ok := services[ds[0]]
			if !ok {
				dta = make(map[string]string)
				services[ds[0]] = dta
			}
			dta[ds[1]] = svc
		} else {
			dtaname := "TXML_SVR"
			dta, ok := services[dtaname]
			if !ok {
				dta = make(map[string]string)
				services[dtaname] = dta
			}
			dta[kv[0]] = svc
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}

func ParseGetSvcnameByProcode(wg *sync.WaitGroup) {
	defer wg.Done()
	GetSvcNameMAP = get_svcname_by_procode()
	log.Print("get_svcname_by_procode.txt parse success")
}

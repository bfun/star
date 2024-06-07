package main

import (
	"fmt"
	"github.com/bfun/cjsonsource"
	"log"
	"os"
	"sync"
)

var ESADMIN ESAdmin
var PROJECT Project
var DTALIST []string
var ALALIST []string
var DTAMAP map[string]DataTransferAdapter
var SVCMAP map[string]map[string]Service
var RUTMAP map[string]map[string]Entrance
var FMTMAP map[string]Format
var FLOWMAP map[string]Flow
var LOGICMAP map[string]map[string]Logic
var JSONMAP map[string]map[string]cjsonsource.SvcFunc

var PORT string = "8080"

// Init
// if named init, can't run unit test
func Init() {
	fmt.Printf("len(os.Args): %d\n", len(os.Args))
	if len(os.Args) > 1 {
		PORT = os.Args[1]
	}
	ESADMIN = ParseESAdminFile()
	log.Print("ESAdmin.xml parse success")
	PROJECT = ParseProjectFile()
	for _, v := range PROJECT.PubDtas {
		DTALIST = append(DTALIST, v.Name)
	}
	for _, v := range PROJECT.Apps {
		for _, sub := range v.SubApps {
			for _, dta := range sub.Dtas {
				DTALIST = append(DTALIST, dta.Name)
			}
			ALALIST = append(ALALIST, sub.Name)
		}
	}
	log.Print("Project.xml parse success")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go ParseAllDtaParmXml(wg)
	wg.Add(1)
	go ParseAllServiceXml(wg)
	wg.Add(1)
	go ParseAllRouteXml(wg)
	wg.Add(1)
	go ParseAllFormatXml(wg)
	wg.Add(1)
	go ParseAllFlowXml(wg)
	wg.Add(1)
	go ParseAllSvcLogicXml(wg)
	wg.Add(1)
	go ParseAllJsonC(wg)
	wg.Wait()
	linkServicesToDtas(SVCMAP, DTAMAP)
}

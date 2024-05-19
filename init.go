package star

import (
	"log"
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

func init() {
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
	wg.Wait()
	linkServicesToDtas(SVCMAP, DTAMAP)
}

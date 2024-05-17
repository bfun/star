package star

import "log"

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
	DTAMAP = ParseAllDtaParmXml()
	log.Print("DtaParm.xml parse success")
	SVCMAP = ParseAllServiceXml()
	log.Print("Service.xml parse success")
	RUTMAP = ParseAllRouteXml()
	log.Print("Route.xml parse success")
	FMTMAP = ParseAllFormatXml()
	log.Print("Format.xml parse success")
	linkServicesToDtas(SVCMAP, DTAMAP)
}

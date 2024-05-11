package star

var ESADMIN ESAdmin
var PROJECT Project
var DTALIST []string
var ALALIST []string
var DTAMAP map[string]DataTransferAdapter
var SVCMAP map[string]map[string]Service

func init() {
	ESADMIN = ParseESAdminFile()
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
	DTAMAP = ParseAllDtaParmXml()
	SVCMAP = ParseAllServiceXml()
}

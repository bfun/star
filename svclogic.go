package star

import (
	"fmt"
	"log"
	"path"
	"strings"
	"sync"
)

type LogicTab struct {
	Logics []Logic `xml:"Logic"`
}
type Logic struct {
	Name      string     `xml:"Name,attr"`
	Desc      string     `xml:"Desc,attr"`
	FlowModel string     `xml:"FlowModel,attr"`
	CustProcs []CustProc `xml:"CustProcTab>CustProc,omitempty"`
}

type CustProc struct {
	SeqNo     int    `xml:"SeqNo,attr"`
	BeginProc string `xml:"BeginProc>Proc>ExprProc,omitempty"`
	EndProc   string `xml:"EndProc>Proc>ExprProc,omitempty"`
	ErrProc   string `xml:"ErrProc>Proc>ExprProc,omitempty"`
}

func trimLogicCDATA(v *LogicTab) {
	for _, logic := range v.Logics {
		for _, custProc := range logic.CustProcs {
			custProc.BeginProc = strings.TrimSpace(custProc.BeginProc)
			custProc.EndProc = strings.TrimSpace(custProc.EndProc)
			custProc.ErrProc = strings.TrimSpace(custProc.ErrProc)

		}
	}
}

func logicArrayToMap(logics []Logic) map[string]Logic {
	m := map[string]Logic{}
	for _, logic := range logics {
		m[logic.Name] = logic
	}
	return m
}

func parseOneSvcLogicXml(fileName string) map[string]Logic {
	fullPath := path.Join(getRootDir(), fileName)
	decoder := getStarFileDecoder(fullPath)
	var v LogicTab
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	trimLogicCDATA(&v)
	return logicArrayToMap(v.Logics)
}

func ParseAllSvcLogicXml(wg *sync.WaitGroup) {
	defer wg.Done()
	m := make(map[string]map[string]Logic)
	files := getSvcLogicFiles()
	for ala, file := range files {
		m[ala] = parseOneSvcLogicXml(file)
	}
	LOGICMAP = m
	fmt.Printf("%#v\n", m)
	log.Print("SvcLogic.xml parse success")
}

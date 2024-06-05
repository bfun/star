package star

import (
	"fmt"
	"log"
	"path"
	"sync"
)

type FlowTab struct {
	Flows []Flow `xml:"Flow"`
}

type Flow struct {
	Name      string     `xml:"Name,attr"`
	FlowDesc  string     `xml:"FlowDesc,attr"`
	FlowSteps []FlowStep `xml:"FlowStepTab>FlowStep,omitempty"`
}
type FlowStep struct {
	SeqNo      int        `xml:"SeqNo,attr"`
	StepType   string     `xml:"StepType,attr"`
	StepDesc   string     `xml:"StepDesc,attr"`
	CallType   string     `xml:"CallType,attr"`
	Expression string     `xml:"Expression,chardata"`
	Condition  string     `xml:"Condition,cdata"`
	NextSteps  []NextStep `xml:"NextStepTab>NextStep"`
}
type NextStep struct {
	Value string `xml:"Value,attr"`
	SeqNo int    `xml:"SeqNo,attr"`
}

func trimFlowCDATA(v *FlowTab) {
	for _, f := range v.Flows {
		for _, s := range f.FlowSteps {
			fmt.Printf("%v\t%v\t[%v][%v]\n", f.Name, s.SeqNo, s.Expression, s.Condition)
		}
	}
}

func parseOneFlowXml(fileName string) []Flow {
	fullPath := path.Join(getRootDir(), fileName)
	decoder := getStarFileDecoder(fullPath)
	var v FlowTab
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	trimFlowCDATA(&v)
	return v.Flows
}

func ParseAllFlowXml(wg *sync.WaitGroup) {
	defer wg.Done()
	m := make(map[string]Flow)
	files := getFlowFiles()
	for _, file := range files {
		flows := parseOneFlowXml(file)
		for _, flow := range flows {
			m[flow.Name] = flow
		}
	}
	FLWMAP = m
	log.Print("Flow.xml parse success")
}

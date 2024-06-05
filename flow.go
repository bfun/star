package star

import (
	"fmt"
	"html/template"
	"log"
	"path"
	"strings"
	"sync"
)

type FlowTab struct {
	Flows []Flow `xml:"Flow"`
}

type Flow struct {
	Name      string     `xml:"Name,attr"`
	FlowDesc  string     `xml:"FlowDesc,attr"`
	FlowSteps []FlowStep `xml:"FlowStepTab>FlowStep,omitempty"`
	Chart     string
}
type FlowStep struct {
	SeqNo      int        `xml:"SeqNo,attr"`
	StepType   string     `xml:"StepType,attr"`
	StepDesc   string     `xml:"StepDesc,attr"`
	CallType   string     `xml:"CallType,attr"`
	Expression string     `xml:"Expression"`
	Condition  string     `xml:"Condition"`
	NextSteps  []NextStep `xml:"NextStepTab>NextStep"`
}
type NextStep struct {
	Value string `xml:"Value,attr"`
	SeqNo int    `xml:"SeqNo,attr"`
}

func flowChart(flow Flow) string {
	md := "flowchart TD"
	for _, step := range flow.FlowSteps {
		for _, next := range step.NextSteps {
			if next.Value == "TRUE" {
				md += fmt.Sprintf("\t%v{%v} --->|Y| %v\n", step.SeqNo, template.HTMLEscapeString(step.Condition), next.SeqNo)
			} else if next.Value == "FALSE" {
				md += fmt.Sprintf("\t%v{%v} --->|N| %v\n", step.SeqNo, template.HTMLEscapeString(step.Condition), next.SeqNo)
			} else if step.StepType == "E" {
				md += fmt.Sprintf("\t%v --->|[%v]| %v\n", step.SeqNo, template.HTMLEscapeString(step.Expression), next.SeqNo)
			} else if step.StepType == "D" {
				md += fmt.Sprintf("\t%v --->|[call %v]| %v\n", step.SeqNo, step.CallType, next.SeqNo)
			} else {
				md += fmt.Sprintf("\t%v ---> %v\n", step.SeqNo, next.SeqNo)
			}
		}
	}
	return md
}

func addFlowChart(v *FlowTab) {
	for _, flow := range v.Flows {
		flow.Chart = flowChart(flow)
	}
}

func trimFlowCDATA(v *FlowTab) {
	for _, f := range v.Flows {
		for _, s := range f.FlowSteps {
			// fmt.Printf("%v\t%v\t[%v][%v]\n", f.Name, s.SeqNo, s.Expression, s.Condition)
			s.Expression = strings.TrimSpace(s.Expression)
			s.Condition = strings.TrimSpace(s.Condition)
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
	FLOWMAP = m
	log.Print("Flow.xml parse success")
}

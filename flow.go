package main

import (
	"fmt"
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
	md := "flowchart TD\n"
	for _, step := range flow.FlowSteps {
		for _, next := range step.NextSteps {
			if step.StepType == "W" {
				if next.Value == "TRUE" {
					// md += fmt.Sprintf("\t%v --->|Y| %v\n", step.SeqNo, next.SeqNo)
					md += fmt.Sprintf("\t%v --->|\"%v\"| %v\n", step.SeqNo, PrepareMarkdownText(step.Condition), next.SeqNo)
				} else if next.Value == "FALSE" {
					md += fmt.Sprintf("\t%v --->|N| %v\n", step.SeqNo, next.SeqNo)
				}
			} else if step.StepType == "E" {
				md += fmt.Sprintf("\t%v[\"%v\"] ---> %v\n", step.SeqNo, PrepareMarkdownText(step.Expression), next.SeqNo)
			} else if step.StepType == "D" {
				var text string
				if step.CallType == "sync" {
					text = "同步调用服务方"
				} else if step.CallType == "async" {
					text = "异步调用服务方"
				} else {
					text = "调用服务方"
				}
				md += fmt.Sprintf("\t%v[%v] ---> %v\n", step.SeqNo, text, next.SeqNo)
			} else if step.StepType == "N" {
				md += fmt.Sprintf("\t%v[END]\n", step.SeqNo)
			} else if next.SeqNo != -1 {
				md += fmt.Sprintf("\t%v ---> %v\n", step.SeqNo, next.SeqNo)
			}
		}
		if step.StepType == "W" {
			// md += fmt.Sprintf("\tnote right of %v \"%v\"\n", step.SeqNo, PrepareMarkdownText(step.Condition))
			md += fmt.Sprintf("\t%v{IF}\n", step.SeqNo)
		}
	}
	return md
}
func flowDtasInfo(flow Flow) string {
	var s string
	for _, step := range flow.FlowSteps {
		if step.StepType == "E" {
			dtaName, svcName := parseDtaSvcFromExpression(step.Expression)
			dta, ok := DTAMAP[dtaName]
			if ok {
				s += fmt.Sprintf("%s\t%s\n", dtaName, dta.DTADesc)
				for _, node := range dta.Nodes {
					s += fmt.Sprintf("\t%s\t%s\t%s:%s\n", node.Name, node.Desc, node.IP, node.Port)
				}
			}
			_ = svcName
		}
	}
	return s
}
func addFlowChart(v *FlowTab) {
	for _, flow := range v.Flows {
		flow.Chart = flowChart(flow)
	}
}

func trimFlowCDATA(v *FlowTab) {
	for _, f := range v.Flows {
		for j, s := range f.FlowSteps {
			// fmt.Printf("%v\t%v\t[%v][%v]\n", f.Name, s.SeqNo, s.Expression, s.Condition)
			s.Expression = strings.TrimSpace(s.Expression)
			s.Condition = strings.TrimSpace(s.Condition)
			f.FlowSteps[j] = s
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

package main

import (
	"encoding/xml"
	"log"
	"path"
	"regexp"
	"strings"
	"sync"
)

type DataTransferAdapter struct {
	XMLName          xml.Name      `xml:"DataTransferAdapter"`
	Name             string        `xml:"Name,attr"`
	DTADesc          string        `xml:"DTADesc,attr"`
	EvtIprtcfmtBegin string        `xml:"EvtIprtcfmtBegin"`
	EvtIprtcfmtEnd   string        `xml:"EvtIprtcfmtEnd"`
	EvtOprtcfmtBegin string        `xml:"EvtOprtcfmtBegin"`
	EvtIfmtEnd       string        `xml:"EvtIfmtEnd"`
	EvtOfmtBegin     string        `xml:"EvtOfmtBegin"`
	EvtAcallBegin    string        `xml:"EvtAcallBegin"`
	Nodes            []DtaParmNode `xml:"NodeTab>Node"`
	ConvertPin       bool
	Services         map[string]Service
	NESB_SDTA_NAME   string
	NESB_DDTA_NAME   string
}

type DtaParmNode struct {
	Name string `xml:"Name,attr"`
	Desc string `xml:"Desc,attr"`
	IP   string
	Port string
}

func trimDtaParmCDATA(d *DataTransferAdapter) {
	d.EvtIprtcfmtBegin = strings.TrimSpace(d.EvtIprtcfmtBegin)
	d.EvtIprtcfmtEnd = strings.TrimSpace(d.EvtIprtcfmtEnd)
	d.EvtIfmtEnd = strings.TrimSpace(d.EvtIfmtEnd)
	d.EvtOfmtBegin = strings.TrimSpace(d.EvtOfmtBegin)
	d.EvtOprtcfmtBegin = strings.TrimSpace(d.EvtOprtcfmtBegin)
}

func judgeConvertPin(dtas map[string]DataTransferAdapter) {
	target := "nesbConvertPin"
	for k, v := range dtas {
		if strings.Contains(v.EvtIfmtEnd, target) {
			v.ConvertPin = true
			dtas[k] = v
		}
	}
}

func parseNESB_SDTA_NAME(dtas map[string]DataTransferAdapter) {
	target := "$NESB_SDTA_NAME"
	re := regexp.MustCompile(`\$NESB_SDTA_NAME="(.*?)"`)
	for k, v := range dtas {
		if strings.Contains(v.EvtIprtcfmtBegin, target) {
			s := re.FindStringSubmatch(v.EvtIprtcfmtBegin)
			if len(s) == 2 {
				v.NESB_SDTA_NAME = s[1]
			}
		}
		if strings.Contains(v.EvtIprtcfmtEnd, target) {
			s := re.FindStringSubmatch(v.EvtIprtcfmtEnd)
			if len(s) == 2 {
				v.NESB_SDTA_NAME = s[1]
			}
		}
		if v.NESB_SDTA_NAME != "" {
			dtas[k] = v
		}
	}
}
func parseNESB_DDTA_NAME(dtas map[string]DataTransferAdapter) {
	target := "$NESB_DDTA_NAME"
	re := regexp.MustCompile(`\$NESB_DDTA_NAME="(.*?)"`)
	for k, v := range dtas {
		if strings.Contains(v.EvtOfmtBegin, target) {
			s := re.FindStringSubmatch(v.EvtOfmtBegin)
			if len(s) == 2 {
				v.NESB_DDTA_NAME = s[1]
			}
		}
		if strings.Contains(v.EvtAcallBegin, target) {
			s := re.FindStringSubmatch(v.EvtAcallBegin)
			if len(s) == 2 {
				v.NESB_DDTA_NAME = s[1]
			}
		}
		if v.NESB_DDTA_NAME != "" {
			dtas[k] = v
		}
	}
}
func parseNodeInfo(dtas map[string]DataTransferAdapter) {
	for dtaName, dtap := range dtas {
		for _, edta := range ESADMIN.DtaParms {
			if edta.DtaName == dtaName {
				for _, edtaNode := range edta.Nodes {
					for k, dtapNode := range dtap.Nodes {
						if edtaNode.NodeName == dtapNode.Name {
							dtapNode.IP = edtaNode.Port.NodeIP
							dtapNode.Port = edtaNode.Port.NodePort
							dtap.Nodes[k] = dtapNode
						}
					}
				}
				break
			}
		}
	}
}
func parseOneDtaParmXml(fileName string) DataTransferAdapter {
	fullPath := path.Join(getRootDir(), fileName)
	var v DataTransferAdapter
	decoder := getStarFileDecoder(fullPath)
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	trimDtaParmCDATA(&v)
	return v
}

func ParseAllDtaParmXml(wg *sync.WaitGroup) {
	defer wg.Done()
	m := make(map[string]DataTransferAdapter)
	files := getDtaParmFiles()
	for dta, file := range files {
		m[dta] = parseOneDtaParmXml(file)
	}
	judgeConvertPin(m)
	parseNESB_SDTA_NAME(m)
	parseNESB_DDTA_NAME(m)
	parseNodeInfo(m)
	DTAMAP = m
	log.Print("DtaParm.xml parse success")
}

func linkServicesToDtas(svcs map[string]map[string]Service, dtas map[string]DataTransferAdapter) {
	for k, v := range svcs {
		d, ok := dtas[k]
		if !ok {
			panic(k + " not found in dtas")
		}
		d.Services = v
		dtas[k] = d
	}
}

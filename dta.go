package star

import (
	"encoding/xml"
	"path"
	"regexp"
	"strings"
)

type DataTransferAdapter struct {
	XMLName          xml.Name `xml:"DataTransferAdapter"`
	Name             string   `xml:"Name,attr"`
	EvtIprtcfmtBegin string   `xml:"EvtIprtcfmtBegin"`
	EvtIprtcfmtEnd   string   `xml:"EvtIprtcfmtEnd"`
	EvtIfmtEnd       string   `xml:"EvtIfmtEnd"`
	EvtOfmtBegin     string   `xml:"EvtOfmtBegin"`
	EvtOprtcfmtBegin string   `xml:"EvtOprtcfmtBegin"`
	ConvertPin       bool
	Services         map[string]Service
	NESB_SDTA_NAME   string
	NESB_DDTA_NAME   string
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
func parseOneDtaParmXml(fileName string) DataTransferAdapter {
	fullPath := path.Join(getRootDir(), fileName)
	var v DataTransferAdapter
	decoder := getGbFileDecoder(fullPath)
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	trimDtaParmCDATA(&v)
	return v
}

func ParseAllDtaParmXml() map[string]DataTransferAdapter {
	m := make(map[string]DataTransferAdapter)
	files := getDtaParmFiles()
	for _, file := range files {
		dta := parseOneDtaParmXml(file)
		m[dta.Name] = dta
	}
	judgeConvertPin(m)
	parseNESB_SDTA_NAME(m)
	return m
}

func getAllConvertPinDtas() map[string]DataTransferAdapter {
	dtas := ParseAllDtaParmXml()
	svcs := ParseAllServiceXml()
	for k, v := range svcs {
		d, ok := dtas[k]
		if ok {
			d.Services = v
			dtas[k] = d
		} else {
			var dta DataTransferAdapter
			dta.ConvertPin = false
			dta.Services = v
			dtas[k] = dta
		}
	}
	return dtas
}

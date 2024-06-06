package main

import (
	"encoding/xml"
	"log"
	"path"
	"regexp"
	"strings"
	"sync"
)

type ServiceTab struct {
	XMLName  xml.Name  `xml:"ServiceTab"`
	Services []Service `xml:"Service"`
}
type Service struct {
	Name           string `xml:"Name,attr"`
	IFmt           string `xml:"IFmt,attr"`
	OFmt           string `xml:"OFmt,attr"`
	EvtIfmtBegin   string `xml:"EvtIfmtBegin"`
	EvtIfmtEnd     string `xml:"EvtIfmtEnd"`
	EvtAcallBegin  string `xml:"EvtAcallBegin"`
	NESB_SDTA_NAME string
	ConvertPin     bool
	PinElems       []PinElem
	Matched        []PinElem
	TcElems        []string
	By             string
}

func (i *Service) Clone() Service {
	var s Service
	s.Name = i.Name
	s.IFmt = i.IFmt
	s.EvtIfmtEnd = i.EvtIfmtEnd
	s.EvtAcallBegin = i.EvtAcallBegin
	s.ConvertPin = i.ConvertPin
	s.PinElems = make([]PinElem, len(i.PinElems))
	for _, v := range i.PinElems {
		s.PinElems = append(s.PinElems, v)
	}
	s.Matched = make([]PinElem, len(i.Matched))
	for _, v := range i.Matched {
		s.Matched = append(s.Matched, v)
	}
	s.TcElems = make([]string, len(i.TcElems))
	for _, v := range i.TcElems {
		s.TcElems = append(s.TcElems, v)
	}
	s.By = i.By
	return s
}

func trimServiceCDATA(st *ServiceTab) {
	sdta := regexp.MustCompile(`\$NESB_SDTA_NAME="(.*?)"`)
	tagdata := regexp.MustCompile(`nesb_get_tagdata\(.*?, *"(.*?)"\)`)
	xmlsign := regexp.MustCompile(`cbs_get_data_by_xmlsign\(.*?,.*?, *"(.*?)"\)`)
	for i, v := range st.Services {
		e := strings.TrimSpace(v.EvtIfmtBegin)
		if strings.Contains(e, "$NESB_SDTA_NAME") {
			s := sdta.FindStringSubmatch(e)
			if len(s) != 2 {
				panic(v)
			}
			v.NESB_SDTA_NAME = s[1]
		}
		v.EvtIfmtBegin = e
		e = strings.TrimSpace(v.EvtIfmtEnd)
		v.EvtIfmtEnd = e
		if strings.Contains(e, "nesbConvertPin") {
			v.ConvertPin = true
		}
		// nesb_get_tagdata("__PACKDATA", "stdpriacno|stdpindata")
		if strings.Contains(e, "nesb_get_tagdata") {
			s := tagdata.FindStringSubmatch(e)
			if len(s) != 2 {
				panic(v)
			}
			v.TcElems = strings.Split(s[1], "|")
		}
		// cbs_get_data_by_xmlsign("0", "__PACKDATA", "stdadddtap|stdpriacno|stdpindata")
		if strings.Contains(e, "cbs_get_data_by_xmlsign") {
			s := xmlsign.FindStringSubmatch(e)
			if len(s) != 2 {
				panic(v)
			}
			v.TcElems = strings.Split(s[1], "|")
		}
		v.EvtAcallBegin = strings.TrimSpace(v.EvtAcallBegin)
		st.Services[i] = v
	}
}

func serviceArrayToMap(services []Service) map[string]Service {
	m := make(map[string]Service)
	for _, v := range services {
		m[v.Name] = v
	}
	return m
}

func parseOneServiceXml(fileName string) map[string]Service {
	fullPath := path.Join(getRootDir(), fileName)
	decoder := getStarFileDecoder(fullPath)
	var v ServiceTab
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	trimServiceCDATA(&v)
	return serviceArrayToMap(v.Services)
}

func ParseAllServiceXml(wg *sync.WaitGroup) {
	defer wg.Done()
	m := make(map[string]map[string]Service)
	files := getServiceFiles()
	for dta, file := range files {
		services := parseOneServiceXml(file)
		m[dta] = services
	}
	SVCMAP = m
	log.Print("Service.xml parse success")
}

func GetServiceNamesByDta(dta string) []string {
	m, ok := SVCMAP[dta]
	if !ok {
		return nil
	}
	var s []string
	for k, _ := range m {
		s = append(s, k)
	}
	return s
}

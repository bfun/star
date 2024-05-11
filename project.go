package star

import (
	"encoding/xml"
	"path"
)

type Project struct {
	XMLName xml.Name     `xml:"Project"`
	PubDtas []ProjectDta `xml:"PubDtaTab>PubDta"`
	Apps    []ProjectApp `xml:"AppTab>App"`
}

type ProjectApp struct {
	Name      string          `xml:"Name,attr"`
	Desc      string          `xml:"Desc,attr"`
	Format    string          `xml:"Format,attr"`
	Sign      string          `xml:"Sign,attr"`
	CustCFunc string          `xml:"CustCFunc,attr"`
	Flow      string          `xml:"Flow,attr"`
	DataElem  string          `xml:"DataElem,attr"`
	SubApps   []ProjectSubApp `xml:"SubAppTab>SubApp"`
}

type ProjectSubApp struct {
	Name      string       `xml:"Name,attr"`
	Desc      string       `xml:"Desc,attr"`
	AlaParm   string       `xml:"AlaParm,attr"`
	Format    string       `xml:"Format,attr"`
	Sign      string       `xml:"Sign,attr"`
	Flow      string       `xml:"Flow,attr"`
	DataElem  string       `xml:"DataElem,attr"`
	CustCFunc string       `xml:"CustCFunc,attr"`
	SvcLogic  string       `xml:"SvcLogic,attr"`
	Dtas      []ProjectDta `xml:"DtaTab>Dta"`
}
type ProjectDta struct {
	Name      string `xml:"Name,attr"`
	Desc      string `xml:"Desc,attr"`
	DtaParm   string `xml:"DtaParm,attr"`
	CustCFunc string `xml:"CustCFunc,attr"`
	DataElem  string `xml:"DataElem,attr"`
	Format    string `xml:"Format,attr"`
	Sign      string `xml:"Sign,attr"`
	Service   string `xml:"Service,attr"`
	Route     string `xml:"Route,attr"`
}

func ParseProjectFile() Project {
	fullpath := path.Join(getRootDir(), "etc/Project.xml")
	decoder := getStarFileDecoder(fullpath)
	var v Project
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	return v
}

func getDtaParmFiles() map[string]string {
	m := make(map[string]string)
	for _, dta := range PROJECT.PubDtas {
		m[dta.Name] = dta.DtaParm
	}
	for _, v := range PROJECT.Apps {
		for _, sub := range v.SubApps {
			for _, dta := range sub.Dtas {
				m[dta.Name] = dta.DtaParm
			}
		}
	}
	return m
}

func getServiceFiles() map[string]string {
	m := make(map[string]string)
	for _, dta := range PROJECT.PubDtas {
		m[dta.Name] = dta.Service
	}
	for _, v := range PROJECT.Apps {
		for _, sub := range v.SubApps {
			for _, dta := range sub.Dtas {
				m[dta.Name] = dta.Service
			}
		}
	}
	return m
}

func getRouteFiles() map[string]string {
	m := make(map[string]string)
	for _, dta := range PROJECT.PubDtas {
		m[dta.Name] = dta.Route
	}
	for _, v := range PROJECT.Apps {
		for _, sub := range v.SubApps {
			for _, dta := range sub.Dtas {
				m[dta.Name] = dta.Route
			}
		}
	}
	return m
}

func getFormatFiles() map[string]string {
	m := make(map[string]string)
	for _, dta := range PROJECT.PubDtas {
		m[dta.Name] = dta.Format
	}
	for _, v := range PROJECT.Apps {
		for _, sub := range v.SubApps {
			for _, dta := range sub.Dtas {
				m[dta.Name] = dta.Format
			}
			m[sub.Name] = sub.Format
		}
		m[v.Name] = v.Format
	}
	return m
}

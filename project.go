package star

import (
	"encoding/xml"
	"path"
	"strings"
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

func trimFilePrefix(p *Project) {
	prefix := "file://"
	for k, v := range p.PubDtas {
		v.Name = strings.TrimPrefix(v.Name, prefix)
		v.Desc = strings.TrimPrefix(v.Desc, prefix)
		v.DtaParm = strings.TrimPrefix(v.DtaParm, prefix)
		v.CustCFunc = strings.TrimPrefix(v.CustCFunc, prefix)
		v.DataElem = strings.TrimPrefix(v.DataElem, prefix)
		v.Format = strings.TrimPrefix(v.Format, prefix)
		v.Sign = strings.TrimPrefix(v.Sign, prefix)
		v.Service = strings.TrimPrefix(v.Service, prefix)
		v.Route = strings.TrimPrefix(v.Route, prefix)
		p.PubDtas[k] = v
	}
	for ka, va := range p.Apps {
		va.Name = strings.TrimPrefix(va.Name, prefix)
		va.Desc = strings.TrimPrefix(va.Desc, prefix)
		va.Format = strings.TrimPrefix(va.Format, prefix)
		va.Sign = strings.TrimPrefix(va.Sign, prefix)
		va.CustCFunc = strings.TrimPrefix(va.CustCFunc, prefix)
		va.Flow = strings.TrimPrefix(va.Flow, prefix)
		va.DataElem = strings.TrimPrefix(va.DataElem, prefix)
		for ks, vs := range va.SubApps {
			vs.Name = strings.TrimPrefix(vs.Name, prefix)
			vs.Desc = strings.TrimPrefix(vs.Desc, prefix)
			vs.AlaParm = strings.TrimPrefix(vs.AlaParm, prefix)
			vs.Format = strings.TrimPrefix(vs.Format, prefix)
			vs.Sign = strings.TrimPrefix(vs.Sign, prefix)
			vs.Flow = strings.TrimPrefix(vs.Flow, prefix)
			vs.DataElem = strings.TrimPrefix(vs.DataElem, prefix)
			vs.CustCFunc = strings.TrimPrefix(vs.CustCFunc, prefix)
			vs.SvcLogic = strings.TrimPrefix(vs.SvcLogic, prefix)
			for kd, vd := range vs.Dtas {
				vd.Name = strings.TrimPrefix(vd.Name, prefix)
				vd.Desc = strings.TrimPrefix(vd.Desc, prefix)
				vd.DtaParm = strings.TrimPrefix(vd.DtaParm, prefix)
				vd.CustCFunc = strings.TrimPrefix(vd.CustCFunc, prefix)
				vd.DataElem = strings.TrimPrefix(vd.DataElem, prefix)
				vd.Format = strings.TrimPrefix(vd.Format, prefix)
				vd.Sign = strings.TrimPrefix(vd.Sign, prefix)
				vd.Service = strings.TrimPrefix(vd.Service, prefix)
				vd.Route = strings.TrimPrefix(vd.Route, prefix)
				vs.Dtas[kd] = vd
			}
			va.SubApps[ks] = vs
		}
		p.Apps[ka] = va
	}

}
func ParseProjectFile() Project {
	fullpath := path.Join(getRootDir(), "etc/Project.xml")
	decoder := getStarFileDecoder(fullpath)
	var v Project
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	trimFilePrefix(&v)
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

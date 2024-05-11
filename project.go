package star

import (
	"encoding/xml"
	"path"
	"strings"
)

var project Project

func init() {
	project = ParseProjectFile()
}

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

func getKindFilenamesFromProject(prefix string) []string {
	file, scanner := fileScanner("etc/Project.xml")
	defer file.Close()
	var files []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		line = strings.TrimSuffix(strings.TrimPrefix(line, prefix), `"`)
		files = append(files, line)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return files
}

func getDtaParmFiles() []string {
	return getKindFilenamesFromProject(`DtaParm="file://`)
}

func getServiceFiles() []string {
	return getKindFilenamesFromProject(`Service="file://`)
}

func getFormatFiles() []string {
	return getKindFilenamesFromProject(`Format="file://`)
}

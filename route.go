package star

import (
	"encoding/xml"
	"path"
)

type RouteTab struct {
	XMLName xml.Name `xml:"RouteTab"`
	Rules   []Rule   `xml:"RuleTab>Rule"`
}

type Rule struct {
	RuleID    string     `xml:"RuleID,attr"`
	RuleType  string     `xml:"RuleType,attr"`
	SvcExpr   string     `xml:"SvcExpr"`
	RouteExpr string     `xml:"RouteExpr"`
	Entrances []Entrance `xml:"EntranceTab>Entrance"`
}
type Entrance struct {
	Destination string `xml:"Destination,attr"`
	DstType     string `xml:"DstType,attr"`
	Expr        string `xml:"Expr"`
}

type CDATA struct {
	Value string `xml:",chardata"`
}

func parseOneRouteXml(fileName string) map[string]Entrance {
	fullPath := path.Join(getRootDir(), fileName)
	decoder := getGbFileDecoder(fullPath)
	var v RouteTab
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	m := make(map[string]Entrance)
	for _, r := range v.Rules {
		for _, e := range r.Entrances {
			m[e.Expr] = e
		}
	}
	return m
}

func ParseAllRouteXml() map[string]map[string]Entrance {
	m := make(map[string]map[string]Entrance)
	files := getRouteFiles()
	for dta, file := range files {
		r := parseOneRouteXml(file)
		m[dta] = r
	}
	return m
}

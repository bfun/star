package main

import (
	"encoding/xml"
	"log"
	"path"
	"strings"
	"sync"
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
	SvcName     string `xml:"SvcName,attr"`
	NodeName    string `xml:"NodeName,attr"`
	Expr        string `xml:"Expr"`
}

func trimSpacesFromCDATA(r *RouteTab) {
	for i, rule := range r.Rules {
		rule.SvcExpr = strings.TrimSpace(rule.SvcExpr)
		rule.RouteExpr = strings.TrimSpace(rule.RouteExpr)
		for j, entrance := range rule.Entrances {
			rule.Entrances[j].Expr = strings.TrimSpace(entrance.Expr)
		}
		r.Rules[i] = rule
	}
}

func parseOneRouteXml(fileName string) map[string]Entrance {
	fullPath := path.Join(getRootDir(), fileName)
	decoder := getStarFileDecoder(fullPath)
	var v RouteTab
	err := decoder.Decode(&v)
	if err != nil {
		panic(err)
	}
	trimSpacesFromCDATA(&v)
	m := make(map[string]Entrance)
	for _, r := range v.Rules {
		for _, e := range r.Entrances {
			m[e.Expr] = e
		}
	}
	return m
}

func ParseAllRouteXml(wg *sync.WaitGroup) {
	defer wg.Done()
	m := make(map[string]map[string]Entrance)
	files := getRouteFiles()
	for dta, file := range files {
		r := parseOneRouteXml(file)
		m[dta] = r
	}
	RUTMAP = m
	log.Print("Route.xml parse success")
}

package main

import (
	"cmp"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"sort"
	"strings"
)

type DtaSum struct {
	Name string
	Desc string
	Port string
	Node string
}

type SvcSum struct {
	Service       Service
	Route         Entrance
	Message       []string
	Request       []FmtSum
	Response      []FmtSum
	RequestItems  []TagMatch
	ResponseItems []TagMatch
}

type FmtSum struct {
	Dta string
	Svc string
	Fmt string
}

type TagMatch struct {
	SDta  string
	SSvc  string
	SFmt  string
	DDta  string
	DSvc  string
	DFmt  string
	Items []MatchItem
}

type MatchItem struct {
	SvrTag string
	Elem   string
	CltTag string
}

func svrsHandler(c *gin.Context) {
	var dtas []DtaSum
	for _, v := range ESADMIN.DtaParms {
		if !isSVR(v.DtaName) {
			continue
		}
		dta := DtaSum{Name: v.DtaName, Desc: v.DtaDesc}
		var ports []string
		for _, v := range v.IPTabItems {
			if v.Port != "" {
				ports = append(ports, v.Port)
			}
		}
		dta.Port = strings.Join(ports, ",")
		dtas = append(dtas, dta)
	}
	c.JSON(http.StatusOK, dtas)
}

func svrHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	DTANAME := strings.ToUpper(dtaName)
	dta, ok := DTAMAP[DTANAME]
	var v any
	if ok {
		v = dta
	} else {
		v = gin.H{dtaName: "not found"}
	}
	c.JSON(http.StatusOK, v)
}

func cltsHandler(c *gin.Context) {
	var dtas []DtaSum
	for _, v := range ESADMIN.DtaParms {
		if !isCLT(v.DtaName) {
			continue
		}
		dta := DtaSum{Name: v.DtaName, Desc: v.DtaDesc}
		var nodes []string
		for _, v := range v.Nodes {
			node := v.Port.NodeIP + ":" + v.Port.NodePort
			nodes = append(nodes, node)
		}
		dta.Node = strings.Join(nodes, ",")
		dtas = append(dtas, dta)
	}
	c.JSON(http.StatusOK, dtas)
}

func cltHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	DTANAME := strings.ToUpper(dtaName)
	dta, ok := DTAMAP[DTANAME]
	var v any
	if ok {
		v = dta
	} else {
		v = gin.H{dtaName: "not found"}
	}
	c.JSON(http.StatusOK, v)
}

func svcsHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	DTANAME := strings.ToUpper(dtaName)
	dta, ok := SVCMAP[DTANAME]
	var v any
	if ok {
		var s []string
		for k, _ := range dta {
			s = append(s, k)
		}
		sort.Strings(s)
		v = s
	} else {
		v = gin.H{dtaName: "not found"}
	}
	c.JSON(http.StatusOK, v)
}
func svcHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	svcName := c.Param("svc")
	DTANAME := strings.ToUpper(dtaName)
	msvc, ok := SVCMAP[DTANAME]
	var v SvcSum
	if ok {
		v.Service, ok = msvc[svcName]
		if ok {
			mrut, ok := RUTMAP[DTANAME]
			if ok {
				v.Route, ok = mrut[svcName]
				if !ok {
					v.Route, ok = mrut["^"+svcName+"$"]
					if !ok {
						v.Message = append(v.Message, fmt.Sprintf("%v.%v route not found", dtaName, svcName))
					}
				}
			} else {
				v.Message = append(v.Message, fmt.Sprintf("%v.%v route not found", dtaName, svcName))
			}
			getServiceFormat(DTANAME, svcName, &v, true)
			if v.Route.DstType == "DTA" && v.Route.Destination != "" && v.Route.SvcName != "" {
				getServiceFormat(v.Route.Destination, v.Route.SvcName, &v, false)
			}
		} else {
			v.Message = append(v.Message, fmt.Sprintf("%v.%v service not found", dtaName, svcName))
		}
	} else {
		v.Message = append(v.Message, dtaName+" not found")
	}
	if len(v.Request) > 0 {
		v.RequestItems = findMatchedTags(v.Request, true)
	}
	if len(v.Response) > 0 {
		v.ResponseItems = findMatchedTags(v.Response, false)
	}
	c.JSON(http.StatusOK, v)
}

func getServiceFormat(dta, svc string, v *SvcSum, first bool) {
	m, ok := SVCMAP[dta]
	if !ok {
		v.Message = append(v.Message, fmt.Sprintf("%v not found", dta))
		return
	}
	s, ok := m[svc]
	if !ok {
		v.Message = append(v.Message, fmt.Sprintf("%v.%v service not found", dta, svc))
		return
	}
	if first {
		v.Request = append(v.Request, FmtSum{Dta: dta, Svc: svc, Fmt: s.IFmt})
		v.Response = append(v.Response, FmtSum{Dta: dta, Svc: svc, Fmt: s.OFmt})
	} else {
		v.Request = append(v.Request, FmtSum{Dta: dta, Svc: svc, Fmt: s.OFmt})
		v.Response = append(v.Response, FmtSum{Dta: dta, Svc: svc, Fmt: s.IFmt})
	}
}

func findMatchedTags(fs []FmtSum, req bool) []TagMatch {
	if len(fs) < 2 {
		return nil
	}
	s := fs[0]
	cs := fs[1:]
	sm := findServiceElems(s.Dta, s.Svc, s.Fmt, req)
	var tms []TagMatch
	for _, c := range cs {
		cm := findServiceElems(c.Dta, c.Svc, c.Fmt, req)
		var tm TagMatch
		tm.SDta = s.Dta
		tm.SSvc = s.Svc
		tm.SFmt = s.Fmt
		tm.DDta = c.Dta
		tm.DSvc = c.Svc
		tm.DFmt = c.Fmt
		var svrOnly []MatchItem
		for se, st := range sm {
			var mi MatchItem
			mi.SvrTag = st
			mi.Elem = se
			ct, ok := cm[se]
			if ok {
				mi.CltTag = ct
				tm.Items = append(tm.Items, mi)
			} else {
				svrOnly = append(svrOnly, mi)
			}
		}
		slices.SortFunc(tm.Items, func(a, b MatchItem) int {
			return cmp.Compare(strings.ToLower(a.SvrTag), strings.ToLower(b.SvrTag))
		})
		slices.SortFunc(svrOnly, func(a, b MatchItem) int {
			return cmp.Compare(strings.ToLower(a.SvrTag), strings.ToLower(b.SvrTag))
		})
		tm.Items = append(tm.Items, svrOnly...)
		var cltOnly []MatchItem
		for ce, ct := range cm {
			_, ok := sm[ce]
			if ok {
				continue
			}
			var mi MatchItem
			mi.Elem = ce
			mi.CltTag = ct
			cltOnly = append(cltOnly, mi)
		}
		slices.SortFunc(cltOnly, func(a, b MatchItem) int {
			return cmp.Compare(strings.ToLower(a.CltTag), strings.ToLower(b.CltTag))
		})
		tm.Items = append(tm.Items, cltOnly...)
		if len(tm.Items) > 0 {
			tms = append(tms, tm)
		}
	}
	return tms
}

func rutsHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	DTANAME := strings.ToUpper(dtaName)
	dta, ok := RUTMAP[DTANAME]
	var v any
	if ok {
		v = dta
	} else {
		v = gin.H{dtaName: "not found"}
	}
	c.JSON(http.StatusOK, v)
}

func rutHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	svcName := c.Param("svc")
	DTANAME := strings.ToUpper(dtaName)
	dta, ok := RUTMAP[DTANAME]
	var v any
	if ok {
		v, ok = dta[svcName]
		if !ok {
			v = gin.H{dtaName + "." + svcName: "not found"}
		}
	} else {
		v = gin.H{dtaName: "not found"}
	}
	c.JSON(http.StatusOK, v)
}

func fmtaHandler(c *gin.Context) {
	keys := getKeysInMap(FMTMAP)
	c.JSON(http.StatusOK, keys)
}

func fmtsHandler(c *gin.Context) {
	sub := c.Param("sub")
	keys := getSomeKeysInMap(FMTMAP, sub)
	c.JSON(http.StatusOK, keys)
}

func fmtHandler(c *gin.Context) {
	dta := c.Param("dta")
	svc := c.Param("svc")
	DTA := strings.ToUpper(dta)
	fmt := c.Param("fmt")
	m := findElemsInFormat(DTA, svc, fmt)
	c.JSON(http.StatusOK, m)
}

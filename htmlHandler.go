package star

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
)

func indexHandler(c *gin.Context) {
	var dtas []DtaSum
	for _, d := range ESADMIN.DtaParms {
		if !isSVR(d.DtaName) {
			continue
		}
		dta := DtaSum{Name: d.DtaName}
		var ports []string
		for _, i := range d.IPTabItems {
			if i.Port != "" {
				ports = append(ports, i.Port)
			}
		}
		dta.Port = strings.Join(ports, ",")
		dtas = append(dtas, dta)
	}
	c.HTML(http.StatusOK, "index.html", dtas)
}

type CodesSum struct {
	DtaName string
	Port    string
	Codes   []string
}

func codesHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	DTANAME := strings.ToUpper(dtaName)
	dta, ok := SVCMAP[DTANAME]
	if !ok {
		return
	}
	dtaRuts, ok := RUTMAP[DTANAME]
	if !ok {
		return
	}
	var s []string
	for k, _ := range dta {
		rut, ok := dtaRuts[k]
		if ok {
			if rut.DstType == "ALA" {
				k += "@"
			}
		}
		s = append(s, k)
	}
	sort.Strings(s)
	var v CodesSum
	v.DtaName = dtaName
	for _, d := range ESADMIN.DtaParms {
		if d.DtaName != DTANAME {
			continue
		}
		var ports []string
		for _, i := range d.IPTabItems {
			if i.Port != "" {
				ports = append(ports, i.Port)
			}
		}
		v.Port = strings.Join(ports, ",")
	}
	v.Codes = s
	c.HTML(http.StatusOK, "codes.html", v)
}

func detailHandler(c *gin.Context) {
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
					if v.Route.DstType == "ALA" {
						flowWrapperHandler(c)
						return
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
	c.HTML(http.StatusOK, "detail.html", v)
}

type FlowSum struct {
	DtaName string
	SvcName string
}

func flowWrapperHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	svcName := c.Param("svc")
	var v = FlowSum{dtaName, svcName}
	c.HTML(http.StatusOK, "flow-wrapper.html", v)
}

func flowHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	svcName := c.Param("svc")
	var v = FlowSum{dtaName, svcName}
	c.HTML(http.StatusOK, "flow.html", v)
}

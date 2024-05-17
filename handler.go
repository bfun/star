package star

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type DtaSum struct {
	Name string
	Desc string
	Port string
	Node string
}

type SvcSum struct {
	Service  Service
	Route    Entrance
	Message  []string
	Request  []string
	Response []string
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
	// fmt.Printf("len(SVCMAP)=%d\n", len(SVCMAP))
	dta, ok := SVCMAP[DTANAME]
	// fmt.Printf("dta=%#v, ok=%v\n", dta, ok)
	var v any
	if ok {
		var s []string
		for k, _ := range dta {
			s = append(s, k)
		}
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
			getServiceFormat(DTANAME, svcName, v, true)
			if v.Route.DstType == "DTA" && v.Route.Destination != "" && v.Route.SvcName != "" {
				getServiceFormat(v.Route.Destination, v.Route.SvcName, v, false)
			}
		} else {
			v.Message = append(v.Message, fmt.Sprintf("%v.%v service not found", dtaName, svcName))
		}
	} else {
		v.Message = append(v.Message, dtaName+" not found")
	}
	c.JSON(http.StatusOK, v)
}

func getServiceFormat(dta, svc string, v SvcSum, first bool) {
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
		v.Request = append(v.Request, s.IFmt)
		v.Response = append(v.Response, s.OFmt)
	} else {
		v.Request = append(v.Request, s.OFmt)
		v.Response = append(v.Response, s.IFmt)
	}
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

package star

import (
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
	SVCNAME := strings.ToUpper(svcName)
	dta, ok := SVCMAP[DTANAME]
	var v any
	if ok {
		v, ok = dta[SVCNAME]
		if !ok {
			v = gin.H{dtaName + "." + svcName: "not found"}
		}
	} else {
		v = gin.H{dtaName: "not found"}
	}
	c.JSON(http.StatusOK, v)
}

func rutsHandler(c *gin.Context) {}

func rutHandler(c *gin.Context) {}

func fmtsHandler(c *gin.Context) {}

func fmtHandler(c *gin.Context) {}

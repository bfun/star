package star

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
)

func indexHandler(c *gin.Context) {
	var dtas []DtaSum
	for _, v := range ESADMIN.DtaParms {
		if !isSVR(v.DtaName) {
			continue
		}
		dta := DtaSum{Name: v.DtaName}
		var ports []string
		for _, v := range v.IPTabItems {
			if v.Port != "" {
				ports = append(ports, v.Port)
			}
		}
		dta.Port = strings.Join(ports, ",")
		dtas = append(dtas, dta)
	}
	c.HTML(http.StatusOK, "index.html", dtas)
}

func codesHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	DTANAME := strings.ToUpper(dtaName)
	dta, ok := SVCMAP[DTANAME]
	if !ok {
		return
	}
	var s []string
	for k, _ := range dta {
		s = append(s, k)
	}
	sort.Strings(s)
	c.HTML(http.StatusOK, "codes.html", s)
}

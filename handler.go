package star

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func dtaHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	dtas := ParseAllDtaParmXml()
	var v any
	if dtaName != "" {
		DTANAME := strings.ToUpper(dtaName)
		dta, ok := dtas[DTANAME]
		if ok {
			v = dta
		} else {
			v = gin.H{dtaName: "not found"}
		}
	} else {
		v = dtas
	}
	c.JSON(http.StatusOK, v)
}

func svcHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	svcName := c.Param("svc")
	svcs := ParseAllServiceXml()
	var v any
	if dtaName != "" {
		DTANAME := strings.ToUpper(dtaName)
		dta, ok := svcs[DTANAME]
		if ok {
			if svcName != "" {
				SVCNAME := strings.ToUpper(svcName)
				v, ok = dta[SVCNAME]
				if !ok {
					v = gin.H{svcName: "not found"}
				}
			} else {
				v = gin.H{dtaName: dta}
			}
		} else {
			v = gin.H{dtaName: "not found"}
		}
	} else {
		v = svcs
	}
	c.JSON(http.StatusOK, v)
}

package star

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func dtaHandler(c *gin.Context) {
	dtaName := c.Query("dta")
	dtas := ParseAllDtaParmXml()
	var v any
	if dtaName != "" {
		dta, ok := dtas[dtaName]
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
	dtaName := c.Query("dta")
	svcName := c.Query("svc")
	svcs := ParseAllServiceXml()
	var v any
	if dtaName != "" {
		dta, ok := svcs[dtaName]
		if ok {
			v, ok = dta[svcName]
			if !ok {
				v = gin.H{svcName: "not found"}
			}
		} else {
			v = gin.H{dtaName: "not found"}
		}
	} else {
		v = svcs
	}
	c.JSON(http.StatusOK, v)
}

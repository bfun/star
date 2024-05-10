package star

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func dtasHandler(c *gin.Context) {
	dtas := ParseAllDtaParmXml()
	var v []string
	for k, _ := range dtas {
		v = append(v, k)
	}
	c.JSON(http.StatusOK, v)
}

func dtaHandler(c *gin.Context) {
	dtaName := c.Param("dta")
	DTANAME := strings.ToUpper(dtaName)
	dtas := ParseAllDtaParmXml()
	dta, ok := dtas[DTANAME]
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
	svcs := ParseAllServiceXml()
	fmt.Printf("len(svcs)=%d\n", len(svcs))
	dta, ok := svcs[DTANAME]
	fmt.Printf("dta=%#v, ok=%v\n", dta, ok)
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
	svcs := ParseAllServiceXml()
	dta, ok := svcs[DTANAME]
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

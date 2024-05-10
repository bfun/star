package star

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func dtaHandler(c *gin.Context) {
	c.JSON(http.StatusOK, ParseAllDtaParmXml())
}

func svcHandler(c *gin.Context) {
	c.JSON(http.StatusOK, ParseAllServiceXml())
}

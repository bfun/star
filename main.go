package star

import "github.com/gin-gonic/gin"

func Main() {
	r := gin.Default()
	r.GET("/dta/:dta", dtaHandler)
	r.GET("/svc/:dta/:svc", svcHandler)
	r.Run(":8080")
}

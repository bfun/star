package star

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/svrs", svrsHandler)
	r.GET("/svr/:dta", svrHandler)
	r.GET("/clts", cltsHandler)
	r.GET("/clt/:dta", cltHandler)
	r.GET("/svcs/:dta", svcsHandler)
	r.GET("/svc/:dta/:svc", svcHandler)
	r.GET("/ruts/:dta", rutsHandler)
	r.GET("/rut/:dta/:svc", rutHandler)
	r.GET("/fmta", fmtaHandler)
	r.GET("/fmts/:sub", fmtsHandler)
	r.GET("/fmt/:dta/:svc/:fmt", fmtHandler)
	r.Run(":8000")
}

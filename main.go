package star

import "github.com/gin-gonic/gin"

func Main() {
	r := gin.Default()
	r.GET("/svrs", svrsHandler)
	r.GET("/svr/:dta", svrHandler)
	r.GET("/clts", cltsHandler)
	r.GET("/clt/:dta", cltHandler)
	r.GET("/svcs/:dta", svcsHandler)
	r.GET("/svc/:dta/:svc", svcHandler)
	r.GET("/ruts/:dta", rutsHandler)
	r.GET("/rut/:dta/:svc", rutHandler)
	r.GET("/fmts/:sub", fmtsHandler)
	r.GET("/fmt/:fmt", fmtHandler)
	r.Run(":8080")
}

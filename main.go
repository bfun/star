package star

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
)

func Main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Static("/static/", "./static")
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})
	r.LoadHTMLGlob("templates/*.html")
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
	r.GET("/codes/:dta", codesHandler)
	r.GET("/detail/:dta/:svc", detailHandler)
	r.GET("/", indexHandler)
	r.Run(":8000")
}

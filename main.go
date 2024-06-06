package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Static("/static/", "./static")
	r.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"len": func(a []string) int {
			return len(a)
		},
		"hasSuffix":  strings.HasSuffix,
		"trimSuffix": strings.TrimSuffix,
	})
	r.LoadHTMLGlob("templates/*.html")
	go func() {
		for {
			r.LoadHTMLGlob("templates/*.html")
			time.Sleep(3 * time.Second)
		}
	}()
	/*
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
	*/
	r.GET("/codes/:dta", codesHandler)
	r.GET("/detail/:dta/:svc", detailHandler)
	r.GET("/flow/:dta/:svc", flowHandler)
	r.GET("/", indexHandler)
	r.Run(":" + PORT)
}

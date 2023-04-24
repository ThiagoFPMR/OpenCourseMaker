package routes

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/controllers"
	"github.com/gin-gonic/gin"
	"html/template"
	"strings"
)

func HandleRequests() {
	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", controllers.Index)
	r.GET("/register", controllers.RegisterGETHandler)
	r.POST("/register", controllers.RegisterPOSTHandler)
	r.Run(":8000")
}

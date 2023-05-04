package routes

import (
	"html/template"
	"strings"

	"github.com/ThiagoFPMR/OpenCourseMaker/controllers"
	"github.com/gin-gonic/gin"
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
	r.GET("/login", controllers.LoginGETHandler)
	r.POST("/login", controllers.LoginPOSTHandler)
	r.GET("/player", controllers.PlayerGET)
	r.Run(":8000")
}

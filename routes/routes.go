package routes

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/controllers"
	"github.com/ThiagoFPMR/OpenCourseMaker/middlewares"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	public := r.Group("/api")
	//r.SetFuncMap(template.FuncMap{
	//	"upper": strings.ToUpper,
	//})

	//authR := r.Group("/user", controllers.Auth)
	//authR.GET("/profile", controllers.ProfileGETHandler)
	public.GET("/", controllers.Index)
	public.GET("/register", controllers.RegisterGETHandler)
	public.POST("/register", controllers.RegisterPOSTHandler)
	public.GET("/login", controllers.LoginGETHandler)
	public.POST("/login", controllers.LoginPOSTHandler)
	public.GET("/logout", controllers.LogoutGETHandler)
	public.GET("/teste", controllers.TesteGETHandler)
	public.GET("/player", controllers.PlayerGET)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.Auth())
	protected.GET("/user", controllers.CurrentUser)
	r.Run(":8000")
}

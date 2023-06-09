package routes

import (
	"html/template"

	"github.com/ThiagoFPMR/OpenCourseMaker/controllers"
	"github.com/ThiagoFPMR/OpenCourseMaker/middlewares"
	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"GetLoggedInStatus": GetLoggedInStatus,
	})
	r.LoadHTMLGlob("templates/*.html")

	public := r.Group("/")
	{
		public.GET("/", controllers.Index)
		public.GET("/register", controllers.RegisterGETHandler)
		public.POST("/register", controllers.RegisterPOSTHandler)
		public.GET("/login", controllers.LoginGETHandler)
		public.POST("/login", controllers.LoginPOSTHandler)

	}

	authMiddleware := middlewares.JwtAuthMiddleware()

	protected := r.Group("/")
	protected.Use(authMiddleware)
	protected.GET("/user", controllers.CurrentUser)
	protected.GET("/logout", controllers.LogoutGETHandler)
	protected.GET("/dashboard", controllers.DashboardGETHandler)
	protected.POST("/create_course", controllers.CreateCoursePOSTHandler)
	protected.GET("/cursos/:id", controllers.CurseseInfoGETHandler)
	protected.GET("/enroll/:id", controllers.EnrollHandler)
	protected.GET("/new_topic/:id", controllers.NewTopicGETHandler)
	protected.POST("/course/:id/topic", controllers.AddTopico)
	protected.GET("/courses/:cursoId/topico/:topicoId/editar", controllers.EditTopicGETHandler)
	protected.POST("/courses/:cursoId/topico/:topicId/salvar", controllers.SaveTopicPOSTHandler)
	protected.POST("/excluir/:cursoId/topico/:topicoID", controllers.DeleteTopicPOSTHandler)
	protected.GET("/player", controllers.PlayerGET)
	r.Static("/ide", "./templates/ide")
	r.Run(":8000")
}

func GetLoggedInStatus(c *gin.Context) bool {
	logged_in, err := c.Cookie("logged_in")
	if err != nil || logged_in != "true" {
		return false
	}
	return true
}

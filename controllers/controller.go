package controllers

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/db"
	"github.com/ThiagoFPMR/OpenCourseMaker/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func RegisterGETHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func RegisterPOSTHandler(c *gin.Context) {
	var user models.User
	user.Nome = c.PostForm("nome")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	password2 := c.PostForm("password2")

	if user.Password != password2 {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "As senhas não conferem",
		})
		return
	}

	res := db.BD.Create(&user)
	if res.Error != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	location := url.URL{Path: "/"}
	c.Redirect(http.StatusMovedPermanently, location.RequestURI())

}

//func Autenticar(w http.ResponseWriter, r *http.Request) {
//}

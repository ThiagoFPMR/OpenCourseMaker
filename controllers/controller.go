package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ThiagoFPMR/OpenCourseMaker/db"
	"github.com/ThiagoFPMR/OpenCourseMaker/user"
	"github.com/ThiagoFPMR/OpenCourseMaker/user/login"
	"github.com/ThiagoFPMR/OpenCourseMaker/user/signup"
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func RegisterGETHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func RegisterPOSTHandler(c *gin.Context) {
	nome := c.PostForm("nome")
	email := c.PostForm("email")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")
	if password != password2 {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "As senhas não conferem",
		})
		return
	}

	res, err := signup.Signup(db.BD, &signup.Request{
		Nome:     nome,
		Email:    email,
		Password: password,
	})

	if err != nil {
		switch err.(type) {
		case *user.EmailAlreadyExistsError:
			c.HTML(http.StatusBadRequest, err.Error(), nil)
			return
		default:
			c.HTML(http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}

	fmt.Println("Created: ", res.Id)

	location := url.URL{Path: "/"}
	c.Redirect(http.StatusMovedPermanently, location.RequestURI())

}

func LoginGETHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func LoginPOSTHandler(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	res, err := login.Login(db.BD, &login.Request{
		Email:    email,
		Password: password,
	})
	if err != nil {
		switch err.(type) {
		case *login.PasswordIncorrectError:
			c.HTML(http.StatusBadRequest, err.Error(), nil)
			return
		default:
			c.HTML(http.StatusInternalServerError, err.Error(), nil)
			return
		}
	}
	fmt.Println("Logged in: ", res.User.Email)
}

func PlayerGET(c *gin.Context) {
	c.HTML(http.StatusOK, "player.html", nil)
}

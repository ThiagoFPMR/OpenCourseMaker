package controllers

import (
	"fmt"
	"github.com/ThiagoFPMR/OpenCourseMaker/services"

	"github.com/ThiagoFPMR/OpenCourseMaker/db"
	"github.com/ThiagoFPMR/OpenCourseMaker/user"
	"github.com/ThiagoFPMR/OpenCourseMaker/user/login"
	"github.com/ThiagoFPMR/OpenCourseMaker/user/signup"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

const userkey = "user"

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Auth(c *gin.Context) {
	//session, _ := store.Get(c.Request, "session")
	//user, password, _ := c.Request.BasicAuth()

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
	} else {
		//access_token, err := user.GenerateToken(res.User.ID)
		token, err := services.NewJWTService().GenerateToken(res.User.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		//c.JSON(200, gin.H{"token": token})
		//if err != nil {
		//	c.JSON(http.StatusInternalServerError, err.Error())
		//	return
		//}

		c.SetCookie("access_token", token, 3600, "/", "localhost", false, true)
		c.SetCookie("logged_in", "true", 3600, "/", "localhost", false, true)
		//c.JSON(http.StatusOK, gin.H{"message": "success", "access_token": access_token})
		c.JSON(http.StatusOK, gin.H{"message": "success", "access_token": token})
		c.Redirect(http.StatusOK, "/teste")
	}
}

func LogoutGETHandler(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func CurrentUser(c *gin.Context) {
	user_id, err := user.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	u, err := user.FindById(db.BD, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

func PlayerGET(c *gin.Context) {
	c.HTML(http.StatusOK, "player.html", nil)
}

func TesteGETHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "teste.html", nil)
}

package controllers

import (
	"fmt"
	"github.com/ThiagoFPMR/OpenCourseMaker/course"
	"github.com/ThiagoFPMR/OpenCourseMaker/course/newCourse"
	"github.com/ThiagoFPMR/OpenCourseMaker/services"
	"gorm.io/gorm"
	"strconv"

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
	logged_in := GetLoggedInStatus(c)

	cursos, err := course.GetAllCursos(db.BD)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"logged_in": logged_in,
		"cursos":    cursos,
	})
}

func RegisterGETHandler(c *gin.Context) {
	logged_in := GetLoggedInStatus(c)
	c.HTML(http.StatusOK, "register.html", gin.H{
		"logged_in": logged_in,
	})
}

func RegisterPOSTHandler(c *gin.Context) {
	nome := c.PostForm("nome")
	email := c.PostForm("email")
	tipo := c.PostForm("tipo")
	password := c.PostForm("password")
	password2 := c.PostForm("password2")
	if password != password2 {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "As senhas não conferem",
		})
		return
	}

	tipoFormatado, err := strconv.Atoi(tipo)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"error": "Tipo de usuário inválido",
		})
		return
	}

	res, err := signup.Signup(db.BD, &signup.Request{
		Nome:     nome,
		Email:    email,
		Tipo:     tipoFormatado,
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
		token, err := services.NewJWTService().GenerateToken(res.User.ID, res.User.Nome, res.User.Email)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}

		c.SetCookie("access_token", token, 3600, "/", "localhost", false, true)
		c.SetCookie("logged_in", "true", 3600, "/", "localhost", false, true)
		c.Redirect(http.StatusSeeOther, "/dashboard")
	}
}

func LogoutGETHandler(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("logged_in", "false", -1, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/")
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

func DashboardGETHandler(c *gin.Context) {
	logged_in := GetLoggedInStatus(c)

	// Obtém o nome de usuário e o e-mail do usuário atual da variável de contexto
	nome, exists := c.Get("nome")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	email, exists := c.Get("email")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, _ := user.FindByEmail(db.BD, email.(string))

	cursos, err := course.GetCursosById(db.BD, user.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	cursosMatriculados, err := course.FindCursosByStudentID(db.BD, user.ID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tipo := user.Tipo

	var tipoConta string
	if tipo == 1 {
		tipoConta = "Professor"
	} else {
		tipoConta = "Aluno"
	}

	// Renderiza a página de dashboard com as informações do usuário atual
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"nome":       nome,
		"email":      email,
		"tipo":       tipoConta,
		"cursos":     cursos,
		"matriculas": cursosMatriculados,
		"logged_in":  logged_in,
	})
}

func CreateCoursePOSTHandler(c *gin.Context) {
	// Obtém os valores enviados pelo formulário
	nome := c.PostForm("nome")
	descricao := c.PostForm("descricao")
	carga_horaria := c.PostForm("carga_horaria")

	email, exists := c.Get("email")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, _ := user.FindByEmail(db.BD, email.(string))
	professorID := user.ID

	// Cria o curso
	res, err := newCourse.NewCourse(db.BD, &newCourse.Request{
		Nome:         nome,
		Descricao:    descricao,
		CargaHoraria: carga_horaria,
		ProfessorID:  professorID,
	})

	if err != nil {
		c.HTML(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	fmt.Println("Created: ", res.Id)

	location := url.URL{Path: "/dashboard"}
	c.Redirect(http.StatusMovedPermanently, location.RequestURI())
}

func CurseseInfoGETHandler(c *gin.Context) {
	logged_in := GetLoggedInStatus(c)
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var id uint = uint(id64)

	curso, err := course.FindCursoById(db.BD, id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	professor_responsavel, err := user.FindById(db.BD, curso.ProfessorID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println(curso)

	c.HTML(http.StatusOK, "course_info.html", gin.H{
		"curso":                 curso,
		"professor_responsavel": professor_responsavel,
		"logged_in":             logged_in,
	})
}

func EnrollHandler(c *gin.Context) {
	logged_in := GetLoggedInStatus(c)
	currentUserID, _ := user.ExtractTokenID(c)
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var courseID uint = uint(id64)
	curso, err := course.FindCursoById(db.BD, courseID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Verificar se o usuário já está matriculado no curso
	enrollment, err := course.FindEnrollmentByCourseAndStudent(db.BD, courseID, currentUserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if enrollment == nil {
		// Criar a matrícula
		newEnrollment := &course.Enrollment{
			CursoID: courseID,
			AlunoID: currentUserID,
		}
		err = db.BD.Create(newEnrollment).Error
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.HTML(http.StatusOK, "enroll.html", gin.H{
		"curso":     curso,
		"logged_in": logged_in,
	})
}

func GetLoggedInStatus(c *gin.Context) bool {
	logged_in, err := c.Cookie("logged_in")
	if err != nil || logged_in != "true" {
		return false
	}
	return true
}

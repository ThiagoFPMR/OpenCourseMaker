package controllers

import (
	"fmt"
	"strconv"

	"github.com/ThiagoFPMR/OpenCourseMaker/course"
	"github.com/ThiagoFPMR/OpenCourseMaker/course/newCourse"
	"github.com/ThiagoFPMR/OpenCourseMaker/services"
	"gorm.io/gorm"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/ThiagoFPMR/OpenCourseMaker/db"
	"github.com/ThiagoFPMR/OpenCourseMaker/user"
	"github.com/ThiagoFPMR/OpenCourseMaker/user/login"
	"github.com/ThiagoFPMR/OpenCourseMaker/user/signup"
	"github.com/gin-gonic/gin"
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
	user, err := user.FindById(db.BD, currentUserID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println(user)
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

	// Verificar se o usuário já está matriculado no curso e se é um professor
	if user.Tipo == 2 {
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
	}

	topicos, err := course.GetTopicosByIdCurso(db.BD, courseID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Verificar se o usuário é o criador do curso
	isOwner := currentUserID == curso.ProfessorID

	c.HTML(http.StatusOK, "enroll.html", gin.H{
		"curso":     curso,
		"usuario":   user,
		"topicos":   topicos,
		"logged_in": logged_in,
		"is_owner":  isOwner,
	})
}

func NewTopicGETHandler(c *gin.Context) {
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

	c.HTML(http.StatusOK, "new_topic.html", gin.H{
		"curso":     curso,
		"logged_in": logged_in,
	})
}

func AddTopico(c *gin.Context) {
	currentUserID, _ := user.ExtractTokenID(c)

	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var courseID uint = uint(id64)

	// Verifica se o curso existe
	curso, err := course.FindCursoById(db.BD, uint(courseID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Verifica se o usuário é o professor responsável pelo curso
	if curso.ProfessorID != currentUserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse dos dados do formulário
	titulo := c.PostForm("titulo")
	videoURL := c.PostForm("video_url")
	descricao := c.PostForm("descricao")

	videoID := extractVideoID(videoURL)

	// Cria o novo tópico
	topico := course.Topico{
		Titulo:      titulo,
		VideoURL:    &videoID,
		Desc:        &descricao,
		CursoID:     curso.ID,
		ProfessorID: curso.ProfessorID,
		CriadoEm:    time.Now(),
	}

	// Salva o novo tópico no banco de dados
	if err := db.BD.Create(&topico).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create topic"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Topic created successfully"})
}

func EditTopicGETHandler(c *gin.Context) {
	logged_in := GetLoggedInStatus(c)

	// Obter o ID do curso da URL
	cursoID := c.Param("cursoId")

	// Obter o ID do tópico da URL
	topicID := c.Param("topicoId")

	// Recuperar as informações do tópico do banco de dados usando o ID
	var topic course.Topico
	err := db.BD.Where("id = ?", topicID).First(&topic).Error
	if err != nil {
		c.HTML(http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// Renderizar o formulário de edição com os dados do tópico
	c.HTML(http.StatusOK, "edit_topic.html", gin.H{
		"cursoid":   cursoID,
		"titulo":    topic.Titulo,
		"videoURL":  topic.VideoURL,
		"desc":      topic.Desc,
		"topicID":   topicID,
		"logged_in": logged_in,
	})
}

func SaveTopicPOSTHandler(c *gin.Context) {
	// Redirecionar para a página do curso
	cursoID := c.Param("cursoId")
	cursoIDUint, err := strconv.ParseUint(cursoID, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Recuperar o ID do tópico
	topicID := c.Param("topicId")
	topicIDUint, err := strconv.ParseUint(topicID, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Recuperar os dados do formulário
	titulo := c.PostForm("titulo")
	videoURL := c.PostForm("video_url")
	descricao := c.PostForm("descricao")

	// Atualizar o tópico no banco de dados
	err = course.UpdateTopic(db.BD, uint(cursoIDUint), uint(topicIDUint), titulo, videoURL, descricao)
	if err != nil {
		// Caso ocorra algum erro ao atualizar o tópico
		// Retornar um erro 500 (Internal Server Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/enroll/%s", cursoID))
}

func DeleteTopicPOSTHandler(c *gin.Context) {
	// Recuperar o ID do curso
	cursoID := c.Param("cursoId")
	cursoIDUint, err := strconv.ParseUint(cursoID, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Recuperar o ID do tópico
	topicoID := c.Param("topicoID") // alterar de "topicId" para "topicoID"
	topicoIDUint, err := strconv.ParseUint(topicoID, 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Excluir o tópico do banco de dados
	err = course.DeleteTopic(db.BD, uint(cursoIDUint), uint(topicoIDUint)) // alterar de "topicId" para "topicoID"
	if err != nil {
		// Caso ocorra algum erro ao excluir o tópico
		// Retornar um erro 500 (Internal Server Error)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/enroll/%s", cursoID))
}

func GetLoggedInStatus(c *gin.Context) bool {
	logged_in, err := c.Cookie("logged_in")
	if err != nil || logged_in != "true" {
		return false
	}
	return true
}

func extractVideoID(videoURL string) string {
	videoID := strings.TrimSpace(strings.TrimPrefix(videoURL, "https://www.youtube.com/watch?v="))
	last11Chars := videoID[len(videoID)-11:]
	return last11Chars
}

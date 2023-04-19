package routes

import (
	"OpenCourseMaker/controllers"
	"net/http"
)

func CarregaRotas() {
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/novocadastro", controllers.NovoCadastro)
	http.HandleFunc("/cadastrar", controllers.Cadastrar)
	http.HandleFunc("/login", controllers.Autenticar)
}

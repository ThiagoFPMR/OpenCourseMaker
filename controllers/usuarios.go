package controllers

import (
	"OpenCourseMaker/models"
	"log"
	"net/http"
	"strconv"
)

func NovoCadastro(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "NovoCadastro", nil)
}

func Cadastrar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		nome := r.FormValue("nome")
		sobrenome := r.FormValue("sobrenome")
		datanasc := r.FormValue("datanasc")
		tipo := r.FormValue("tipo")
		password := r.FormValue("password")
		password2 := r.FormValue("password2")

		if password != password2 {
			log.Println("Senhas não conferem")
			http.Redirect(w, r, "/", 301)
		}

		tipoConvertidoParaInt, err := strconv.Atoi(tipo)
		if err != nil {
			log.Println("Erro ao converter a data para int", err)
			http.Redirect(w, r, "/", 301)
		}

		if password == password2 {
			models.CadastrarUsuario(username, nome, sobrenome, datanasc, password, tipoConvertidoParaInt)
			http.Redirect(w, r, "/", 301)
		}
	}
}

func Autenticar(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "Login", nil)
}

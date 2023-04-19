package models

import "OpenCourseMaker/db"

type Usuario struct {
	Id, Tipo  int
	Nome      string
	Sobrenome string
	Datanasc  string
	username  string
	password  string
}

func CadastrarUsuario(username, nome, sobrenome, datanasc, password string, tipo int) {
	db := db.ConectaBD()

	insert, err := db.Prepare("insert into usuarios(nome, sobrenome, datanasc, tipo, username, password) values($1, $2, $3, $4, $5, $6)")
	if err != nil {
		panic(err.Error())
	}

	insert.Exec(nome, sobrenome, datanasc, tipo, username, password)
	defer db.Close()
}

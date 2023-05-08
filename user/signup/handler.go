package signup

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/models"
	"github.com/ThiagoFPMR/OpenCourseMaker/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Request struct {
	Nome     string
	Email    string
	Tipo     int
	Password string
}

type Response struct {
	Id uint
}

func Signup(db *gorm.DB, req *Request) (*Response, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := &models.User{
		Nome:     req.Nome,
		Email:    req.Email,
		Tipo:     req.Tipo,
		Pswdhash: string(passwordHash),
	}
	id, err := user.Create(db, newUser)
	if err != nil {
		return nil, err
	}
	return &Response{Id: id}, nil
}

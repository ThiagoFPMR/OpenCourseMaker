package login

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/models"
	"github.com/ThiagoFPMR/OpenCourseMaker/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Request struct {
	Email    string
	Password string
}

type Response struct {
	User *models.User
}

type PasswordIncorrectError struct{}

func (e *PasswordIncorrectError) Error() string {
	return "Password incorrect"
}

func Login(db *gorm.DB, req *Request) (*Response, error) {
	user, err := user.FindByEmail(db, req.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Pswdhash), []byte(req.Password))
	if err != nil {
		return nil, &PasswordIncorrectError{}
	}
	return &Response{User: user}, nil
}

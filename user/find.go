package user

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/models"
	"gorm.io/gorm"
)

type EmailNotFoundError struct{}

func (e *EmailNotFoundError) Error() string {
	return "Email not found"
}

func FindByEmail(db *gorm.DB, email string) (*models.User, error) {
	var user models.User
	res := db.Find(&user, &models.User{Email: email})
	if res.Error != nil {
		return nil, &EmailNotFoundError{}
	}
	return &user, nil
}

func FindById(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	res := db.Find(&user, &models.User{ID: id})
	if res.Error != nil {
		return nil, &EmailNotFoundError{}
	}
	return &user, nil
}

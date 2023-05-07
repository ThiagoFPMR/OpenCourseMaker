package user

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/models"
	"gorm.io/gorm"
)

func Create(db *gorm.DB, user *models.User) (uint, error) {
	err := db.Create(user).Error
	if err != nil {
		//if postgres.IsUniqueConstraintError(err, UniqueConstraintUserEmail) {
		//	return 0, &EmailAlreadyExistsError{Email: user.Email}
		//}
		return 0, err
	}
	return user.ID, nil
}

package course

import (
	"gorm.io/gorm"
)

func Create(db *gorm.DB, curso *Curso) (uint, error) {
	err := db.Create(curso).Error
	if err != nil {
		//if postgres.IsUniqueConstraintError(err, UniqueConstraintUserEmail) {
		//	return 0, &EmailAlreadyExistsError{Email: user.Email}
		//}
		return 0, err
	}
	return curso.ID, nil
}

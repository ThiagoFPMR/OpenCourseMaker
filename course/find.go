package course

import (
	"gorm.io/gorm"
)

func GetAllCursos(db *gorm.DB) ([]Curso, error) {
	var cursos []Curso
	if err := db.Find(&cursos).Error; err != nil {
		return nil, err
	}

	return cursos, nil
}

func GetCursosById(db *gorm.DB, id uint) ([]Curso, error) {
	var cursos []Curso
	if err := db.Find(&cursos, &Curso{ProfessorID: id}).Error; err != nil {
		return nil, err
	}

	return cursos, nil
}

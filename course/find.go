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

func FindCursoById(db *gorm.DB, id uint) (*Curso, error) {
	var curso Curso
	res := db.Find(&curso, &Curso{ID: id})
	if res.Error != nil {
		return nil, res.Error
	}
	return &curso, nil
}

func FindEnrollmentByCourseAndStudent(db *gorm.DB, courseID, studentID uint) (*Enrollment, error) {
	var enrollment Enrollment
	res := db.Where("curso_id = ? AND aluno_id = ?", courseID, studentID).First(&enrollment)
	if res.Error != nil {
		return nil, res.Error
	}
	return &enrollment, nil
}

func FindCursosByStudentID(db *gorm.DB, studentID uint) ([]Curso, error) {
	var cursos []Curso
	res := db.Joins("JOIN enrollments ON enrollments.curso_id = cursos.id").
		Where("enrollments.aluno_id = ?", studentID).
		Find(&cursos)
	if res.Error != nil {
		return nil, res.Error
	}
	return cursos, nil
}

func GetTopicosByIdCurso(db *gorm.DB, id uint) ([]Topico, error) {
	var topicos []Topico
	if err := db.Find(&topicos, &Topico{CursoID: id}).Error; err != nil {
		return nil, err
	}

	return topicos, nil
}

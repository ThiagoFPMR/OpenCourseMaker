package newCourse

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/course"
	"gorm.io/gorm"
)

type Request struct {
	Nome         string
	Descricao    string
	CargaHoraria string
	ProfessorID  uint
}

type Response struct {
	Id uint
}

func NewCourse(db *gorm.DB, req *Request) (*Response, error) {
	newCourse := &course.Curso{
		Nome:         req.Nome,
		Descricao:    req.Descricao,
		CargaHoraria: req.CargaHoraria,
		ProfessorID:  req.ProfessorID,
	}
	id, err := course.Create(db, newCourse)
	if err != nil {
		return nil, err
	}
	return &Response{Id: id}, nil
}

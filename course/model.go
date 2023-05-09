package course

import (
	"time"
)

type Curso struct {
	ID           uint   `gorm:"primary_key"`
	Nome         string `gorm:"not null"`
	Descricao    string
	CargaHoraria string `gorm:"not null"`
	ProfessorID  uint
	Topicos      []Topico
}

type Topico struct {
	ID       uint      `gorm:"primary_key"`
	Titulo   string    `gorm:"not null"`
	VideoURL *string   // campo de vídeo é opcional
	Desc     *string   // campo de descrição é opcional
	CursoID  uint      `gorm:"not null"`
	CriadoEm time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Enrollment struct {
	ID      uint   `gorm:"primary_key"`
	AlunoID uint   `gorm:"not null"`
	CursoID uint   `gorm:"not null"`
	Status  string `gorm:"not null"`
}

package db

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/course"
	"github.com/ThiagoFPMR/OpenCourseMaker/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var (
	BD  *gorm.DB
	err error
)

func ConectBD() {
	conect := "user=postgres dbname=open_course_maker password=postgres host=localhost sslmode=disable"
	con, err := gorm.Open(postgres.Open(conect), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados")
	}
	con.AutoMigrate(&models.User{})
	con.AutoMigrate(&course.Topico{})
	con.AutoMigrate(&course.Curso{})
	con.AutoMigrate(&course.Enrollment{})

	BD = con
}

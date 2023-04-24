package db

import (
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
	DB, err := gorm.Open(postgres.Open(conect))
	if err != nil {
		log.Panic("Erro ao conectar com o banco de dados")
	}
	DB.AutoMigrate(&models.User{})
}

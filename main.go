package main

import (
	"github.com/ThiagoFPMR/OpenCourseMaker/db"
	"github.com/ThiagoFPMR/OpenCourseMaker/routes"
)

func main() {
	db.ConectBD()
	routes.HandleRequests()
}

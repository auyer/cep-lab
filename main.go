package main

import (
	"pessoalAPI-gingonic/controllers"
	"pessoalAPI-gingonic/db"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db.Init()
	defer db.GetDB().Db.Close()
	pessoal := new(controllers.PessoalController) //Controller instance

	r.GET("/api/servidores", pessoal.GetPessoal)             //Simple route
	r.GET("/api/servidor/:matricula", pessoal.GetPessoalMat) //Route with URL parameter

	r.Run() // listen and serve on 0.0.0.0:8080
}

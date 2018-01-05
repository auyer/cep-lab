package main

import (
	"pessoalAPI-gingonic/controllers"
	"pessoalAPI-gingonic/db"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	db.Init()
	pessoal := new(controllers.PessoalController)

	r.GET("/pessoal", pessoal.GetPessoal)
	//r.POST("/pong", ping.Pong)

	r.Run() // listen and serve on 0.0.0.0:8080
}

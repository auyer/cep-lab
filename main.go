package main

import (
	"github.com/latitude-RESTsec-lab/api-gingonic/controllers"
	"github.com/latitude-RESTsec-lab/api-gingonic/db"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db.Init()
	defer db.GetDB().Db.Close()
	servidor := new(controllers.ServidorController) //Controller instance

	router.GET("/api/servidores", servidor.GetServidores)           //Simple route
	router.GET("/api/servidor/:matricula", servidor.GetServidorMat) //Route with URL parameter
	router.POST("/api/servidor/", servidor.PostServidor)

	router.Run() // listen and serve on 0.0.0.0:8080
}

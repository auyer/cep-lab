package main

import (
	"github.com/latitude-RESTsec-lab/api-gingonic/controllers"
	"github.com/latitude-RESTsec-lab/api-gingonic/db"

	"github.com/gin-gonic/gin"
)

func main() {
	httpsRouter := gin.Default()

	db.Init()
	defer db.GetDB().Db.Close()
	servidor := new(controllers.ServidorController) //Controller instance

	httpsRouter.GET("/api/servidores", servidor.GetServidores)           //Simple route
	httpsRouter.GET("/api/servidor/:matricula", servidor.GetServidorMat) //Route with URL parameter
	httpsRouter.POST("/api/servidor/", servidor.PostServidor)

	httpsRouter.RunTLS(":443", "./devssl/server.pem", "./devssl/server.key") // listen and serve on 0.0.0.0:8080
}

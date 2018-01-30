package main

import (
	"fmt"
	"log"
	"os"

	"github.com/latitude-RESTsec-lab/api-gingonic/controllers"
	"github.com/latitude-RESTsec-lab/api-gingonic/db"

	"github.com/gin-gonic/gin"
)

func main() {
	file, fileErr := os.Create("server.log")
	if fileErr != nil {
		fmt.Println(fileErr)
		file = os.Stdout
	}
	log.SetOutput(file)
	gin.DefaultWriter = file
	gin.SetMode(gin.ReleaseMode)
	// BEGIN HTTPS

	httpsRouter := gin.Default()

	db.Init()
	defer db.GetDB().Db.Close()
	servidor := new(controllers.ServidorController) //Controller instance

	httpsRouter.GET("/api/servidores", servidor.GetServidores)           //Simple route
	httpsRouter.GET("/api/servidor/:matricula", servidor.GetServidorMat) //Route with URL parameter
	httpsRouter.POST("/api/servidor/", servidor.PostServidor)

	// BEGIN HTTP
	httpRouter := gin.Default()

	httpRouter.GET("/api/servidores/", func(c *gin.Context) {
		c.Redirect(302, fmt.Sprint("https://", c.Request.Host, ".", c.Request.URL.Path))
	})
	httpRouter.GET("/api/servidor/:matricula", func(c *gin.Context) {
		c.Redirect(302, fmt.Sprint("https://", c.Request.Host, ".", c.Request.URL.Path))
	})

	go httpRouter.Run(":80")
	httpsRouter.RunTLS(":443", "./devssl/server.pem", "./devssl/server.key") // listen and serve on 0.0.0.0:8080
}

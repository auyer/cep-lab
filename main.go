package main

import (
	"fmt"
	"log"

	"github.com/latitude-RESTsec-lab/api-gingonic/controllers"
	"github.com/latitude-RESTsec-lab/api-gingonic/db"

	"github.com/gin-gonic/gin"
)

func main() {

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

	httpRouter.GET("/", func(c *gin.Context) {
		log.Print(c.Request.Host)
		log.Print(c.Request.URL.Path)
		c.JSON(200, nil)
	})
	httpRouter.GET("/api/servidores/", func(c *gin.Context) {
		c.Redirect(302, fmt.Sprint("https://", c.Request.Host, ".", c.Request.URL.Path))
	})
	httpRouter.GET("/api/servidor/:matricula", func(c *gin.Context) {
		c.Redirect(302, fmt.Sprint("https://", c.Request.Host, ".", c.Request.URL.Path))
	})

	go httpRouter.Run(":80")
	httpsRouter.RunTLS(":443", "./devssl/server.pem", "./devssl/server.key") // listen and serve on 0.0.0.0:8080
}

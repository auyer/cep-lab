package main

import (
	"fmt"
	"log"

	"github.com/restsec/api-gingonic/config"
	"github.com/restsec/api-gingonic/controllers"
	"github.com/restsec/api-gingonic/db"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Starting Gin Gonic API")
	err := config.ReadConfig()
	if err != nil {
		fmt.Print("Error reading configuration file")
		log.Print(err.Error())
		return
	}

	log.SetOutput(config.LogFile)
	gin.DefaultWriter = config.LogFile
	if config.ConfigParams.Debug != "true" {
		gin.SetMode(gin.ReleaseMode)
	}
	// BEGIN HTTPS

	httpsRouter := gin.Default()

	db.Init()
	defer db.GetDB().Db.Close()
	servidor := new(controllers.ServidorController) //Controller instance

	httpsRouter.GET("/api/servidores", servidor.GetServidores)           //Simple route
	httpsRouter.GET("/api/servidor/:matricula", servidor.GetServidorMat) //Route with URL parameter
	httpsRouter.POST("/api/servidor/", servidor.PostServidor)
	httpsRouter.POST("/api/calculo/", servidor.Calculate)

	// BEGIN HTTP
	// httpRouter := gin.Default()

	// httpRouter.GET("/api/servidores/", func(c *gin.Context) {
	// 	c.Redirect(302, fmt.Sprint("https://", c.Request.Host, ".", c.Request.URL.Path))
	// })
	// httpRouter.GET("/api/servidor/:matricula", func(c *gin.Context) {
	// 	c.Redirect(302, fmt.Sprint("https://", c.Request.Host, ".", c.Request.URL.Path))
	// })

	// go httpRouter.Run(":" + config.ConfigParams.HttpPort)
	err = httpsRouter.RunTLS(":"+config.ConfigParams.HttpsPort, config.ConfigParams.TLSCertLocation, config.ConfigParams.TLSKeyLocation) // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return
	}
}

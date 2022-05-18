// Package main contains: http app, mqtt app, swagger doc, ORM setup
package main

import (
	"fmt"
	"os"

	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	"github.com/ecoprohcm/DMS_BackendServer/initializers"
	"github.com/gin-gonic/gin"                 // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title          Backend API
// @version        1.0
// @description    This is API document for DMS backend server
// @contact.name   Mr. Khai

// @host      http://iot.hcmue.space:8002
// @BasePath  /v1

func main() {
	cc, _, err := initializers.InitApplication("./.env")
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}
	// HTTP Serve
	r := handlers.SetupRouter(cc.HandlerOptions)
	initSwagger(r)
	r.Run(":8080")
	cc.MqttClient.Disconnect(250)
}

func initSwagger(r *gin.Engine) {
	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

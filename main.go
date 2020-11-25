package main

import (
	"github.com/0sman/godemo/app/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	controller.InitController()

	r.GET("/histories", controller.FindAllHistories)
	r.GET("/histories/:id", controller.FindHistory)
	r.PUT("/histories/:id", controller.UpdateHistory)
	r.POST("/histories", controller.CreateHistory)

	r.GET("/gi", controller.FindAllGeneralInformations)
	r.GET("/gi/:id", controller.FindGeneralInformation)
	r.PUT("/gi/:id", controller.UpdateGeneralInformation)
	r.POST("/gi", controller.CreateGeneralInformation)

	r.GET("/users", controller.FindAllUsers)
	r.GET("/users/:id", controller.FindUser)
	r.PUT("/users/:id", controller.UpdateUser)
	r.POST("/users", controller.CreateUser)

	r.POST("/auth", controller.AuthUser)

	r.Run()
}

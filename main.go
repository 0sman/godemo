package main

import (
	"net/http"

	"github.com/0sman/godemo/app/controller"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	controller.InitConttroller()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	r.GET("/histories", controller.FindHistories)
	r.PUT("/histories/:id", controller.UpdateHistory)
	r.POST("/histories", controller.CreateHistory)

	r.GET("/gi", controller.FindGeneralInformations)
	r.PUT("/gi/:id", controller.UpdateGeneralInformation)
	r.POST("/gi", controller.CreateGeneralInformation)

	r.GET("/users", controller.FindUsers)
	r.PUT("/users/:id", controller.UpdateUser)
	r.POST("/users", controller.CreateUser)

	r.Run()

	/*	d := service.ReadSecuredModel(appmodel.History{})
		fmt.Println("output:", d)

		st := "new name"
		id := 1
		service.UpdateSecuredModel(id, appmodel.History{CourseName: &st})

		stt := "new insert name"
		service.CreateSecuredModel(appmodel.History{CourseName: &stt})
	*/
}

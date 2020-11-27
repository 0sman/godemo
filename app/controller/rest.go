package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/0sman/godemo/app/appmodel"
	"github.com/0sman/godemo/app/appservice"
	"github.com/0sman/godemo/perm/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	dataSourceName := "root:osman123@tcp(localhost:3306)/godemo?parseTime=True"
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	/*
		// db.Exec("CREATE DATABASE godemo")
		// db.Exec("USE godemo")

		 db.AutoMigrate(
			&appmodel.GeneralInformation{},
			&appmodel.History{},
			&appmodel.User{})

		db.AutoMigrate(
			&dal.Object{},
			&dal.Owner{},
			&dal.Column{},
			&dal.Group{},
			&dal.Permission{})
	*/

	return db
}

func InitController() {
	var db = initDB()
	service.InitService(db)
	appservice.InitService(db)
}

func FindAllHistories(c *gin.Context) {
	validateToken(c)
	histories := appservice.ReadAllHistories()
	c.JSON(http.StatusOK, gin.H{"data": histories})
}

func FindHistory(c *gin.Context) {
	validateToken(c)
	id, _ := strconv.Atoi(c.Param("id"))
	history, err := appservice.ReadHistory(id)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": history})
}

func UpdateHistory(c *gin.Context) {
	checkUserLogged(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var history appmodel.History
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newHistory, err := appservice.UpdateHistory(id, history)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": newHistory})
}

func CreateHistory(c *gin.Context) {
	userID := checkUserLogged(c)
	var history appmodel.History
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newHistory, err := appservice.CreateHistory(history, userID)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": newHistory})
}

func FindAllGeneralInformations(c *gin.Context) {
	validateToken(c)
	gis := appservice.ReadAllGeneralInformations()
	c.JSON(http.StatusOK, gin.H{"data": gis})
}

func FindGeneralInformation(c *gin.Context) {
	validateToken(c)
	id, _ := strconv.Atoi(c.Param("id"))
	gi, err := appservice.ReadGeneralInformation(id)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gi})
}

func UpdateGeneralInformation(c *gin.Context) {
	checkUserLogged(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var gi appmodel.GeneralInformation
	if err := c.ShouldBindJSON(&gi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newGi, err := appservice.UpdateGeneralInformation(id, gi)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": newGi})
}

func CreateGeneralInformation(c *gin.Context) {
	userID := checkUserLogged(c)
	var gi appmodel.GeneralInformation
	if err := c.ShouldBindJSON(&gi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newGi, err := appservice.CreateGeneralInformation(gi, userID)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": newGi})
}

func FindAllUsers(c *gin.Context) {
	validateToken(c)
	users := appservice.ReaAlldUsers()
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func FindUser(c *gin.Context) {
	validateToken(c)
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := appservice.ReadUser(id)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdateUser(c *gin.Context) {
	checkUserLogged(c)
	id, _ := strconv.Atoi(c.Param("id"))

	var user appmodel.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := appservice.UpdateUser(id, user)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}

func CreateUser(c *gin.Context) {
	userID := checkUserLogged(c)

	var user appmodel.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newUser, err := appservice.CreateUser(user, userID)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": newUser})
}

func AuthUser(c *gin.Context) {
	var user appmodel.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := appservice.AuthUser(user)
	if handleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": session})

}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return true
	}
	return false
}

func getToken(tk []string) string {
	token := ""
	if len(tk) > 0 {
		token = tk[0]
	}
	return token
}

func checkUserLogged(c *gin.Context) int {
	token := getToken(c.Request.Header["Token"])
	err := appservice.ValidateSession(token)
	if err != nil {
		c.JSON(401, gin.H{"error": errors.New("Unauthorized")})
	}

	userID, _ := appservice.GetUserId(token)
	return userID
}

func validateToken(c *gin.Context) {
	token := getToken(c.Request.Header["Token"])
	appservice.ValidateSession(token)
}

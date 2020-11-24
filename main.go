package main

import (
	"fmt"

	"github.com/0sman/godemo/app/appmodel"
	"github.com/0sman/godemo/perm/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var db = initDB()

	service.InitService(db)
	service.InitPermissions(1, 1, 5)

	d := service.ReadSecuredModel(appmodel.History{})
	fmt.Println("output:", d)

	st := "new name"
	id := 1
	service.UpdateSecuredModel(id, appmodel.History{CourseName: &st})

	stt := "new insert name"
	service.CreateSecuredModel(appmodel.History{CourseName: &stt})
}

func initDB() *gorm.DB {
	dataSourceName := "root:osman123@tcp(localhost:3306)/godemo?parseTime=True"
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// db.Exec("CREATE DATABASE godemo")
	// db.Exec("USE godemo")

	/* db.AutoMigrate(
		&appmodel.GeneralInformation{},
		&appmodel.History{},
		&appmodel.User{})

	db.AutoMigrate(
		&dal.Object{},
		&dal.Owner{},
		&dal.Column{},
		&dal.Group{},
		&dal.Permission{})  */

	return db
}

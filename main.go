package main

import (
	"fmt"
	"reflect"

	"github.com/0sman/godemo/app/appmodel"
	"github.com/0sman/godemo/perm/service"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	var db = initDB()
	db.LogMode(true)
	service.InitService(db)
	service.InitPermissions(1, 1, 5)

	d := service.ReadSecuredModel(appmodel.History{})
	g := InterfaceToSlice(d)
	fmt.Println("output1:", g)

	st := "new name"
	id := 1
	service.UpdateSecuredModel(appmodel.History{HistoryID: &id, CourseName: &st})
}

func InterfaceToSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice).Elem()

	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func initDB() *gorm.DB {
	dataSourceName := "root:osman123@tcp(localhost:3306)/godemo?parseTime=True"
	var db, err = gorm.Open("mysql", dataSourceName)

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

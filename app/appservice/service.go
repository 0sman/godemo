package appservice

import (
	"time"

	"github.com/0sman/godemo/app/appmodel"
	"github.com/0sman/godemo/perm/service"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitService(dbRef *gorm.DB) {
	service.InitService(dbRef)
	service.InitPermissions(1, 1, 5)
}

func ReadHistories() []appmodel.History {
	var histMapList = service.ReadSecuredModel(appmodel.History{})
	var histories []appmodel.History
	for _, m := range histMapList {
		hist := historyFromMap(m)
		histories = append(histories, hist)
	}
	return histories
}

func UpdateHistory(id int, history appmodel.History) {
	service.UpdateSecuredModel(id, history)
}

func CreateHistory(history appmodel.History) {
	service.CreateSecuredModel(history)
}

func ReadGeneralInformations() []appmodel.GeneralInformation {
	var giMapList = service.ReadSecuredModel(appmodel.GeneralInformation{})
	var generalInformations []appmodel.GeneralInformation
	for _, m := range giMapList {
		gi := generalInformationFromMap(m)
		generalInformations = append(generalInformations, gi)
	}
	return generalInformations
}

func UpdateGeneralInformation(id int, gi appmodel.GeneralInformation) {
	service.UpdateSecuredModel(id, gi)
}

func CreateGeneralInformation(gi appmodel.GeneralInformation) {
	service.CreateSecuredModel(gi)
}

func ReadUsers() []appmodel.User {
	var userMapList = service.ReadSecuredModel(appmodel.User{})
	var users []appmodel.User
	for _, m := range userMapList {
		u := userFromMap(m)
		users = append(users, u)
	}
	return users
}

func UpdateUser(id int, user appmodel.User) {
	service.UpdateSecuredModel(id, user)
}

func CreateUser(user appmodel.User) {
	service.CreateSecuredModel(user)
}

func historyFromMap(m map[string]interface{}) appmodel.History {
	var ct *time.Time
	if m["course_time"] != nil {
		ct = m["course_time"].(*time.Time)
	}

	var cs *float64
	if m["course_score"] != nil {
		v := m["course_score"].(float64)
		cs = &v
	}

	hist := appmodel.History{
		HistoryID:   getIntRefValue(m, "history_id"),
		GiID:        getIntRefValue(m, "gi_id"),
		CourseTime:  ct,
		CourseName:  getStringRefValue(m, "course_name"),
		CourseScore: cs,
	}

	return hist
}

func generalInformationFromMap(m map[string]interface{}) appmodel.GeneralInformation {
	gi := appmodel.GeneralInformation{
		GiID:      getIntRefValue(m, "gi_id"),
		FirstName: getStringRefValue(m, "first_name"),
		LastName:  getStringRefValue(m, "last_name"),
		Email:     getStringRefValue(m, "email"),
		Phone:     getStringRefValue(m, "phone"),
		FinCode:   getStringRefValue(m, "fin_code"),
		Position:  getStringRefValue(m, "position"),
		Education: getStringRefValue(m, "education"),
	}

	return gi
}

func userFromMap(m map[string]interface{}) appmodel.User {
	user := appmodel.User{
		UserID:        getIntRefValue(m, "user_id"),
		Username:      getStringRefValue(m, "username"),
		Password:      getStringRefValue(m, "password"),
		Email:         getStringRefValue(m, "email"),
		AccessGroupId: getIntRefValue(m, "access_group_id"),
	}

	return user
}

func getStringRefValue(m map[string]interface{}, key string) *string {
	var st *string
	if m[key] != nil {
		v := m[key].(string)
		st = &v
	}
	return st
}

func getIntRefValue(m map[string]interface{}, key string) *int {
	var i *int
	if m[key] != nil {
		v := int(m[key].(int64))
		i = &v
	}
	return i
}

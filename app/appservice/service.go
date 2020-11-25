package appservice

import (
	"errors"
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

func ReadAllHistories() []appmodel.History {
	var histMapList = service.ReadAllSecuredModels(appmodel.History{})
	var histories []appmodel.History
	for _, m := range histMapList {
		hist, _ := historyFromMap(m)
		histories = append(histories, hist)
	}
	return histories
}

func ReadHistory(id int) (appmodel.History, error) {
	var histMap = service.ReadSecuredModel(id, appmodel.History{})
	var history, err = historyFromMap(histMap)
	return history, err
}

func UpdateHistory(id int, history appmodel.History) (appmodel.History, error) {
	newID, err := service.UpdateSecuredModel(id, history)
	if err == nil {
		return ReadHistory(newID.(int))
	}
	return appmodel.History{}, err
}

func CreateHistory(history appmodel.History) (appmodel.History, error) {
	newID, err := service.CreateSecuredModel(history)
	if err == nil {
		return ReadHistory(newID.(int))
	}
	return appmodel.History{}, err
}

func ReadAllGeneralInformations() []appmodel.GeneralInformation {
	var giMapList = service.ReadAllSecuredModels(appmodel.GeneralInformation{})
	var generalInformations []appmodel.GeneralInformation
	for _, m := range giMapList {
		gi, _ := generalInformationFromMap(m)
		generalInformations = append(generalInformations, gi)
	}
	return generalInformations
}

func ReadGeneralInformation(id int) (appmodel.GeneralInformation, error) {
	var giMap = service.ReadSecuredModel(id, appmodel.GeneralInformation{})
	var gi, err = generalInformationFromMap(giMap)
	return gi, err
}

func UpdateGeneralInformation(id int, gi appmodel.GeneralInformation) (appmodel.GeneralInformation, error) {
	newID, err := service.UpdateSecuredModel(id, gi)
	if err == nil {
		return ReadGeneralInformation(newID.(int))
	}
	return appmodel.GeneralInformation{}, err
}

func CreateGeneralInformation(gi appmodel.GeneralInformation) (appmodel.GeneralInformation, error) {
	newID, err := service.CreateSecuredModel(gi)
	if err == nil {
		return ReadGeneralInformation(newID.(int))
	}
	return appmodel.GeneralInformation{}, err
}

func ReaAlldUsers() []appmodel.User {
	var userMapList = service.ReadAllSecuredModels(appmodel.User{})
	var users []appmodel.User
	for _, m := range userMapList {
		u, _ := userFromMap(m)
		users = append(users, u)
	}
	return users
}

func ReadUser(id int) (appmodel.User, error) {
	var userMap = service.ReadSecuredModel(id, appmodel.User{})
	var user, err = userFromMap(userMap)
	return user, err
}

func UpdateUser(id int, user appmodel.User) (appmodel.User, error) {
	newID, err := service.UpdateSecuredModel(id, user)
	if err == nil {
		return ReadUser(newID.(int))
	}
	return appmodel.User{}, err
}

func CreateUser(user appmodel.User) (appmodel.User, error) {
	newID, err := service.CreateSecuredModel(user)
	if err == nil {
		return ReadUser(newID.(int))
	}
	return appmodel.User{}, err
}

func historyFromMap(m map[string]interface{}) (appmodel.History, error) {
	if isMapNull(m) {
		return appmodel.History{}, errors.New("Not found")
	}
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

	return hist, nil
}

func generalInformationFromMap(m map[string]interface{}) (appmodel.GeneralInformation, error) {
	if isMapNull(m) {
		return appmodel.GeneralInformation{}, errors.New("Not found")
	}

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

	return gi, nil
}

func userFromMap(m map[string]interface{}) (appmodel.User, error) {
	if isMapNull(m) {
		return appmodel.User{}, errors.New("Not found")
	}

	user := appmodel.User{
		UserID:        getIntRefValue(m, "user_id"),
		Username:      getStringRefValue(m, "username"),
		Password:      getStringRefValue(m, "password"),
		Email:         getStringRefValue(m, "email"),
		AccessGroupId: getIntRefValue(m, "access_group_id"),
	}

	return user, nil
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

func isMapNull(m map[string]interface{}) bool {
	var isNull = true
	for k := range m {
		if m[k] != nil {
			isNull = false
			break
		}
	}
	return isNull
}

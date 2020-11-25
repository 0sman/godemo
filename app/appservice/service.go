package appservice

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/0sman/godemo/app/appmodel"
	"github.com/0sman/godemo/perm/service"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitService(dbRef *gorm.DB) {
	db = dbRef
	service.InitService(dbRef)
	service.InitPermissions(1, 1, 5)
}

func ReadAllHistories() []appmodel.History {
	var histMapList, _ = service.ReadAllSecuredModels(appmodel.History{})
	var histories []appmodel.History
	for _, m := range histMapList {
		hist, _ := historyFromMap(m)
		histories = append(histories, hist)
	}
	return histories
}

func ReadHistory(id int) (appmodel.History, error) {
	var histMap, err = service.ReadSecuredModel(id, appmodel.History{})
	if err == nil {
		var history, er = historyFromMap(histMap)
		return history, er
	}
	return appmodel.History{}, err
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
		return ReadHistory(int(newID.(int64)))
	}
	return appmodel.History{}, err
}

func ReadAllGeneralInformations() []appmodel.GeneralInformation {
	var giMapList, _ = service.ReadAllSecuredModels(appmodel.GeneralInformation{})
	var generalInformations []appmodel.GeneralInformation
	for _, m := range giMapList {
		gi, _ := generalInformationFromMap(m)
		generalInformations = append(generalInformations, gi)
	}
	return generalInformations
}

func ReadGeneralInformation(id int) (appmodel.GeneralInformation, error) {
	var giMap, err = service.ReadSecuredModel(id, appmodel.GeneralInformation{})
	if err == nil {
		var gi, er = generalInformationFromMap(giMap)
		return gi, er
	}
	return appmodel.GeneralInformation{}, err
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
		return ReadGeneralInformation(int(newID.(int64)))
	}
	return appmodel.GeneralInformation{}, err
}

func ReaAlldUsers() []appmodel.User {
	var userMapList, _ = service.ReadAllSecuredModels(appmodel.User{})
	var users []appmodel.User
	for _, m := range userMapList {
		u, _ := userFromMap(m)
		users = append(users, u)
	}
	return users
}

func ReadUser(id int) (appmodel.User, error) {
	var userMap, err = service.ReadSecuredModel(id, appmodel.User{})
	if err == nil {
		var user, er = userFromMap(userMap)
		return user, er
	}
	return appmodel.User{}, err
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
		return ReadUser(int(newID.(int64)))
	}
	return appmodel.User{}, err
}

func AuthUser(um appmodel.User) (string, error) {
	user := appmodel.User{}
	db.Where("username = ? and password = ?", um.Username, um.Password).Take(&user)
	if user.UserID != nil && user.AccessGroupId != nil {
		uid := *user.UserID
		gid := *user.AccessGroupId
		//ids in plain text as 'session'
		session := strings.Join([]string{strconv.Itoa(uid), strconv.Itoa(gid)}, "-")
		return session, nil
	}
	return "", nil
}

func ValidateSession(session string) error {
	uid, grid, err := parseSession(session)
	owGrid := getOwnerGroupId()
	guestGrid := getGuestGroupId()
	if err != nil {
		grid = guestGrid //login as guest
	}

	service.InitPermissions(uid, grid, owGrid)
	return err
}

func parseSession(session string) (int, int, error) {
	ids := strings.Split(session, "-")
	if len(ids) == 2 {
		userID, err1 := strconv.Atoi(ids[0])
		groupID, err2 := strconv.Atoi(ids[1])
		if err1 == nil && err2 == nil {
			return userID, groupID, nil
		}
		return -1, -1, errors.New("invalid session")
	}
	return -1, -1, errors.New("invalid session")
}

func getOwnerGroupId() int {
	return 4
}

func getGuestGroupId() int {
	return 5
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

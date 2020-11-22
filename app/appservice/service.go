package appservice

import (
	"github.com/0sman/godemo/app/appmodel"
	"github.com/0sman/godemo/perm/service"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func InitService(dbRef *gorm.DB) {
	db = dbRef
}

func ReadHistory() []appmodel.History {
	var histories []appmodel.History
	if service.CanRead(histories, 0) {
		db.Find(&histories)
	}

	return histories
}

package appmodel

import "time"

type GeneralInformation struct {
	GiID      *uint `gorm:"primary_key"`
	FirstName *string
	LastName  *string
	Email     *string
	Phone     *string
	FinCode   *string
	Position  *string
	Education *string
}

type History struct {
	HistoryID   *uint              `gorm:"primary_key" perm:"history_id"`
	GiID        *uint              `perm:"gi_id"`
	Gi          GeneralInformation `gorm:"foreignKey:GiID"`
	CourseTime  *time.Time         `perm:"course_time"`
	CourseName  *string            `perm:"course_name"`
	CourseScore *float32           `perm:"course_score"`
}

type User struct {
	UserID        *uint `gorm:"primary_key"`
	Username      *string
	Password      *string
	Email         *string
	Gi            GeneralInformation `gorm:"foreignKey:Email"`
	AccessGroupId *uint
}

func (history History) GetTableName() string {
	return "histories"
}

func (history History) GetDal() interface{} {
	return history
}

package appmodel

import (
	"time"
)

type GeneralInformation struct {
	GiID      *int    `gorm:"primary_key" perm:"gi_id"`
	FirstName *string `perm:"first_name"`
	LastName  *string `perm:"last_name"`
	Email     *string `perm:"email"`
	Phone     *string `perm:"phone"`
	FinCode   *string `perm:"fin_code"`
	Position  *string `perm:"position"`
	Education *string `perm:"education"`
}

type History struct {
	HistoryID   *int               `gorm:"primary_key" perm:"history_id"`
	GiID        *int               `perm:"gi_id"`
	Gi          GeneralInformation `gorm:"foreignKey:GiID"`
	CourseTime  *time.Time         `perm:"course_time"`
	CourseName  *string            `perm:"course_name"`
	CourseScore *float64           `perm:"course_score"`
}

type User struct {
	UserID        *int               `gorm:"primary_key" perm:"user_id"`
	Username      *string            `perm:"username"`
	Password      *string            `perm:"password"`
	Email         *string            `perm:"email"`
	Gi            GeneralInformation `gorm:"foreignKey:Email"`
	AccessGroupId *int               `perm:"access_group_id"`
}

func (gi GeneralInformation) GetTableName() string {
	return "general_informations"
}

func (gi GeneralInformation) GetDal() interface{} {
	return gi
}

func (history History) GetTableName() string {
	return "histories"
}

func (history History) GetDal() interface{} {
	return history
}

func (user User) GetTableName() string {
	return "users"
}

func (user User) GetDal() interface{} {
	return user
}

package appmodel

import (
	"time"
)

type GeneralInformation struct {
	GiID      *int    `json:"gi_id" gorm:"primary_key" perm:"gi_id"`
	FirstName *string `json:"first_name" perm:"first_name"`
	LastName  *string `json:"last_name" perm:"last_name"`
	Email     *string `json:"email" perm:"email"`
	Phone     *string `json:"phone" perm:"phone"`
	FinCode   *string `json:"fin_code" perm:"fin_code"`
	Position  *string `json:"position" perm:"position"`
	Education *string `json:"education" perm:"education"`
}

type History struct {
	HistoryID   *int               `json:"history_id" gorm:"primary_key" perm:"history_id"`
	GiID        *int               `json:"gi_id" perm:"gi_id"`
	Gi          GeneralInformation `json:"-" gorm:"foreignKey:GiID"`
	CourseTime  *time.Time         `json:"course_time" perm:"course_time"`
	CourseName  *string            `json:"course_name" perm:"course_name"`
	CourseScore *float64           `json:"course_score" perm:"course_score"`
}

type User struct {
	UserID        *int               `json:"user_id" gorm:"primary_key" perm:"user_id"`
	Username      *string            `json:"username" perm:"username"`
	Password      *string            `json:"password" perm:"password"`
	Email         *string            `json:"email" perm:"email"`
	Gi            GeneralInformation `json:"-" gorm:"foreignKey:Email"`
	AccessGroupId *int               `json:"access_group_id" perm:"access_group_id"`
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

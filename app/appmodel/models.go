package appmodel

import (
	"time"
)

type GeneralInformation struct {
	GiID      *int    `json:"gi_id,omitempty" gorm:"primary_key" perm:"gi_id"`
	FirstName *string `json:"first_name,omitempty" perm:"first_name"`
	LastName  *string `json:"last_name,omitempty" perm:"last_name"`
	Email     *string `json:"email,omitempty" perm:"email"`
	Phone     *string `json:"phone,omitempty" perm:"phone"`
	FinCode   *string `json:"fin_code,omitempty" perm:"fin_code"`
	Position  *string `json:"position,omitempty" perm:"position"`
	Education *string `json:"education,omitempty" perm:"education"`
}

type History struct {
	HistoryID   *int               `json:"history_id,omitempty" gorm:"primary_key" perm:"history_id"`
	GiID        *int               `json:"gi_id,omitempty" perm:"gi_id"`
	Gi          GeneralInformation `json:"-" gorm:"foreignKey:GiID"`
	CourseTime  *time.Time         `json:"course_time,omitempty" perm:"course_time"`
	CourseName  *string            `json:"course_name,omitempty" perm:"course_name"`
	CourseScore *float64           `json:"course_score,omitempty" perm:"course_score"`
}

type User struct {
	UserID        *int               `json:"user_id,omitempty" gorm:"primary_key" perm:"user_id"`
	Username      *string            `json:"username,omitempty" perm:"username"`
	Password      *string            `json:"password,omitempty" perm:"password"`
	Email         *string            `json:"email,omitempty" perm:"email"`
	Gi            GeneralInformation `json:"-" gorm:"foreignKey:Email"`
	AccessGroupId *int               `json:"access_group_id,omitempty" perm:"access_group_id"`
}

func (gi GeneralInformation) GetTableName() string {
	return "general_informations"
}

func (history History) GetTableName() string {
	return "histories"
}

func (user User) GetTableName() string {
	return "users"
}

func (gi GeneralInformation) GetDal() interface{} {
	return gi
}

func (history History) GetDal() interface{} {
	return history
}

func (user User) GetDal() interface{} {
	return user
}

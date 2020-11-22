package service

import (
	"fmt"
	"reflect"

	"github.com/0sman/godemo/perm/dal"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

var permMap = make(map[string]PermContainer)

type PermContainer struct {
	// map of column name -> permission mask
	SecuredColumns      map[string]int
	SecuredOwnedColumns map[string]int
	OwnedRowIds         []int
}

func InitService(dbRef *gorm.DB) {
	db = dbRef
}

const (
	PermissionRead   int = 1
	PermissionWrite  int = 2
	PermissionCreate int = 4
)

func InitPermissions(userId int, groupId int, ownerGrId int) {

	var objects []dal.Object
	db.Find(&objects)

	for _, obj := range objects {

		var columnMap = make(map[int]string)
		var columns []dal.Column
		db.Where(&dal.Column{ObjectId: obj.ObjectId}).Find(&columns)

		var columnsIds []int
		for _, cl := range columns {
			columnsIds = append(columnsIds, cl.ColumnId)
			columnMap[cl.ColumnId] = cl.Name
		}

		var permissions []dal.Permission
		db.Where("group_id = ? and column_id IN (?)", groupId, columnsIds).Find(&permissions)

		var ownedPermissions []dal.Permission
		db.Where("group_id = ? and column_id IN (?)", ownerGrId, columnsIds).Find(&ownedPermissions)

		var ownedRows []dal.Owner
		db.Where(&dal.Owner{UserId: userId, ObjectId: obj.ObjectId}).Find(&ownedRows)

		var cmMap = make(map[string]int)
		for _, perm := range permissions {
			cmMap[columnMap[perm.ColumnId]] = perm.PermMask
		}

		var opMap = make(map[string]int)
		for _, perm := range ownedPermissions {
			opMap[columnMap[perm.ColumnId]] = perm.PermMask
		}

		var ownedRowIds []int
		for _, ow := range ownedRows {
			ownedRowIds = append(ownedRowIds, ow.TargetRow)
		}

		var pc PermContainer
		pc.OwnedRowIds = ownedRowIds
		pc.SecuredColumns = cmMap
		pc.SecuredOwnedColumns = opMap

		permMap[obj.Name] = pc

	}
}

type SecuredModel interface {
	GetTableName() string
	GetAllColumns() []string
	GetDal() interface{}
}

func ReadSecuredModel(secModel SecuredModel) interface{} {
	var ac = GetAllowedReadColumns(secModel, 1)
	var tp = reflect.TypeOf(secModel.GetDal())
	var results = reflect.New(reflect.SliceOf(tp)).Interface()
	db.Select(ac).Table(secModel.GetTableName()).Find(results)
	return results
}

func UpdateSecuredModel(secModel SecuredModel) {
	//var ac = GetAllowedWriteColumns(secModel, 1)

	var tp = reflect.TypeOf(secModel.GetDal())
	var rv = reflect.ValueOf(secModel.GetDal())
	for i := 0; i < tp.NumField(); i++ {
		if value, ok := tp.Field(i).Tag.Lookup("perm"); ok {
			if ok {

				fmt.Println("field:", rv.Field(i))
				fmt.Println("value:", value)
			}

		}
	}
}

func GetAllowedReadColumns(secModel SecuredModel, row int) []string {
	return getAllowedColumns(PermissionRead, secModel, row)
}

func GetAllowedWriteColumns(secModel SecuredModel, row int) []string {
	return getAllowedColumns(PermissionWrite, secModel, row)
}

func GetAllowedCreateColumns(secModel SecuredModel, row int) []string {
	return getAllowedColumns(PermissionCreate, secModel, row)
}

func getAllowedColumns(permission int, secModel SecuredModel, row int) []string {
	var result []string
	for _, cl := range secModel.GetAllColumns() {
		if checkPermission(permission, secModel.GetTableName(), cl, row) == true {
			result = append(result, cl)
		}
	}
	return result
}

func CanRead(secModel SecuredModel, row int) bool {
	return checkSecuredModel(PermissionRead, secModel, row)
}

func CanWrite(secModel SecuredModel, row int) bool {
	return checkSecuredModel(PermissionWrite, secModel, row)
}

func CanCreate(secModel SecuredModel, row int) bool {
	return checkSecuredModel(PermissionCreate, secModel, row)
}

func checkSecuredModel(permission int, secModel SecuredModel, row int) bool {
	for _, clm := range secModel.GetAllColumns() {
		if !checkPermission(PermissionRead, secModel.GetTableName(), clm, row) {
			return false
		}
	}
	return true
}

func checkPermission(permission int, table string, column string, row int) bool {
	pc := permMap[table]

	var totalPermission int

	for _, v := range pc.OwnedRowIds {
		if v == row { // owner
			totalPermission |= pc.SecuredOwnedColumns[column]
		}
	}

	totalPermission |= pc.SecuredColumns[column]
	return (totalPermission&permission > 0)
}

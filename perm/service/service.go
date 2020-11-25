package service

import (
	"fmt"
	"reflect"

	"github.com/0sman/godemo/perm/dal"
	"gorm.io/gorm"
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
	GetDal() interface{}
}

func ReadSecuredModel(secModel SecuredModel) []map[string]interface{} {
	var ac = GetAllowedReadColumns(secModel, 1)
	var tp = reflect.TypeOf(secModel.GetDal())
	var results = reflect.New(reflect.SliceOf(tp)).Interface()
	db.Select(ac).Table(secModel.GetTableName()).Find(results)

	var sl = interfaceToSlice(results)
	var res []map[string]interface{}

	for _, s := range sl {
		var mp = make(map[string]interface{})
		var rv = reflect.ValueOf(s)
		for i := 0; i < tp.NumField(); i++ {
			fl := rv.Field(i)
			if tagName, ok := tp.Field(i).Tag.Lookup("perm"); ok {
				if fl.Kind() == reflect.Ptr {
					if !fl.IsNil() {
						var propValue = GetPropValue(fl)
						mp[tagName] = propValue
					}
				}
			}
		}
		res = append(res, mp)
	}

	return res
}

func UpdateSecuredModel(id interface{}, secModel SecuredModel) {
	modelMap, ok := buildModelMap(secModel)
	if ok {
		pkName := getPKColumnName(secModel)
		pkMap := make(map[string]interface{})
		pkMap[pkName] = id
		db.Table(secModel.GetTableName()).Where(pkMap).Updates(modelMap)
	} else {
		fmt.Println("Update is not allowed. Please, check your permissions.")
	}
}

func CreateSecuredModel(secModel SecuredModel) {
	modelMap, ok := buildModelMap(secModel)
	if ok {
		db.Table(secModel.GetTableName()).Create(modelMap)
	} else {
		fmt.Println("Create is not allowed. Please, check your permissions.")
	}
}

func buildModelMap(secModel SecuredModel) (modelMap map[string]interface{}, ok bool) {
	var allowedColumns = GetAllowedWriteColumns(secModel, 1)

	var tp = reflect.TypeOf(secModel.GetDal())
	var rv = reflect.ValueOf(secModel.GetDal())

	modelMap = make(map[string]interface{})

	for i := 0; i < tp.NumField(); i++ {
		if clName, ok := tp.Field(i).Tag.Lookup("perm"); ok {
			fl := rv.Field(i)
			if !fl.IsNil() {
				if !isColumnInList(allowedColumns, clName) {
					return nil, false
				}
				var propValue = GetPropValue(fl)
				modelMap[clName] = propValue
			}
		}
	}

	return modelMap, true
}

func isColumnInList(slice []string, target string) bool {
	for _, ac := range slice {
		if ac == target {
			return true
		}
	}
	return false
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
	for _, cl := range getAllColumns(secModel) {
		if checkPermission(permission, secModel.GetTableName(), cl, row) == true {
			result = append(result, cl)
		}
	}
	return result
}

func getAllColumns(secModel SecuredModel) []string {
	var res []string
	var tp = reflect.TypeOf(secModel.GetDal())
	for i := 0; i < tp.NumField(); i++ {
		if value, ok := tp.Field(i).Tag.Lookup("perm"); ok {
			res = append(res, value)
		}
	}
	return res
}

func getPKColumnName(secModel SecuredModel) string {
	var tp = reflect.TypeOf(secModel.GetDal())
	for i := 0; i < tp.NumField(); i++ {
		tag := tp.Field(i).Tag
		if value, ok := tag.Lookup("gorm"); ok {
			if value == "primary_key" {
				if value, ok := tag.Lookup("perm"); ok {
					return value
				}
			}
		}
	}
	return ""
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

func interfaceToSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice).Elem()

	if s.Kind() != reflect.Slice {
		panic("interfaceToSlice() given a non-slice type")
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

func GetPropValue(propValue reflect.Value) interface{} {
	var val = propValue.Elem()
	var res interface{}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		res = val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		res = val.Uint()
	case reflect.String:
		res = val.String()
	case reflect.Bool:
		res = val.Bool()
	case reflect.Float32, reflect.Float64:
		res = val.Float()
	default:
		res = propValue.Interface()
	}
	return res
}

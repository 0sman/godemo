package service

import (
	"errors"
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
	PermissionUpdate int = 2
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

func ReadAllSecuredModels(secModel SecuredModel) ([]map[string]interface{}, error) {
	var ac = getAllowedColumns(PermissionRead, secModel, false)
	if len(ac) == 0 {
		return nil, errors.New("Read is not allowed")
	}

	pkName := getPKColumnName(secModel)
	var userOwnRows []int

	idAdded := false
	if !isColumnInList(ac, pkName) {
		idAdded = true
		ac = append(ac, pkName)
	}

	var tp = reflect.TypeOf(secModel)
	var results = reflect.New(reflect.SliceOf(tp)).Interface()
	db.Select(ac).Table(secModel.GetTableName()).Find(results)

	var sl = interfaceToSlice(results)
	var res []map[string]interface{}

	for _, s := range sl {
		mp := toModelMap(s, tp)
		id := int(mp[pkName].(int64))

		if checkRowForOwner(secModel.GetTableName(), id) {
			userOwnRows = append(userOwnRows, id)
		} else {
			if idAdded {
				mp[pkName] = nil
			}
			res = append(res, mp)
		}
	}

	if len(userOwnRows) > 0 {
		totalCl, err := getTotalAllowedColumns(PermissionRead, secModel)
		if err != nil {
			return nil, err
		}
		results = reflect.New(reflect.SliceOf(tp)).Interface()
		db.Select(totalCl).Table(secModel.GetTableName()).Where(userOwnRows).Find(results)
		sl = interfaceToSlice(results)

		for _, s := range sl {
			mp := toModelMap(s, tp)
			res = append(res, mp)
		}
	}

	return res, nil
}

func ReadSecuredModel(id interface{}, secModel SecuredModel) (map[string]interface{}, error) {
	var ac, err = getTotalAllowedColumnsForRow(PermissionRead, secModel, id)
	if err != nil {
		return nil, err
	}

	var tp = reflect.TypeOf(secModel)
	var result = reflect.New(tp).Interface()

	pkName := getPKColumnName(secModel)
	pkMap := make(map[string]interface{})
	pkMap[pkName] = id

	db.Select(ac).Table(secModel.GetTableName()).Where(pkMap).Take(result)

	res := toModelMap(result, tp)
	return res, nil
}

func UpdateSecuredModel(id interface{}, secModel SecuredModel) (interface{}, error) {
	var ac, err = getTotalAllowedColumnsForRow(PermissionUpdate, secModel, id)
	if err != nil {
		return nil, err
	}

	modelMap, ok := buildModelMap(ac, secModel)
	if ok {
		pkName := getPKColumnName(secModel)
		pkMap := make(map[string]interface{})
		pkMap[pkName] = id
		db.Table(secModel.GetTableName()).Where(pkMap).Updates(modelMap)
		return pkMap[pkName], nil
	}
	return nil, errors.New("Update is not allowed")
}

func CreateSecuredModel(secModel SecuredModel, userID int) (interface{}, error) {
	var ac, err = getTotalAllowedColumns(PermissionCreate, secModel)
	if err != nil {
		return nil, err
	}

	var tp = reflect.TypeOf(secModel)
	var result = reflect.New(tp).Interface()

	db.Table(secModel.GetTableName()).Create(result)
	keyValue := getPKColumnValue(result.(SecuredModel))

	modelMap, ok := buildModelMap(ac, secModel)
	if ok {
		pkName := getPKColumnName(secModel)
		pkMap := make(map[string]interface{})
		pkMap[pkName] = keyValue
		db.Table(secModel.GetTableName()).Where(pkMap).Updates(modelMap)
		id := int(pkMap[pkName].(int64))
		addOwnerForRow(secModel.GetTableName(), id, userID)
		return pkMap[pkName], nil
	}
	return nil, errors.New("Create is not allowed")
}

func addOwnerForRow(table string, row int, userID int) {
	var obj dal.Object
	db.Where(&dal.Object{Name: table}).First(&obj)

	owner := dal.Owner{
		ObjectId:  obj.ObjectId,
		TargetRow: row,
		UserId:    userID,
	}
	db.Create(&owner)
}

func buildModelMap(allowedColumns []string, secModel SecuredModel) (modelMap map[string]interface{}, ok bool) {
	var tp = reflect.TypeOf(secModel)
	var rv = reflect.ValueOf(secModel)

	modelMap = make(map[string]interface{})

	for i := 0; i < tp.NumField(); i++ {
		if clName, ok := tp.Field(i).Tag.Lookup("perm"); ok {
			fl := rv.Field(i)
			if !fl.IsNil() {
				if !isColumnInList(allowedColumns, clName) {
					return nil, false
				}
				var propValue = getPropValue(fl)
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

func getTotalAllowedColumnsForRow(permission int, secModel SecuredModel, id interface{}) ([]string, error) {
	var ac = getAllowedColumns(permission, secModel, false)
	if len(ac) == 0 {
		return nil, fmt.Errorf("Permission %d is not allowed", permission)
	}

	//current implementation supports only int IDs
	if checkRowForOwner(secModel.GetTableName(), id.(int)) {
		var aoc = getAllowedColumns(permission, secModel, true)
		for _, cl := range aoc {
			if !isColumnInList(ac, cl) {
				ac = append(ac, cl)
			}
		}
	}

	return ac, nil
}

func getTotalAllowedColumns(permission int, secModel SecuredModel) ([]string, error) {
	var ac = getAllowedColumns(permission, secModel, false)
	if len(ac) == 0 {
		return nil, fmt.Errorf("Permission %d is not allowed", permission)
	}
	var aoc = getAllowedColumns(permission, secModel, true)
	for _, cl := range aoc {
		if !isColumnInList(ac, cl) {
			ac = append(ac, cl)
		}
	}

	return ac, nil
}

func getAllowedColumns(permission int, secModel SecuredModel, isOwner bool) []string {
	var result []string
	for _, cl := range getAllColumns(secModel) {
		if checkPermission(permission, secModel.GetTableName(), cl, isOwner) == true {
			result = append(result, cl)
		}
	}
	return result
}

func getAllColumns(secModel SecuredModel) []string {
	var res []string
	var tp = reflect.TypeOf(secModel)
	for i := 0; i < tp.NumField(); i++ {
		if value, ok := tp.Field(i).Tag.Lookup("perm"); ok {
			res = append(res, value)
		}
	}
	return res
}

func getPKColumnValue(secModel SecuredModel) interface{} {
	var tp = reflect.TypeOf(secModel.GetDal())
	var rv = reflect.ValueOf(secModel.GetDal())
	for i := 0; i < tp.NumField(); i++ {
		if value, ok := tp.Field(i).Tag.Lookup("gorm"); ok {
			if value == "primary_key" {
				if _, ok := tp.Field(i).Tag.Lookup("perm"); ok {
					var propValue = getPropValue(rv.Field(i))
					return propValue
				}
			}
		}
	}
	return nil
}

func getPKColumnName(secModel SecuredModel) string {
	var tp = reflect.TypeOf(secModel)
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

func checkPermission(permission int, table string, column string, isOwner bool) bool {
	pc := permMap[table]
	if isOwner {
		return (pc.SecuredOwnedColumns[column]&permission > 0)
	}
	return (pc.SecuredColumns[column]&permission > 0)
}

func checkRowForOwner(table string, row int) bool {
	pc := permMap[table]
	for _, v := range pc.OwnedRowIds {
		if v == row {
			return true
		}
	}
	return false
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

func toModelMap(sm interface{}, tp reflect.Type) map[string]interface{} {
	var mp = make(map[string]interface{})
	var rv = reflect.ValueOf(sm)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		fl := rv.Field(i)
		if tagName, ok := tp.Field(i).Tag.Lookup("perm"); ok {
			if fl.Kind() == reflect.Ptr {
				if !fl.IsNil() {
					var propValue = getPropValue(fl)
					mp[tagName] = propValue
				}
			}
		}
	}

	return mp
}

func getPropValue(propValue reflect.Value) interface{} {
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

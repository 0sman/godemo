package dal

type Object struct {
	ObjectId int `gorm:"primary_key"`
	Name     string
	IdColumn string
}

type Column struct {
	ColumnId int `gorm:"primary_key"`
	ObjectId int
	Object   Object `gorm:"foreignKey:ObjectId"`
	Name     string
}

type Owner struct {
	OwnerId   int `gorm:"primary_key"`
	ObjectId  int
	Object    Object `gorm:"foreignKey:ObjectId"`
	TargetRow int
	UserId    int
}

type Group struct {
	GroupId int `gorm:"primary_key"`
	Name    string
}

type Permission struct {
	PermId   int `gorm:"primary_key"`
	ColumnId int
	Column   Column `gorm:"foreignKey:ColumnId"`
	GroupId  int
	Group    Group `gorm:"foreignKey:GroupId"`
	PermMask int
}

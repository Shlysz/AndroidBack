package respo

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitSQL() *gorm.DB {
	//初始化数据库
	username := "root"
	password := "123456"
	host := "127.0.0.1"
	port := "3306"
	database := "androidwork"
	db, err := gorm.Open(
		mysql.Open(username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	return db
}

var GolbalDB = InitSQL()

package respo

import (
	"crypto/md5"
	"encoding/base64"
)

type Account struct {
	Username string `gorm:"primarykey"`
	Password string
	Key      string
}

func (a Account) TableName() string {
	return "account"
}

// Login account的子函数
func (a *Account) Login() bool {
	account := Account{}
	//查找数据库中时候有没有这个账号
	GolbalDB.First(&account).Where("username = ?", a.Username)
	//如果没有这个账号
	if account.Username == "" {
		return false
	}
	//如果有这个账号，但是密码不对
	if account.Password != a.Password {
		return false
	}
	//如果没有key,则分配一个key
	if account.Key == "" {
		//使用md5生成一个key
		key := md5.Sum([]byte(a.Username))
		//使用base64编码
		encodedString := base64.StdEncoding.EncodeToString(key[:])
		a.Key = encodedString
		GolbalDB.Model(&account).Update("key", encodedString)
		return true
	}
	a.Key = account.Key
	return true
}

func (a Account) GetUsernameByKey() string {
	account := Account{}
	tx := GolbalDB.First(&account).Where("key = ?", a.Key)
	if tx.Error != nil {
		return ""
	}
	return account.Username
}

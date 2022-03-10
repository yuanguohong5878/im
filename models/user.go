package models

import (
	"errors"
	"github.com/github-user/IMv_3/helper"
)

type User struct {
	ID       int    `gorm:"primary_key" json:"id"` // 用户id
	UserName string `json:"userName"`              // 账号,用账号做sessionid
	Password string `json:"password"`              // 密码
	Name     string `json:"name"`                  // 名字
}

type RegisterReq struct {
	UserName string `json:"userName"` // 账号
	Password string `json:"password"` // 密码
	Name     string `json:"name"`     // 名字
}

// RegisterReq POST /login login请求
type LoginReq struct {
	UserName string `json:"userName"` // 账号
	Password string `json:"password"` // 密码
}

func Register(name, username, password string) error {
	var users []User
	db.Where("user_name = ?", username).Find(&users)
	if len(users) != 0 {
		return errors.New("email has exits")
	}
	var data []byte = []byte(password)
	password = helper.GetSHA256HashCode(data)
	var user = User{Name: name, UserName: username, Password: password}
	db.Create(&user)
	return nil
}

func IsExistUsername(username string) bool {
	var users []User
	db.Where("user_name = ?", username).Find(&users)
	return len(users) > 0
}

func FindUserByUsername(username string) User {
	var user User
	db.Where("user_name = ?", username).First(&user)
	return user
}

func Login(username, password string) error {
	var users []User
	db.Where("user_name = ?", username).Find(&users)
	if len(users) == 0 {
		return errors.New("user not exits")
	}
	var data []byte = []byte(password)
	hashCode := helper.GetSHA256HashCode(data)
	if hashCode != users[0].Password {
		return errors.New("password error")
	}
	return nil
}

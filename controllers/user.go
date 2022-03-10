package controllers

import (
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/github-user/IMv_3/models"
)

type Response struct {
	Code    int         `bson:"code" json:"code"`
	Type    string      `bson:"type" json:"type"`
	Message interface{} `bson:"message" json:"message"`
}

type UserInfoRes struct {
	Code  int
	State string
	Data  models.User
}

const (
	StatusSuccess  = "success"
	StatusBadReq   = "bad_req"
	StatusNotLogin = "not_login"
	StatusNotAllow = "not_allow"
	StatusExist    = "had_exist"
	StatusError    = "error"
	StatusNotValid = "not_invalid"
)

func Register(c *gin.Context) {

	var user models.RegisterReq
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil || user.UserName == "" || user.Password == "" {
		c.JSON(http.StatusOK, &Response{
			Code:    400,
			Type:    "fail",
			Message: "username and password do not match",
		})
		return
	}
	err = models.Register(user.Name, user.UserName, user.Password)
	if err != nil {
		c.JSON(http.StatusOK, &Response{
			Code:    400,
			Type:    "fail",
			Message: "username and password do not match",
		})
		return
	}
	c.JSON(http.StatusOK, &Response{
		Code: 200,
		Type: "success",
	})
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	option := sessions.Options{MaxAge: 3600 * 24, Path: "/"}
	session.Options(option)
	var user models.LoginReq
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, &Response{
			Code:    400,
			Type:    "fail",
			Message: "bad_eq",
		})
	}
	exist := models.IsExistUsername(user.UserName)
	if exist {
		err = models.Login(user.UserName, user.Password)
		if err != nil {
			c.JSON(http.StatusOK, &Response{
				Code:    400,
				Type:    "fail",
				Message: "username and password do not match",
			})
		} else {
			var u models.User
			u = models.FindUserByUsername(user.UserName)
			fmt.Println(u.UserName)
			session.Set("sessionid", u.UserName)
			session.Save()
			///////////////

			///////////////
			c.JSON(http.StatusOK, &Response{
				Code:    http.StatusOK,
				Type:    "success",
				Message: u.UserName,
			})
		}

	} else {
		c.JSON(http.StatusOK, &Response{
			Code:    400,
			Type:    "fail",
			Message: "the user not exist",
		})
	}
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("sessionid")
		fmt.Println("cookie: ", v)
		if v != nil {
			// 验证通过，会继续访问下一个中间件
			c.Next()
		} else {
			// 验证不通过，不再调用后续的函数处理
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
			return
		}
	}
}

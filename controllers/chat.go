package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/github-user/IMv_3/models"
	"github.com/github-user/IMv_3/service"
	"github.com/gorilla/websocket"
	"net/http"
	"sort"
	"strconv"
)

const (
	PCStatus     = "PC"
	MOBILEStatus = "MOBILE"
	WEBStatus    = "WEB"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectWeb(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("sessionid")
	connectType := c.GetHeader("User-Agent")
	switch connectType {
	case PCStatus:
		connectType = PCStatus
	case MOBILEStatus:
		connectType = MOBILEStatus
	default:
		connectType = WEBStatus
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		http.NotFound(c.Writer, c.Request)
		return
	}
	sessionid := v.(string)
	service.ConnectWeb(conn, sessionid, connectType)
}

func FindMessages(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("sessionid")
	userName := v.(string)
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	startTime, _ := strconv.ParseUint(startDate, 10, 32)
	endTime, _ := strconv.ParseUint(endDate, 10, 32)
	sendMessages := models.FindMessagesBySendid(userName, startTime, endTime)
	recMessages := models.FindMessagesByRecid(userName, startTime, endTime)
	messages := append(sendMessages, recMessages...)
	sort.Sort(models.MessageSort(messages))
	c.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Msg:  StatusSuccess,
		Data: messages,
	})
}

func FindChats(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("sessionid")
	userName := v.(string)
	//查询置顶的会话列表
	topChats := models.FindChatsBySendidTop(userName)
	//查询非置顶的会话列表
	notTopChats := models.FindChatsBySendidNotTop(userName)
	chats := append(topChats, notTopChats...)
	c.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Msg:  StatusSuccess,
		Data: chats,
	})
}

func UpdateChat(c *gin.Context) {
	var chat models.Chat
	json.NewDecoder(c.Request.Body).Decode(&chat)
	state := chat.State
	id := c.Param("id")
	if err := models.UpdateChatState(id, state); err != nil {
		c.JSON(http.StatusInternalServerError, &Response{
			Code: http.StatusInternalServerError,
			Msg:  StatusError,
		})
	} else {
		c.JSON(http.StatusOK, &Response{
			Code: http.StatusOK,
			Msg:  StatusSuccess,
		})
	}
}

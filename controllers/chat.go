package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/github-user/IMv_3/kafka"
	"github.com/github-user/IMv_3/models"
	"github.com/github-user/IMv_3/service"
	"github.com/gorilla/websocket"
	"net/http"
	"sort"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMessage(c *gin.Context) {
	var receiveMsg models.Message
	err := json.NewDecoder(c.Request.Body).Decode(&receiveMsg)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, &Response{
			Code:    400,
			Type:    "fail",
			Message: "bad_eq",
		})
	}
	message, _ := json.Marshal(receiveMsg)
	kafka.SendMessage(message)
	c.JSON(http.StatusOK, &Response{
		Code:    http.StatusOK,
		Type:    "success",
		Message: "",
	})
}

func ConnectWeb(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("sessionid")
	connectType := c.Query("type")
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
	userName := c.Query("userName")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	startTime, _ := strconv.ParseUint(startDate, 10, 32)
	endTime, _ := strconv.ParseUint(endDate, 10, 32)
	sendMessages := models.FindMessagesBySendid(userName, startTime, endTime)
	recMessages := models.FindMessagesByRecid(userName, startTime, endTime)
	messages := append(sendMessages, recMessages...)
	sort.Sort(models.MessageSort(messages))
	c.JSON(http.StatusOK, &Response{
		Code:    http.StatusOK,
		Type:    "success",
		Message: messages,
	})
}

func FindChats(c *gin.Context) {
	userName := c.Query("userName")
	//查询置顶的会话列表
	topChats := models.FindChatsBySendidTop(userName)
	//查询非置顶的会话列表
	notTopChats := models.FindChatsBySendidNotTop(userName)
	chats := append(topChats, notTopChats...)
	c.JSON(http.StatusOK, &Response{
		Code:    http.StatusOK,
		Type:    "success",
		Message: chats,
	})
}

func UpdateChat(c *gin.Context) {
	var chat models.Chat
	json.NewDecoder(c.Request.Body).Decode(&chat)
	state := chat.State
	id := c.Param("id")
	models.UpdateChatState(id, state)
	c.JSON(http.StatusOK, &Response{
		Code: http.StatusOK,
		Type: "success",
	})
}

package service

import (
	"encoding/json"
	"fmt"
	"github.com/github-user/IMv_3/models"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"sync"
)

type ClientManger struct {
	clientsLock  sync.RWMutex
	clients      map[*Client]string
	singleClient map[string]map[string]*Client
}

type Client struct {
	socket  *websocket.Conn
	message ReceiveMessage
	//mutex   sync.Mutex
}

type ReceiveMessage struct {
	Uuid string `json:"uuid"`
	Type string `json:"type"`
}

var manager = ClientManger{
	clients:      make(map[*Client]string, 1000),
	singleClient: make(map[string]map[string]*Client, 1000),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnectWeb(conn *websocket.Conn, sessionid string, connectType string) {
	var receiveMsg ReceiveMessage
	receiveMsg.Uuid = sessionid
	receiveMsg.Type = connectType
	client := &Client{socket: conn}
	client.message = receiveMsg
	manager.register(client, receiveMsg)
	//查询未读消息
	messages := models.FindUnReadMessageByRecid(sessionid)
	//把未读消息推送
	for _, message := range messages {
		if toClients, ok := manager.singleClient[sessionid]; ok {
			msg, _ := json.Marshal(message)
			for _, toClient := range toClients {
				toClient.write(string(msg))
			}
		}
	}
	models.UpdateUnReadMessageByRecid(sessionid)

}

func HandleMessage(receiveMsg models.Message) {
	if toClients, ok := manager.singleClient[receiveMsg.Recid]; ok {
		message, _ := json.Marshal(receiveMsg)
		for _, toClient := range toClients {
			toClient.write(string(message))
		}
		go models.AddMessage(receiveMsg, 1)
	} else {
		go models.AddMessage(receiveMsg, 0)
		fmt.Println("没有发现client")
	}
	if toClients, ok := manager.singleClient[receiveMsg.Sendid]; ok {
		message, _ := json.Marshal(receiveMsg)
		for _, toClient := range toClients {
			toClient.write(string(message))
		}
	}
	//维护会话列表
	go UpdateChat(receiveMsg.Sendid, receiveMsg.Recid, receiveMsg.PublishedTime)
}

func UpdateChat(sendid string, recid string, publishedTime uint) {
	sendChats := models.FindChats(sendid, recid)
	if len(sendChats) > 0 {
		models.UpdateChatPublishedTime(sendid, recid, publishedTime)
	} else {
		models.AddChat(sendid, recid, publishedTime, 0)
	}

	recChats := models.FindChats(recid, sendid)
	if len(recChats) > 0 {
		models.UpdateChatPublishedTime(recid, sendid, publishedTime)
	} else {
		models.AddChat(recid, sendid, publishedTime, 0)
	}
}

func (c *Client) write(msg string) {
	err := c.socket.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		c.socket.Close()
		fmt.Println(err)
	}
	fmt.Print("发送成功")
}

func (manager *ClientManger) register(client *Client, receive ReceiveMessage) {
	_, ok := manager.clients[client]
	if !ok {
		manager.AddClients(client, receive.Uuid, receive.Type)
	}
}

func (manager *ClientManger) AddClients(client *Client, uuid string, connectType string) {
	manager.clientsLock.Lock()
	defer manager.clientsLock.Unlock()
	manager.clients[client] = uuid + "-" + connectType
	if _, ok := manager.singleClient[uuid]; ok {
		manager.singleClient[uuid][connectType] = client
	} else {
		manager.singleClient[uuid] = make(map[string]*Client, 10)
		manager.singleClient[uuid][connectType] = client
	}
	fmt.Println("add: ", uuid)
	fmt.Println("client:  ", client)
}

func (manager *ClientManger) DelClients(client *Client, uuid string) {
	manager.clientsLock.Lock()
	defer manager.clientsLock.Unlock()
	if _, ok := manager.clients[client]; ok {
		delete(manager.clients, client)
	}
	_, ok := manager.singleClient[uuid]
	if ok {
		splits := strings.Split(uuid, "-")
		delete(manager.singleClient[splits[0]], splits[1])
	}
}

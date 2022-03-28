package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/github-user/IMv_3/kafka"
	"github.com/github-user/IMv_3/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
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
	// Buffered channel of outbound messages.
	send chan []byte
}

type ReceiveMessage struct {
	Uuid string `json:"uuid"`
	Type string `json:"type"`
}

type WsPush struct {
	// 消息类型 1:登陆成功，2:收到新消息，3:PONG，4:发送消息的ACK, 5: 错误
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

const (
	ACKPushType   = 1
	ErrorPushType = -1
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var manager = ClientManger{
	clients:      make(map[*Client]string, 10000),
	singleClient: make(map[string]map[string]*Client, 10000),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	//启动消费者(轮询)
	revMessage()
}

func ConnectWeb(conn *websocket.Conn, sessionid string, connectType string) {
	var receiveMsg ReceiveMessage
	receiveMsg.Uuid = sessionid
	receiveMsg.Type = connectType
	client := &Client{socket: conn, send: make(chan []byte, 256)}
	client.message = receiveMsg
	manager.register(client, receiveMsg)
	//查询未读消息
	messages := models.FindUnReadMessageByRecid(sessionid)
	//把未读消息推送
	for _, message := range messages {
		if toClients, ok := manager.singleClient[sessionid]; ok {
			msg, _ := json.Marshal(message)
			for _, toClient := range toClients {
				toClient.write(msg)
			}
		}
		models.UpdateUnReadMessageById(message.ID)
	}

	go client.read()

}

func (c *Client) read() {
	defer func() {
		c.socket.Close()
	}()
	c.socket.SetReadLimit(maxMessageSize)
	//c.socket.SetReadDeadline(time.Now().Add(pongWait))
	c.socket.SetPongHandler(func(string) error { c.socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.send <- message
		go receiveMessage(c.send, c)
	}
}

func receiveMessage(send chan []byte, c *Client) {
	wsPush := WsPush{Type: ErrorPushType}
	flag := true
	msg := <-send
	var receiveMsg models.Message
	err := json.Unmarshal(msg, &receiveMsg)
	if err != nil {
		flag = false
		wsPush.Data = err.Error()
	}
	//保存消息
	err = models.AddMessage(&receiveMsg)
	if err != nil {
		flag = false
		wsPush.Data = err.Error()
	}
	msg, err = json.Marshal(receiveMsg)
	if err != nil {
		flag = false
		wsPush.Data = err.Error()
	}
	//发送消息到kafka
	kafka.SendMessage(msg)
	if err != nil {
		flag = false
		wsPush.Data = err.Error()
	}
	if flag {
		wsPush = WsPush{Type: ACKPushType}
		pushMessage, _ := json.Marshal(wsPush)
		c.write(pushMessage)
	} else {
		pushMessage, _ := json.Marshal(wsPush)
		c.write(pushMessage)
	}
}

func HandleMessage(receiveMsg models.Message) {
	if toClients, ok := manager.singleClient[receiveMsg.Recid]; ok {
		message, _ := json.Marshal(receiveMsg)
		for _, toClient := range toClients {
			toClient.write(message)
		}
		//把消息状态改为已发送(已读)
		models.UpdateUnReadMessageById(receiveMsg.ID)
	} else {
		fmt.Println("没有发现client")
	}
	//维护会话列表
	UpdateChat(receiveMsg.Sendid, receiveMsg.Recid, receiveMsg.PublishedTime)
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

func (c *Client) write(msg []byte) {
	err := c.socket.WriteMessage(websocket.TextMessage, msg)
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

func revMessage() {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("im") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("im", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		//defer pc.AsyncClose()
		// 异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("consumer: Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
				var receiveMsg models.Message
				json.Unmarshal(msg.Value, &receiveMsg)
				HandleMessage(receiveMsg)
			}
		}(pc)
	}
}

package models

type Message struct {
	ID            int    `gorm:"primary_key" json:"id"`
	Sendid        string `json:"sendid"`
	Recid         string `json:"recid"` //接收者的id是用户的id
	Content       string `json:"content"`
	PublishedTime uint   `json:"publishedTime"`
	State         int    `json:"state"` //0:未读 1:已读
}

type MessageSort []Message

//PersonSort 实现sort SDK 中的Interface接口

func (s MessageSort) Len() int {
	//返回传入数据的总数
	return len(s)
}
func (s MessageSort) Swap(i, j int) {
	//两个对象满足Less()则位置对换
	//表示执行交换数组中下标为i的数据和下标为j的数据
	s[i], s[j] = s[j], s[i]
}
func (s MessageSort) Less(i, j int) bool {
	//按字段比较大小,此处是降序排序
	//返回数组中下标为i的数据是否小于下标为j的数据
	return s[i].PublishedTime > s[j].PublishedTime
}

func AddMessage(message *Message) error {
	return db.Create(&message).Error
}

func FindUnReadMessageByRecid(recid string) []Message {
	var messages []Message
	db.Where("recid = ? and state = ?", recid, 0).Order("published_time desc").Find(&messages)
	return messages
}

func UpdateUnReadMessageById(id int) error {
	return db.Model(&Message{}).Where("id = ? and state = ?", id, 0).Update("state", 1).Error
}

func FindMessagesByRecid(recid string, startTime, endTime uint64) []Message {
	var messages []Message
	db.Where("recid = ? and  published_time >= ? and published_time <= ?", recid, startTime, endTime).Find(&messages)
	return messages
}

func FindMessagesBySendid(sendid string, startTime, endTime uint64) []Message {
	var messages []Message
	db.Where("sendid = ? and  published_time >= ? and published_time <= ?", sendid, startTime, endTime).Find(&messages)
	return messages
}

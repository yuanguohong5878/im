package models

type Chat struct {
	ID            int    `gorm:"primary_key" json:"id"`
	Sendid        string `json:"sendid"`
	Recid         string `json:"recid"` //接收者的id是用户的id
	PublishedTime uint   `json:"publishedTime"`
	State         int    `json:"state"` //0:未读 1:已读 2:置顶
}

func AddChat(sendid string, recid string, publishedTime uint, state int) error {
	return db.Create(&Chat{Sendid: sendid, Recid: recid, PublishedTime: publishedTime, State: state}).Error
}

func UpdateChatState(id string, state int) error {
	return db.Model(&Chat{}).Where("id = ? ", id).Update("state", state).Error
}

func UpdateChatPublishedTime(sendid string, recid string, publishedTime uint) error {
	return db.Model(&Chat{}).Where("sendid = ? and recid = ?", sendid, recid).Update("published_time", publishedTime).Error
}

func FindChats(sendid string, recid string) []Chat {
	var chats []Chat
	db.Where("sendid = ? and recid = ?", sendid, recid).Find(&chats)
	return chats
}

func FindChatsBySendidTop(sendid string) []Chat {
	var chats []Chat
	db.Where("sendid = ? and state = ?", sendid, 2).Order("published_time desc").Find(&chats)
	return chats
}

func FindChatsBySendidNotTop(sendid string) []Chat {
	var chats []Chat
	db.Where("sendid = ? and state != ?", sendid, 2).Order("published_time desc").Find(&chats)
	return chats
}

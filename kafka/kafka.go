package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

var client sarama.SyncProducer

// Setup initializes the database instance
func Setup() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true

	// 连接kafka
	var err error
	client, _ = sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}

}

func SendMessage(message []byte) error {
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "im"
	msg.Value = sarama.StringEncoder(message)

	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
	}
	fmt.Printf("produce: pid:%v offset:%v\n", pid, offset)
	return err
}

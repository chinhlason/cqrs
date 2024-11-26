package sync

import (
	"context"
	"cqrs-postgres-elastic-search-debezium/config"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type ConsumerGroupHandler struct{}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func ConsumeMessage(broker, topic, groupId string) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_6_0_0

	consumer, err := sarama.NewConsumerGroup([]string{broker}, groupId, config)
	if err != nil {
		return err
	}
	defer consumer.Close()

	handler := &ConsumerGroupHandler{}

	fmt.Println("Start consuming messages...")

	for {
		err := consumer.Consume(context.Background(), []string{topic}, handler)
		if err != nil {
			return err
		}
	}
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: topic = %s, partition = %d, offset = %d, key = %s, value = %s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		if msg.Topic == config.USER_TOPIC {
		}
		if msg.Topic == config.ORDER_TOPIC {
			var message OrderMessage
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Println("Error unmarshalling message: ", err)
				return err
			}
			log.Printf("Decoded message: %+v %s", *message.After, message.Op)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

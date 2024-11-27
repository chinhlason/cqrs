package sync

import (
	"context"
	"cqrs-postgres-elastic-search-debezium/query"
	"cqrs-postgres-elastic-search-debezium/utils"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type ConsumerGroupHandler struct {
	es *query.ESClient
}

func NewConsumerGroupHandler(es *query.ESClient) *ConsumerGroupHandler {
	return &ConsumerGroupHandler{es: es}
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeMessage(broker, topic, groupId string) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.V2_6_0_0

	consumer, err := sarama.NewConsumerGroup([]string{broker}, groupId, config)
	if err != nil {
		return err
	}
	defer consumer.Close()

	fmt.Printf("Start consuming messages at topic %s...\n", topic)

	for {
		err := consumer.Consume(context.Background(), []string{topic}, h)
		if err != nil {
			return err
		}
	}
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	errChan := make(chan error)
	for msg := range claim.Messages() {
		log.Printf("Message claimed: topic = %s, partition = %d, offset = %d, key = %s, value = %s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		if msg.Topic == utils.USER_TOPIC {
			var message UserMessage
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Println("Error unmarshalling message: ", err)
				errChan <- err
			}
			idDocument := fmt.Sprintf("user_%d", message.After.Id)
			if message.Op == "c" {
				err := h.es.InsertDocument(utils.USER_INDEX, idDocument, message.After)
				if err != nil {
					log.Println("Error inserting document into ES: ", err)
					errChan <- err
				}
			}
			if message.Op == "u" {
				err := h.es.UpdateDocument(idDocument, utils.USER_INDEX, message.After)
				if err != nil {
					log.Println("Error updating document into ES: ", err)
					errChan <- err
				}
			}
		}
		if msg.Topic == utils.ORDER_TOPIC {
			var message OrderMessage
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Println("Error unmarshalling message: ", err)
				errChan <- err
			}
			log.Printf("Decoded message: %+v %s", *message.After, message.Op)
			if message.Op == "c" {
				idOrder := fmt.Sprintf("order_%d", message.After.Id)
				idUser := fmt.Sprintf("user_%d", message.After.UserId)
				if err := h.es.InsertDocument(utils.ORDER_INDEX, idOrder, message.After); err != nil {
					log.Println("Error inserting document into ES: ", err)
					errChan <- err
				}
				if err := h.updateUserOrders(idUser, message.After); err != nil {
					log.Println("Error updating user orders: ", err)
					errChan <- err
				}
			}
			if message.Op == "u" {
				idOrder := fmt.Sprintf("order_%d", message.After.Id)
				idUser := fmt.Sprintf("user_%d", message.After.UserId)
				if err := h.es.UpdateDocument(idOrder, utils.ORDER_INDEX, message.After); err != nil {
					log.Println("Error updating document into ES: ", err)
					errChan <- err
				}
				if err := h.updateUserOrders(idUser, message.After); err != nil {
					log.Println("Error updating user orders: ", err)
					errChan <- err
				}
			}
		}
		go func() {
			if err := <-errChan; err != nil {
				session.MarkMessage(msg, "")
				return
			}
		}()
		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *ConsumerGroupHandler) updateUserOrders(idUser string, order *Order) error {
	record, err := h.es.GetDocument(idUser, utils.USER_INDEX)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	var recordJson ESRecord
	var userDoc User
	if err := json.Unmarshal(jsonData, &recordJson); err != nil {
		return err
	}

	userJson, err := json.Marshal(recordJson.Source)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(userJson, &userDoc); err != nil {
		return err
	}

	idFound := false

	for i, o := range userDoc.Orders {
		if o.Id == order.Id {
			userDoc.Orders[i] = *order
			idFound = true
			break
		}
	}

	if !idFound {
		userDoc.Orders = append(userDoc.Orders, *order)
	}

	if err := h.es.UpdateDocument(idUser, utils.USER_INDEX, userDoc); err != nil {
		return err
	}

	return nil
}

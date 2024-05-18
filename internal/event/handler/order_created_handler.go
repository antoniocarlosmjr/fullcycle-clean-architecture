package handler

import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"sync"

	"github.com/fullcycle-clean-architecture/pkg/events"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp091.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp091.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order created: %#v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp091.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	err := h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitmq,  // message to publish
	)

	if err != nil {
		return
	}
}

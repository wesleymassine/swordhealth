package notification

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/wesleymassine/swordhealth/task-management/domain"
)

var rabbitConn *amqp.Connection

// SetupRabbitMQConnection establishes a connection to RabbitMQ
func SetupRabbitMQConnection() (*amqp.Connection, error) {
	if rabbitConn == nil || rabbitConn.IsClosed() { // Check if connection is nil or closed
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/") // TODO VARS
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %s", err) // Log instead of exiting
			return nil, err
		}
		rabbitConn = conn
	}

	return rabbitConn, nil
}

// PublishToTopicExchange publishes a message to the specified topic exchange
func PublishToTopicExchange(routingKey string, taskMsg domain.Task) error {
	// Ensure we have a valid connection before publishing
	if rabbitConn == nil || rabbitConn.IsClosed() {
		_, err := SetupRabbitMQConnection() // Reconnect if needed
		if err != nil {
			return err // Return error if reconnection fails
		}
	}

	ch, err := rabbitConn.Channel() // Allocate a channel
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the topic exchange
	err = ch.ExchangeDeclare(
		"task_notifications", // name of the exchange
		"topic",              // type of exchange
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return err
	}

	taskMsgBody, err := json.Marshal(taskMsg)
	if err != nil {
		return err
	}

	// Publish the message to the exchange
	err = ch.Publish(
		"task_notifications", // exchange
		routingKey,           // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        taskMsgBody,
		})
	if err != nil {
		return err
	}

	log.Printf("Message published to exchange %s with routing key: %s:", "task_notifications", routingKey)

	return nil
}

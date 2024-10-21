package mq

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/streadway/amqp"
	"github.com/wesleymassine/swordhealth/user-notification/domain"
)

type RabbitMQConsumer struct {
	channel *amqp.Channel
}

func NewRabbitMQConsumer(channel *amqp.Channel) domain.NotificationEvent {

	// Declare the topic exchange for task notifications
	err := channel.ExchangeDeclare(
		"task_notifications", // exchange name
		"topic",              // exchange type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Declare the task_status_create_queue
	_, err = channel.QueueDeclare(
		"task_status_create_queue", // queue name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare task_status_create_queue: %v", err)
	}

	// Bind the task_status_create_queue to the exchange with routing key task.status.create
	err = channel.QueueBind(
		"task_status_create_queue", // queue name
		"task.status.create",       // routing key
		"task_notifications",       // exchange name
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind task_status_create_queue: %v", err)
	}

	// Declare the task_status_update_queue
	_, err = channel.QueueDeclare(
		"task_status_update_queue", // queue name
		true,                       // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare task_status_update_queue: %v", err)
	}

	// Bind the task_status_update_queue to the exchange with routing key task.status.update
	err = channel.QueueBind(
		"task_status_update_queue", // queue name
		"task.status.update",       // routing key
		"task_notifications",       // exchange name
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind task_status_update_queue: %v", err)
	}

	// Declare the dead_letter_queue and its exchange
	err = channel.ExchangeDeclare(
		"dead_letter_exchange", // exchange name
		"topic",                // exchange type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare dead_letter_exchange: %v", err)
	}

	_, err = channel.QueueDeclare(
		"dead_letter_queue", // queue name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare dead_letter_queue: %v", err)
	}

	return &RabbitMQConsumer{channel: channel}
}

// StartConsuming consumes messages from RabbitMQ and pushes them into the provided task channel
// Uses context for graceful shutdown and wg for concurrency control
func (c *RabbitMQConsumer) StartConsuming(ctx context.Context, queueName string, taskChannel chan<- domain.Task, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the WaitGroup counter when done

	msgs, err := c.channel.Consume(
		queueName,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		log.Fatalf("Failed to start consuming RabbitMQ messages: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done(): // Graceful shutdown signal from context
			log.Println("Shutting down RabbitMQ consumer...")
			close(taskChannel) // Close task channel when done
			return

		case d := <-msgs: // Received message from RabbitMQ
			var task domain.Task
			err := json.Unmarshal(d.Body, &task)

			if err != nil {
				log.Printf("Error decoding message: %v", err)
				c.sendToDeadLetterQueue(d) // Send failed message to dead-letter queue
				continue
			}

			task.Event = d.RoutingKey
			taskChannel <- task // Send the task to the channel
		}
	}
}

// sendToDeadLetterQueue handles sending failed messages to the dead-letter queue
func (c *RabbitMQConsumer) sendToDeadLetterQueue(d amqp.Delivery) {
	err := c.channel.Publish(
		"dead_letter_exchange", // dead-letter exchange
		"dead_letter",          // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        d.Body, // Original message body
		},
	)

	if err != nil {
		log.Printf("Failed to publish message to dead-letter queue: %v", err)
		return
	}
}

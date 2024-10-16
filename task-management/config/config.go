package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

func NewMySQLConnection() (*sql.DB, error) {
    db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/task")
    if err != nil {
        return nil, err
    }
    return db, nil
}

func NewRabbitMQChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return nil, err
	}

	return channel, nil
}
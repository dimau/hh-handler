package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func initializeRabbitMQConnection(rabbitMQUser, rabbitMQPass, rabbitMQServerName, rabbitMQPort string) (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	// Create connection to RabbitMQ
	connectionURL := fmt.Sprintf("amqp://%v:%v@%v:%v/", rabbitMQUser, rabbitMQPass, rabbitMQServerName, rabbitMQPort)
	conn, err := amqp.Dial(connectionURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	// Create channel to RabbitMQ
	rabbitChannel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	// Create HeadHunter queue (if it's already exist just get this channel)
	args := make(amqp.Table)
	args["x-message-ttl"] = int32(86400000)
	q, err := rabbitChannel.QueueDeclare(
		"HeadHunter", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		args,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return conn, rabbitChannel, q
}

func getMessagesFromQueue(rabbitChannel *amqp.Channel, q amqp.Queue) <-chan amqp.Delivery {
	messages, err := rabbitChannel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	return messages
}

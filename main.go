package main

import (
	"encoding/json"
	"github.com/dimau/hh-api-client-go"
	secrets "github.com/ijustfool/docker-secrets"
	"log"
	"os"
)

func main() {
	// Get environment vars for the application
	rabbitMQServerName := os.Getenv("RABBIT_MQ_SERVER_NAME")
	rabbitMQPort := os.Getenv("RABBIT_MQ_PORT")
	postgresServerName := os.Getenv("POSTGRES_SERVER_NAME")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_DB")

	// Get Docker secrets
	dockerSecrets, _ := secrets.NewDockerSecrets("")
	rabbitMQUser, _ := dockerSecrets.Get("rabbitmq_user")
	rabbitMQPass, _ := dockerSecrets.Get("rabbitmq_passwd")
	postgresUser, _ := dockerSecrets.Get("postgres_user")
	postgresPass, _ := dockerSecrets.Get("postgres_passwd")

	// Initialize RabbitMQ connection
	rabbitConn, rabbitChannel, rabbitHHQueue := initializeRabbitMQConnection(rabbitMQUser, rabbitMQPass, rabbitMQServerName, rabbitMQPort)
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	// Initialize Postgres connection
	db := initializePostgresConnection(postgresUser, postgresPass, postgresServerName, postgresPort, postgresDB)
	defer db.Close()

	// Handle messages from RabbitMQ
	messages := getMessagesFromQueue(rabbitChannel, rabbitHHQueue)
	for msg := range messages {

		// Unmarshall vacancy
		var vacancy hh.Vacancy
		err := json.Unmarshal(msg.Body, &vacancy)
		failOnError(err, "Failed Unmarshalling message from RabbitMQ to Vacancy struct")

		// Save a new vacancy to Database
		insertVacancy(db, &vacancy)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

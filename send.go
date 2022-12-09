package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn := con()
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Faild to open channel")
	defer ch.Close()

	/* create queue*/
	q, err := ch.QueueDeclare(
		"queue_yellow",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Faild to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Yellow green blue"
	err = ch.PublishWithContext(
		ctx,    // context
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

func con() *amqp.Connection {
	conn, err := amqp.Dial("amqp://rabbit:password@localhost:5672")
	failOnError(err, "Faild to connect to RabbitMQ")

	return conn
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

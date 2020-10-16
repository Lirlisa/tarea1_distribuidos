package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"sync"

	"./CamionLogistica"
	"./ClienteLogistica"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

func main() {
	var wait sync.WaitGroup
	//setear seed para los números de seguimiento
	rand.Seed(time.Now().UnixNano())

	listenerCliente, err1 := net.Listen("tcp", ":9000")
	listenerCamion, err2 := net.Listen("tcp", ":9001")
	if err1 != nil {
		log.Fatalf("Se ha producido un error: %s", err1)
	}
	if err2 != nil {
		log.Fatalf("Se ha producido un error: %s", err2)
	}

	fmt.Println("Iniciado servidor en escucha en puerto 9000")

	wait.Add(3)
	go func() {
		escuchar_cliente(listenerCliente)
		wait.Done()
	}()
	go func() {
		escuchar_camion(listenerCamion)
		wait.Done()
	}()
	go func() {
		conectar_finanas()
		wait.Done()
	}()

	wait.Wait()
}

func escuchar_cliente(listener net.Listener) {
	servidorCliente := ClienteLogistica.ServerCliente{}
	grpcServer := grpc.NewServer()
	ClienteLogistica.RegisterInteraccionesServer(grpcServer, &servidorCliente)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func escuchar_camion(listener net.Listener) {
	servidorCamion := CamionLogistica.ServerCamion{}
	grpcServer := grpc.NewServer()
	CamionLogistica.RegisterInteraccionesServer(grpcServer, &servidorCamion)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func conectar_finanas() {
	conn, err := amqp.Dial("amqp://test:test@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
}

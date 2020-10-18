package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/streadway/amqp"
)

var gananciasGeneral int
var perdidasGeneral int
var totalGeneral int

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("Ganancias: %d \n", gananciasGeneral)
		fmt.Printf("Perdidas: %d \n", perdidasGeneral)
		fmt.Printf("Total: %d \n", totalGeneral)
		os.Exit(0)
	}()

	conn, err := amqp.Dial("amqp://admin:password@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	log.Printf("conecta3")
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
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	file, err := os.Create("resumen.txt")
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		for d := range msgs {
			gananciapedido := 0
			perdidapedido := 0
			totalpedido := 0
			estado := ""
			textoGanancias := ""
			textoPerdidas := ""
			textoTotal := ""
			intento := ""

			gananciapedido = 0
			perdidapedido = 0
			totalpedido = 0

			info := make(map[string]interface{})
			json.Unmarshal([]byte(d.Body), &info)

			if str1, ok := (info["tipo"]).(string); ok {
				if str2, ok := (info["valor"]).(string); ok {
					if str3, ok := (info["estado"]).(string); ok {
						gananciapedido = calcularGanancias(str1, str2, str3)
						textoGanancias = "GANANCIAS: " + strconv.Itoa(gananciapedido)
					}
				}
			}
			if intentos, ok := (info["intentos"]).(string); ok {
				intento = "INTENTOS: " + intentos
				i, _ := strconv.Atoi(intentos)
				perdidapedido = 10 * (i - 1)
				textoPerdidas = "PERDIDAS: " + strconv.Itoa(perdidapedido)
			}

			totalpedido = gananciapedido - perdidapedido
			textoTotal = "TOTAL: " + strconv.Itoa(totalpedido)

			if str, ok := (info["estado"]).(string); ok {
				if str == "0" {
					estado = "NO ENTREGADO"
				} else {
					estado = "COMPLETADO"
				}
			}
			if str, ok := (info["id"]).(string); ok {
				file.WriteString(str + " " + estado + " " + intento + " " + textoGanancias + " " + textoPerdidas + " " + textoTotal + "\n")
			}

			gananciasGeneral = gananciasGeneral + gananciapedido
			perdidasGeneral = perdidasGeneral + perdidapedido
			totalGeneral = totalGeneral + totalpedido

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func calcularGanancias(tipo string, valor string, estado string) int {
	g, _ := strconv.Atoi(valor)
	if estado != "0" {
		if tipo == "prioritario" {
			return int(float32(g) * 1.3)
		}
		return g
	}
	if tipo == "retail" {
		return g
	}
	if tipo == "prioritario" {
		f := float32(g) * 0.3
		return int(f)
	}
	return 0

}

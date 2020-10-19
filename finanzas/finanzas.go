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

	//Si se ingresa CTRL C por terminal se muestran los datos correspondientes
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n")
		fmt.Printf("Ganancias: %d \n", gananciasGeneral)
		fmt.Printf("Perdidas: %d \n", perdidasGeneral)
		fmt.Printf("Total: %d \n", totalGeneral)
		os.Exit(0)
	}()

	//Se establece comunicación y estructuras necesarias para esto
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

	//archivo de registro para finanzas
	file, err := os.Create("resumen.txt")
	if err != nil {
		log.Fatal(err)
	}

	//Se leen los mensajes
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

			//Se recibe mensaje
			info := make(map[string]interface{})
			json.Unmarshal([]byte(d.Body), &info)

			//si hay mensaje
			if f, ok := (info["terminado"]).(string); ok {
				if f == "0" {
					//se realiza el cálculo de las ganancias del paquete
					if str1, ok := (info["tipo"]).(string); ok {
						if str2, ok := (info["valor"]).(string); ok {
							if str3, ok := (info["estado"]).(string); ok {
								gananciapedido = calcularGanancias(str1, str2, str3)
								textoGanancias = "GANANCIAS: " + strconv.Itoa(gananciapedido)
							}
						}
					}
					//Se ven las perdidas del paquete
					if intentos, ok := (info["intentos"]).(string); ok {
						intento = "INTENTOS: " + intentos
						i, _ := strconv.Atoi(intentos)
						perdidapedido = 10 * (i - 1)
						textoPerdidas = "PERDIDAS: " + strconv.Itoa(perdidapedido)
					}

					//Total del paquete
					totalpedido = gananciapedido - perdidapedido
					textoTotal = "TOTAL: " + strconv.Itoa(totalpedido)

					//Se mapea el estado
					if str, ok := (info["estado"]).(string); ok {
						if str == "0" {
							estado = "NO ENTREGADO"
						} else {
							estado = "COMPLETADO"
						}
					}
					//Se escribe en el registro
					if str, ok := (info["id"]).(string); ok {
						file.WriteString("Pedido: " + str + " " + estado + " " + intento + " " + textoGanancias + " " + textoPerdidas + " " + textoTotal + "\n")
					}

					//Se actualizan los valores generales
					gananciasGeneral = gananciasGeneral + gananciapedido
					perdidasGeneral = perdidasGeneral + perdidapedido
					totalGeneral = totalGeneral + totalpedido

				}
			}

		}
	}()

	log.Printf("[*] Presione Ctrl+C para terminar ejecución y recibir los totales:\n")
	<-forever
}

//calcula el valor de la ganancia
func calcularGanancias(tipo string, valor string, estado string) int {
	g, _ := strconv.Atoi(valor)
	if estado != "0" { //si se ntrega
		if tipo == "prioritario" {
			return int(float32(g) * 1.3) // si es prioritario es su valor más el 30% de cobro adicional
		}
		return g //cualquier otro caso es el valor
	} //si no se entrega
	if tipo == "retail" {
		return g //si es retail se obtiene el valor del producto
	}
	if tipo == "prioritario" { //si es prioritario se recibe el 30% del producto
		return int(float32(g) * 0.3)
	}
	return 0 //normal no se recibe nada

}

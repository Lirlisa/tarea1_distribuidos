package ClienteLogistica

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"../Estructuras"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

var id_disponible uint32

type ServerCliente struct {
	//server
	placeholder int
}

func (s *ServerCliente) Encargar(ctx context.Context, in *Encargo) (*Producto, error) {
	log.Printf("Se ha recibido encargo: %s", in.TipoLocal)
	if in.TipoLocal == "gg" {
		paqueteFinal()
		return &Producto{ID: 0}, nil
	}
	content := fmt.Sprintf("%s,%s,%d,%s,%s\n", in.GetTipoLocal(), in.GetNombreProducto(), in.GetValor(), in.GetOrigen(), in.GetDestino())
	var candado sync.Mutex
	candado.Lock()
	f, err := os.OpenFile("bitacora.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Ha ocurrido un error con el archivo: %s", err)
	}

	if _, err = f.WriteString(content); err != nil {
		f.Close()
		log.Panicf("Error al escribir en archivo: %s", err)
	}
	f.Close()
	id_disponible++
	idReservada := id_disponible

	test := rand.Uint32()
	for _, existe := Estructuras.SeguimientoAId[test]; existe; _, existe = Estructuras.SeguimientoAId[test] {
		test = rand.Uint32()
	}
	Estructuras.SeguimientoAId[test] = idReservada
	candado.Unlock()

	nuevoRegistro := new(Estructuras.Registro)
	*nuevoRegistro = Estructuras.Registro{
		time.Now(),
		idReservada,
		in.GetTipoLocal(),
		in.GetNombreProducto(),
		in.GetValor(),
		in.GetOrigen(),
		in.GetDestino(),
		int32(test),
	}
	Estructuras.Tabla[idReservada] = nuevoRegistro

	pack := Estructuras.Paquete{
		idReservada,
		uint32(test),
		in.GetTipoLocal(),
		in.GetValor(),
		0,
		0,
	}
	Estructuras.Paquetes[idReservada] = &pack

	switch x := in.GetTipoLocal(); x {
	case "retail":
		candado.Lock()
		Estructuras.ColaRetail = append(Estructuras.ColaRetail, pack)
		candado.Unlock()
	case "prioritario":
		candado.Lock()
		Estructuras.ColaPrioridad = append(Estructuras.ColaPrioridad, pack)
		candado.Unlock()
	case "normal":
		candado.Lock()
		Estructuras.ColaNormal = append(Estructuras.ColaNormal, pack)
		candado.Unlock()
	}

	return &Producto{ID: test}, nil

}
func (s *ServerCliente) EstadoEncargo(ctx context.Context, in *Producto) (*Estatus, error) {
	log.Printf("Solicitado estado de %d", in.ID)

	if valor, existencia := Estructuras.Paquetes[Estructuras.SeguimientoAId[in.ID]]; existencia {
		return &Estatus{Valor: int32(valor.Estado)}, nil
	}
	return &Estatus{Valor: 4}, nil //el extraño caso en que no exista el paquete
}

func paqueteFinal() {
	var candado sync.Mutex
	paquete := new(Estructuras.Paquete)
	paquete.Tipo = "gg"
	candado.Lock()
	for i := 0; i < 3; i++ {
		Estructuras.ColaRetail = append(Estructuras.ColaRetail, *paquete)
		Estructuras.ColaPrioridad = append(Estructuras.ColaPrioridad, *paquete)
		Estructuras.ColaNormal = append(Estructuras.ColaNormal, *paquete)
	}
	candado.Unlock()

	conn, err := amqp.Dial("amqp://admin:password@dist46:5672/")
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
	var body string
	body = `{terminado: "1", estado: "0", intentos: "0", ganancia: "0", tipo: "0", id: "0"}`

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

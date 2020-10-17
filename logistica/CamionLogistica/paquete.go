package CamionLogistica

import (
	"log"
	"strconv"
	"sync"

	"../Estructuras"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
)

func conectarFinanzas(terminado int32, estado uint32, ganancia uint32, tipo string, id uint32) {
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

	body := "
	{
		terminado: '"+strconv.Itoa(terminado)"',
		estado: '"+strconv.Itoa(estado)+"',
		intentos: '"+strconv.Itoa(intentos)+"',
		ganancia: '"+strconv.Itoa(ganancia)+"',
		tipo: '"+tipo+"',
		id: '"+strconv.Itoa(id)+"'
	}"
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

type ServerCamion struct {
	placeholder int
}

func (c *ServerCamion) PedirPaquete(ctx context.Context, in *Tipo) (*Paquete, error) {
	var candado sync.Mutex
	var a string
	if in.GetClase() == 1 {
		a = "Retail"
	} else {
		a = "Normal"
	}
	log.Printf("Pedido de paquete de tipo %s", a)
	elem := new(Estructuras.Paquete)

	if x := in.GetClase(); x == 1 {
		candado.Lock()
		if len(Estructuras.ColaRetail) > 0 {
			*elem = Estructuras.ColaRetail[0]
			Estructuras.ColaRetail = Estructuras.ColaRetail[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
		}
		candado.Unlock()
	} else if x == 2 {
		candado.Lock()
		if len(Estructuras.ColaPrioridad) > 0 {
			*elem = Estructuras.ColaPrioridad[0]
			Estructuras.ColaPrioridad = Estructuras.ColaPrioridad[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
		} else if len(Estructuras.ColaNormal) > 0 {
			*elem = Estructuras.ColaNormal[0]
			Estructuras.ColaNormal = Estructuras.ColaNormal[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
		}
		candado.Unlock()
	} else {
		candado.Lock()
		if len(Estructuras.ColaRetail) > 0 {
			*elem = Estructuras.ColaRetail[0]
			Estructuras.ColaRetail = Estructuras.ColaRetail[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
		} else if len(Estructuras.ColaPrioridad) > 0 {
			*elem = Estructuras.ColaPrioridad[0]
			Estructuras.ColaPrioridad = Estructuras.ColaPrioridad[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
		}
		candado.Unlock()
	}

	if item, existe := Estructuras.Tabla[elem.IDPaquete]; existe {
		log.Printf("Entregado paquete id: %d", elem.IDPaquete)
		return &Paquete{
			IDPaquete:   elem.IDPaquete,
			Seguimiento: elem.Seguimiento,
			Tipo:        elem.Tipo,
			Valor:       elem.Valor,
			Intentos:    elem.Intentos,
			Estado:      elem.Estado,
			Origen:      item.Origen,
			Destino:     item.Destino,
		}, nil
	} else {
		return &Paquete{
			IDPaquete:   0,
			Seguimiento: 0,
			Tipo:        "",
			Valor:       0,
			Intentos:    0,
			Estado:      0,
			Origen:      "",
			Destino:     "",
		}, nil
	}

}

func (c *ServerCamion) DevolverPaquete(ctx context.Context, in *Paquete) (*Paquete, error) {
	Estructuras.Paquetes[in.IDPaquete].Intentos = in.GetIntentos()
	Estructuras.Paquetes[in.IDPaquete].Estado = in.GetEstado()
	go conectarFinanzas(0, in.GetEstado(), in.GetValor(), in.GetTipo, in.GetIDPaquete())
	return in, nil
}

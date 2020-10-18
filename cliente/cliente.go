package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"./pedido"
	"google.golang.org/grpc"
)

var wait sync.WaitGroup
var tiempo int

func main() {
	//semilla para aleatoriedad
	rand.Seed(time.Now().UnixNano())

	//Se establece comunicación
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dist47:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	//Solicitud tiempo entre ordenes
	var tienda int
	fmt.Println("Tiempo de espera entre envio de ordenes (En segundos):")
	fmt.Scanln(&tiempo)
	for {
		fmt.Println("Ingrese comportamiento a seguir (1/2):\n1 RETAIL\n2 PYME")
		fmt.Scanln(&tienda)
		if tienda == 1 || tienda == 2 {
			break
		}
	}

	//Se ejecuta según corresponda
	if tienda == 1 {
		recibir("retail.csv", tienda, conn)
	} else {
		recibir("pymes.csv", tienda, conn)
	}
}

//Lee el archivo correspondiente, se realizan las oredenes correspondientes y se llama a la función encargada de realizar las consultas
func recibir(archivo string, tienda int, conn *grpc.ClientConn) {

	//Se abre el archivo
	file, err := os.Open(archivo)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer file.Close()

	//Se lee el archivo
	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error", err)
	}
	var aux1 int

	//iteramos por los datos
	for value := range record { // for i:=0; i<len(record)
		if value != 0 {

			//Establece el tipo de envío
			var tipo string
			if tienda == 2 {
				if record[value][5] == "1" {
					tipo = "prioritario"
				} else {
					tipo = "normal"
				}
			} else {
				tipo = "retail"
			}

			//Se envía el pedido
			cliente1 := pedido.NewInteraccionesClient(conn)
			aux1, _ = strconv.Atoi(record[value][2])
			response, err := cliente1.Encargar(context.Background(), &pedido.Encargo{
				TipoLocal:      tipo,
				NombreProducto: record[value][1],
				Valor:          uint32(aux1),
				Origen:         record[value][3],
				Destino:        record[value][4],
			})
			if err != nil {
				log.Fatalf("Error al Encargar: %s", err)
			}

			//Si la tienda es pyme se llama a la función encargada de realizar consultas de estado
			if tienda == 2 {
				var seguimiento uint32
				seguimiento = response.ID
				wait.Add(1)
				go cliente(seguimiento, conn)
			}

			//Se espera el tiempo especificado entre pedidos
			time.Sleep(time.Second * time.Duration(tiempo))

		}
	}
	wait.Wait()
	cliente1 := pedido.NewInteraccionesClient(conn)
	response, err := cliente1.Encargar(context.Background(), &pedido.Encargo{
		TipoLocal:      "gg",
		NombreProducto: "",
		Valor:          0,
		Origen:         "",
		Destino:        "",
	})
	if err != nil {
		log.Fatalf("Error al Encargar: %s", err)
	}
	fmt.Printf("Terminado %d", response.ID)

}

//realiza consultas de estado, si en alguna iteración no realiza consultas entonces deja de hacerlo.
func cliente(seguimiento uint32, conn *grpc.ClientConn) {
	var consulta int
	var estado string

	for {
		consulta = rand.Intn(100)
		// 60% de probabilidad de realizar consultas
		if consulta < 40 { //////////////////////////////////////VALOR MODIFICABLE
			break
		}
		//Se realiza consulta
		time.Sleep(time.Second * 10) /////////////////////////////VALOR MODIFICABLE
		fmt.Printf("QUIERO CONSULTAR POR MI PEDIDO: %d\n", seguimiento)

		cliente := pedido.NewInteraccionesClient(conn)

		response, err := cliente.EstadoEncargo(context.Background(), &pedido.Producto{
			ID: seguimiento,
		})
		if err != nil {
			log.Fatalf("Error al consultar estado: %s", err)
		}
		// Mapeo estado edido
		estado = estadoPedido(response.Valor)
		fmt.Printf("EL ESTADO DEL PEDIDO:%d ES %s\n", seguimiento, estado)

	}

	wait.Done()
}

// Mapea el valor, da un significado a la respuesta
func estadoPedido(valor int32) string {
	if valor == 0 {
		return "EN BODEGA"
	} else if valor == 1 {
		return "EN CAMINO"
	} else if valor == 2 {
		return "RECIBIDO"
	}
	return "NO RECIBIDO"
}

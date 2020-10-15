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

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	/*#######################################################*/
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

	if tienda == 1 {
		recibir("retail.csv", tienda, conn)
	} else {
		recibir("pymes.csv", tienda, conn)
	}
}

func recibir(archivo string, tienda int, conn *grpc.ClientConn) {
	file, err := os.Open(archivo)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error", err)
	}
	var aux1 int
	for value := range record { // for i:=0; i<len(record)
		if value != 0 {

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

			if tienda == 2 {
				var seguimiento uint32
				seguimiento = response.ID
				wait.Add(1)
				go cliente(seguimiento, conn)
			}

			time.Sleep(time.Second * time.Duration(tiempo))

		}
	}
	wait.Wait()

}

func cliente(seguimiento uint32, conn *grpc.ClientConn) {
	var consulta int
	var estado string

	for {
		consulta = rand.Intn(100)
		if consulta <= 40 { ////////////////////////////////////
			break
		}
		time.Sleep(time.Second * 5) /////////////////////////////
		fmt.Printf("QUIERO CONSULTAR POR MI PEDIDO: %d\n", seguimiento)

		cliente := pedido.NewInteraccionesClient(conn)

		response, err := cliente.EstadoEncargo(context.Background(), &pedido.Producto{
			ID: seguimiento,
		})
		if err != nil {
			log.Fatalf("Error al consultar estado: %s", err)
		}
		estado = estadoPedido(response.Valor)
		fmt.Printf("EL ESTADO DEL PEDIDO:%d ES %s\n", seguimiento, estado)

	}

	wait.Done()
}

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

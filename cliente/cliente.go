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

	"sync"
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

			cliente := pedido.NewInteraccionesClient(conn)

			response, err := cliente.Encargar(context.Background(), &pedido.Encargo{
				TipoLocal:      tipo,
				NombreProducto: record[value][1],
				Valor:          strconv.Atoi(record[value][2]),
				Origen:         record[value][3],
				Destino:        record[value][4],
			})
			if err != nil {
				log.Fatalf("Error al Encargar: %s", err)
			}

			if tienda == 2 {
				var seguimiento int
				seguimiento = response.ID
				wait.Add(1)
				go cliente(seguimiento, conn)
			}

			time.Sleep(time.Second * time.Duration(tiempo))

		}
	}
	wait.Wait()

}

func cliente(seguimiento int, conn *grpc.ClientConn) {
	var consulta int

	for {
		consulta = rand.Intn(100)
		if consulta <= 40 { ////////////////////////////////////
			break
		}
		time.Sleep(time.Second * 5) /////////////////////////////
		fmt.Println("QUIERO CONSULTAR POR MI PEDIDO: %d", seguimiento)

		cliente := pedido.NewInteraccionesClient(conn)

		response, err := cliente.Estado_encargo(context.Background(), &pedido.Producto{
			ID: seguimiento,
		})
		if err != nil {
			log.Fatalf("Error al consultar estado: %s", err)
		}

		fmt.Printf("ESTADO PEDIDO:%d,%s", seguimiento, response.Valor)

	}

	wait.Done()
}

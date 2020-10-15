package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"./CamionLogistica"
	"google.golang.org/grpc"
)

type paquete struct {
	ID          uint32
	tipo        string
	valor       uint32
	origen      string
	destino     string
	intentos    uint32
	fentrega    string
	seguimiento uint32
}

var tiempo int
var tiempoEntrega int
var wait sync.WaitGroup

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Tiempo de espera de los camiones (En segundos):")
	fmt.Scanln(&tiempo)
	fmt.Println("Tiempo que tardan en entregar los camiones (En segundos):")
	fmt.Scanln(&tiempoEntrega)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	wait.Add(3)
	go camion(1, conn, 1)
	go camion(1, conn, 2)
	go camion(2, conn, 3)
	wait.Wait()

}

func camion(tipo uint32, conn *grpc.ClientConn, vehiculo int) {

	registro := make(map[int]paquete)

	var ind1 int
	var ind2 int

	c := CamionLogistica.NewInteraccionesClient(conn)

	for {
		var envio1 paquete
		var envio2 paquete
		envio1.ID = 0
		envio1.tipo = ""
		envio1.valor = 0
		envio2.ID = 0
		envio2.tipo = ""
		envio2.valor = 0
		ind1 = 0
		ind2 = 0

		for {
			if envio1.ID == 0 {
				response, err := c.PedirPaquete(context.Background(), &CamionLogistica.Tipo{Clase: tipo})
				if err != nil {
					log.Fatalf("Error al pedir paquete: %s", err)
				}
				if response.IDPaquete != 0 {
					fmt.Printf("Camion %d tiene 1 paquete\n", vehiculo)
					envio1.ID = response.IDPaquete
					envio1.tipo = response.Tipo
					envio1.valor = response.Valor
					envio1.origen = response.Origen
					envio1.destino = response.Destino
					envio1.intentos = 0
					envio1.fentrega = "0"
					envio1.seguimiento = response.Seguimiento
					if ind1 == ind2 {
						ind1 = ind1 + 1
					} else {
						ind1 = ind2 + 1
					}
					registro[ind1] = envio1
					break
				}
				time.Sleep(time.Second * 5)
			}
		}
		time.Sleep(time.Second * time.Duration(tiempo))
		response, err := c.PedirPaquete(context.Background(), &CamionLogistica.Tipo{Clase: tipo})
		if err != nil {
			log.Fatalf("Error al pedir paquete: %s", err)
		}
		if response.IDPaquete != 0 {
			fmt.Printf("Camion %d tiene 2 paquete\n", vehiculo)
			envio2.ID = response.IDPaquete
			envio2.tipo = response.Tipo
			envio2.valor = response.Valor
			envio2.origen = response.Origen
			envio2.destino = response.Destino
			envio2.intentos = 0
			envio2.fentrega = "0"
			envio2.seguimiento = response.Seguimiento
			ind2 = ind1 + 1
			registro[ind2] = envio2
		}

		if tipo == 2 { ////camion normal aniliza los casos
			if envio1.tipo == "prioritario" {
				envio1, envio2 = reparto(envio1, envio2)
			} else if envio2.tipo == "prioritario" {
				envio2, envio1 = reparto(envio2, envio1)
			} else if envio1.valor >= envio2.valor {
				envio1, envio2 = reparto(envio1, envio2)
			} else {
				envio2, envio1 = reparto(envio2, envio1)
			}
		} else { //////camion retail analiza los casos
			if envio1.valor >= envio2.valor {
				envio1, envio2 = reparto(envio1, envio2)
			} else {
				envio2, envio1 = reparto(envio2, envio1)
			}
		}

		if envio1.ID != 0 {
			registro[ind1] = envio1
			var entrega uint32
			if envio1.fentrega == "0" {
				entrega = 3
			} else {
				entrega = 2
			}
			response, err := c.DevolverPaquete(context.Background(), &CamionLogistica.Paquete{
				IDPaquete:   envio1.ID,
				Seguimiento: envio1.seguimiento,
				Tipo:        envio1.tipo,
				Valor:       envio1.valor,
				Intentos:    envio1.intentos,
				Estado:      entrega,
				Origen:      envio1.origen,
				Destino:     envio1.destino})
			if err != nil {
				log.Fatalf("Error al pedir paquete: %s", err)
			}
			fmt.Printf("Entregada información pedido: %d\n", response.IDPaquete)

		}
		if envio2.ID != 0 {
			registro[ind2] = envio2
			var entrega uint32
			if envio2.fentrega == "0" {
				entrega = 3
			} else {
				entrega = 2
			}
			response, err := c.DevolverPaquete(context.Background(), &CamionLogistica.Paquete{
				IDPaquete:   envio2.ID,
				Seguimiento: envio2.seguimiento,
				Tipo:        envio2.tipo,
				Valor:       envio2.valor,
				Intentos:    envio2.intentos,
				Estado:      entrega,
				Origen:      envio2.origen,
				Destino:     envio2.destino})
			if err != nil {
				log.Fatalf("Error al pedir paquete: %s", err)
			}
			fmt.Printf("Entregada información pedido: %d\n", response.IDPaquete)

		}
	}
}

func reparto(envio1 paquete, envio2 paquete) (paquete, paquete) { /////simulación del camion en la calle
	var intento1 uint32
	var intento2 uint32
	var break1 int
	var break2 int
	intento1, break1 = cantintentos(envio1)
	intento2, break2 = cantintentos(envio2)
	time.Sleep(time.Second * time.Duration(tiempoEntrega))
	for {
		if envio1.valor == 0 || envio1.fentrega != "0" || break1 == 1 { /////intenta entregar primer pedido
			break1 = 1
		} else {
			envio1, break1 = entrega(envio1, intento1)
		}
		time.Sleep(time.Second * time.Duration(tiempoEntrega))
		if envio2.valor == 0 || envio2.fentrega != "0" || break2 == 1 { /////intenta entregar segundo pedido
			break2 = 1
		} else {
			envio2, break2 = entrega(envio2, intento2)
		}
		if break1 == 1 && break2 == 1 { ////si se entregan ambos o rompe condiciones vuelve
			break
		}
		time.Sleep(time.Second * time.Duration(tiempoEntrega))
	}
	return envio1, envio2
}

func entrega(envio paquete, intentos uint32) (paquete, int) { /// revisar condiciones para realizar entrega, retorna 1 si no se puede entregar
	if envio.valor > 10*envio.intentos && envio.intentos < intentos && envio.fentrega == "0" {
		return probentregar(envio), 0
	}
	return envio, 1
}

func cantintentos(envio paquete) (uint32, int) { /// entrega la cantidad de intentos maximos para el paquete
	if envio.tipo == "retail" {
		return 3, 0
	} else if envio.tipo == "" {
		return 0, 1
	}
	return 2, 0
}

func probentregar(envio paquete) paquete { ///realiza la simulación de entregar en domicilio modifica la fecha si es que entrega
	var prob int
	envio.intentos = envio.intentos + 1
	prob = rand.Intn(100)
	if prob >= 80 {
		envio.fentrega = "0"
	} else {
		present := time.Now()
		envio.fentrega = present.Format("01-02-2006 15:04:05")
	}
	return envio
}

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type paquete struct {
	ID       uint32
	tipo     string
	valor    uint32
	origen   string
	destino  string
	intentos uint32
	fentrega string
}

var tiempo int
var wait sync.WaitGroup

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Tiempo de espera de los camiones (En segundos):")
	fmt.Scanln(&tiempo)
	fmt.Printf("Tiempo:%d\n", tiempo)
	wait.Add(3)
	go camion(1)
	go camion(1)
	go camion(2)
	wait.Wait()
}

func camion(tipo int) {
	//registro := make(map[int]paquete)
	var envio1 paquete
	var envio2 paquete
	envio1.valor = 0
	envio2.valor = 0
	for { ///////modificar
		for {
			if envio1.valor == 0 {
				//esperar primer pedido
				envio1.intentos = 0
				envio1.fentrega = "0"
				break
			}
		}
		//pedir envio 2
		break
	} /////////
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
	} else { //////camion retail anilza los casos
		if envio1.valor >= envio2.valor {
			envio1, envio2 = reparto(envio1, envio2)
		} else {
			envio2, envio1 = reparto(envio2, envio1)
		}
	}
	fmt.Println(envio1.fentrega)
	wait.Done() ////////////////////////////////////////////////////////
	///////NOTIFICAR
}

func reparto(envio1 paquete, envio2 paquete) (paquete, paquete) { /////simulación del camion en la calle
	var intento1 uint32
	var intento2 uint32
	var break1 int
	var break2 int
	intento1 = cantintentos(envio1)
	intento2 = cantintentos(envio2)
	break1 = 0
	break2 = 0
	for {
		if envio1.valor == 0 || envio1.fentrega != "0" || break1 == 1 { /////intenta entregar primer pedido
			break1 = 1
		} else {
			envio1, break1 = entrega(envio1, intento1)
		}
		if envio2.valor == 0 || envio2.fentrega != "0" || break2 == 1 { /////intenta entregar segundo pedido
			break2 = 1
		} else {
			envio2, break2 = entrega(envio2, intento2)
		}
		if break1 == 1 && break2 == 1 { ////si se entregan ambos o rompe condiciones vuelve
			break
		}
		time.Sleep(time.Second * 5) /////////////////////////////////////////////////////7
	}
	return envio1, envio2
}

func entrega(envio paquete, intentos uint32) (paquete, int) { /// revisar condiciones para realizar entrega, retorna 1 si no se puede entregar
	if envio.valor > 10*envio.intentos && envio.intentos < intentos && envio.fentrega == "0" {
		return probentregar(envio), 0
	}
	return envio, 1
}

func cantintentos(envio paquete) uint32 { /// entrega la cantidad de intentos maximos para el paquete
	if envio.tipo == "retail" {
		return 3
	}
	return 2
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

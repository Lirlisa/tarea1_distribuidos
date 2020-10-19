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
	"./Estructuras"
	"google.golang.org/grpc"
)

func main() {
	//variable de sincronización de go routines
	var wait sync.WaitGroup

	//setear seed para los números de seguimiento
	rand.Seed(time.Now().UnixNano())

	//se abren los puertos
	listenerCliente, err1 := net.Listen("tcp", ":9000")
	listenerCamion, err2 := net.Listen("tcp", ":9001")
	if err1 != nil {
		log.Fatalf("Se ha producido un error: %s", err1)
	}
	if err2 != nil {
		log.Fatalf("Se ha producido un error: %s", err2)
	}

	fmt.Println("Iniciado servidor en escucha en puerto 9000 y 9001")

	//se ejecuta cada listener en rutinas distintas
	wait.Add(2)
	go func() {
		escucharCliente(listenerCliente)
		wait.Done()
	}()
	go func() {
		escucharCamion(listenerCamion)
		wait.Done()
	}()

	wait.Wait()
}

func escucharCliente(listener net.Listener) {
	//función encargada del servidor para escuchar al cliente, termina cuando ocurra un error es el servidor grpc
	//o en el listener
	servidorCliente := ClienteLogistica.ServerCliente{}
	Estructuras.GrpcServerCliente = grpc.NewServer()
	ClienteLogistica.RegisterInteraccionesServer(Estructuras.GrpcServerCliente, &servidorCliente)

	if err := Estructuras.GrpcServerCliente.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func escucharCamion(listener net.Listener) {
	//función encargada del servidor para escuchar al camion, termina cuando ocurra un error es el servidor grpc
	//o en el listener
	servidorCamion := CamionLogistica.ServerCamion{}
	Estructuras.GrpcServerCamion = grpc.NewServer()
	CamionLogistica.RegisterInteraccionesServer(Estructuras.GrpcServerCamion, &servidorCamion)

	if err := Estructuras.GrpcServerCamion.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

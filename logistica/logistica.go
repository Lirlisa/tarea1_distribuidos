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
	"google.golang.org/grpc"
)

func main() {
	var wait sync.WaitGroup
	//setear seed para los números de seguimiento
	rand.Seed(time.Now().UnixNano())

	listenerCliente, err1 := net.Listen("tcp", ":9000")
	listenerCamion, err2 := net.Listen("tcp", ":9001")
	if err1 != nil {
		log.Fatalf("Se ha producido un error: %s", err1)
	}
	if err2 != nil {
		log.Fatalf("Se ha producido un error: %s", err2)
	}

	fmt.Println("Iniciado servidor en escucha en puerto 9000")

	wait.Add(2)
	go func() {
		escucharCliente(listenerCliente)
		wait.Done()
	}()
	go func() {
		escucharCamion(listenerCamion)
		wait.Done()
	}()
	/*go func() {
		conectarFinanzas()
		wait.Done()
	}()*/

	wait.Wait()
}

func escucharCliente(listener net.Listener) {
	servidorCliente := ClienteLogistica.ServerCliente{}
	grpcServer := grpc.NewServer()
	ClienteLogistica.RegisterInteraccionesServer(grpcServer, &servidorCliente)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func escucharCamion(listener net.Listener) {
	servidorCamion := CamionLogistica.ServerCamion{}
	grpcServer := grpc.NewServer()
	CamionLogistica.RegisterInteraccionesServer(grpcServer, &servidorCamion)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

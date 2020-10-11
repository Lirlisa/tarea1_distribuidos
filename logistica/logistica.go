package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"./ClienteLogistica"
	"google.golang.org/grpc"
)

func main() {
	//setear seed para los números de seguimiento
	rand.Seed(time.Now().UnixNano())

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Se ha producido un error: %s", err)
	}

	fmt.Println("Iniciado servidor en escucha en puerto 9000")

	servidor := ClienteLogistica.Server{}
	grpcServer := grpc.NewServer()
	ClienteLogistica.RegisterInteraccionesServer(grpcServer, &servidor)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

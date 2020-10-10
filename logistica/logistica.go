package main

import (
	"fmt"
	"log"
	"net"

	"./pedido"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Se ha producido un error: %s", err)
	}

	fmt.Println("Iniciado servidor en escucha en puerto 9000")

	servidor := pedido.Server{}
	grpcServer := grpc.NewServer()
	pedido.RegisterInteraccionesServer(grpcServer, &servidor)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

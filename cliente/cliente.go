package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	cliente := chat.NewChatServiceClient(conn)

	response, err := cliente.Encargar(context.Background(), &chat.Encargo{
		tipo_local="retail",
		nombre_producto="Nada",
		valor=0,
		origen="tiendaA",
		destino="casaA",
	
	})
	if err != nil {
		log.Fatalf("Error al Encargar: %s", err)
	}
	log.Printf("Response from server:")

}

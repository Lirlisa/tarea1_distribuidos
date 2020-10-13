package main

import (
	"context"
	"log"

	"./CamionLogistica"

	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := CamionLogistica.NewInteraccionesClient(conn)

	response, err := c.PedirPaquete(context.Background(), &CamionLogistica.Tipo{Clase: 1})
	if err != nil {
		log.Fatalf("Error al pedir paquete: %s", err)
	}
	log.Printf("Response from server: %d", response.IDPaquete)

}

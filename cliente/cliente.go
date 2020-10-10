package main

import (
	"context"
	"log"

	"./pedido"
	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	cliente := pedido.NewInteraccionesClient(conn)

	response, err := cliente.Encargar(context.Background(), &pedido.Encargo{
		TipoLocal:      "retail",
		NombreProducto: "Nada",
		Valor:          0,
		Origen:         "tiendaA",
		Destino:        "casaA",
	})
	if err != nil {
		log.Fatalf("Error al Encargar: %s", err)
	}
	log.Printf("Response from server: %s", response.TipoLocal)

}

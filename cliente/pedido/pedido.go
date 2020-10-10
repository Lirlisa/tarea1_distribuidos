package pedido

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) Encargar(ctx context.Context, in *Encargo) (*Encargo, error) {
	log.Printf("Se ha recibido encargo")
	return in, nil
}
func (s *Server) EstadoEncargo(ctx context.Context, in *Producto) (*Producto, error) {
	log.Printf("Solicitado estado de %s", in.ID)
	return in, nil
}

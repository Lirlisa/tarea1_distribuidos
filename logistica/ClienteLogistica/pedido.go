package ClienteLogistica

import (
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"
)

var id_disponible uint32

type registro struct {
	timestamp   time.Time
	id          uint32
	tipo        string
	nombre      string
	valor       uint32
	origen      string
	destino     string
	seguimiento int
}

var tabla = make(map[uint32]registro)
var seguimientoAId = make(map[int]uint32)

type Server struct {
	//server
	placeholder int
}

func (s *Server) Encargar(ctx context.Context, in *Encargo) (*Producto, error) {
	idReservada := id_disponible
	id_disponible++

	log.Printf("Se ha recibido encargo: %s", in.TipoLocal)
	test := rand.Int()
	for _, existe := seguimientoAId[test]; existe; _, existe = seguimientoAId[test] {
		test = rand.Int()
	}
	tabla[idReservada] = registro{
		time.Now(),
		idReservada,
		in.GetTipoLocal(),
		in.GetNombreProducto(),
		in.GetValor(),
		in.GetOrigen(),
		in.GetDestino(),
		test,
	}
	seguimientoAId[test] = idReservada
	return &Producto{ID: idReservada}, nil

}
func (s *Server) EstadoEncargo(ctx context.Context, in *Producto) (*Estatus, error) {
	log.Printf("Solicitado estado de %s", in.ID)
	return &Estatus{Valor: 1}, nil
}

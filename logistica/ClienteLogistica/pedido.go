package ClienteLogistica

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"../Estructuras"
	"golang.org/x/net/context"
)

var id_disponible uint32

type ServerCliente struct {
	//server
	placeholder int
}

func (s *ServerCliente) Encargar(ctx context.Context, in *Encargo) (*Producto, error) {
	log.Printf("Se ha recibido encargo: %s", in.TipoLocal)
	var candado sync.Mutex
	candado.Lock()
	id_disponible++
	idReservada := id_disponible

	test := rand.Uint32()
	for _, existe := Estructuras.SeguimientoAId[test]; existe; _, existe = Estructuras.SeguimientoAId[test] {
		test = rand.Uint32()
	}
	Estructuras.SeguimientoAId[test] = idReservada
	candado.Unlock()

	nuevoRegistro := new(Estructuras.Registro)
	*nuevoRegistro = Estructuras.Registro{
		time.Now(),
		idReservada,
		in.GetTipoLocal(),
		in.GetNombreProducto(),
		in.GetValor(),
		in.GetOrigen(),
		in.GetDestino(),
		int32(test),
	}
	Estructuras.Tabla[idReservada] = nuevoRegistro

	pack := Estructuras.Paquete{
		idReservada,
		uint32(test),
		in.GetTipoLocal(),
		in.GetValor(),
		0,
		2,
	}
	Estructuras.Paquetes[idReservada] = &pack

	switch x := in.GetTipoLocal(); x {
	case "retail":
		candado.Lock()
		Estructuras.ColaRetail = append(Estructuras.ColaRetail, pack)
		candado.Unlock()
	case "prioritario":
		candado.Lock()
		Estructuras.ColaPrioridad = append(Estructuras.ColaPrioridad, pack)
		candado.Unlock()
	case "normal":
		candado.Lock()
		Estructuras.ColaNormal = append(Estructuras.ColaNormal, pack)
		candado.Unlock()
	}

	return &Producto{ID: test}, nil

}
func (s *ServerCliente) EstadoEncargo(ctx context.Context, in *Producto) (*Estatus, error) {
	log.Printf("Solicitado estado de %d", in.ID)

	if valor, existencia := Estructuras.Paquetes[Estructuras.SeguimientoAId[in.ID]]; existencia {
		return &Estatus{Valor: int32(valor.Estado)}, nil
	}
	return &Estatus{Valor: 4}, nil //el extraño caso en que no exista el paquete
}

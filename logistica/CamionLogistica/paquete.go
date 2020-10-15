package CamionLogistica

import (
	"log"

	"sync"

	"../Estructuras"
	"golang.org/x/net/context"
)

type ServerCamion struct {
	placeholder int
}

func (c *ServerCamion) PedirPaquete(ctx context.Context, in *Tipo) (*Paquete, error) {
	var candado sync.Mutex
	var a string
	if in.GetClase() == 1 {
		a = "Retail"
	} else {
		a = "Normal"
	}
	log.Printf("Pedido de paquete de tipo %s", a)
	elem := new(Estructuras.Paquete)

	if x := in.GetClase(); x == (1) {
		if len(Estructuras.ColaRetail) > 0 {
			candado.Lock()
			*elem = Estructuras.ColaRetail[0]
			Estructuras.ColaRetail = Estructuras.ColaRetail[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
			candado.Unlock()
		} else if len(Estructuras.ColaPrioridad) > 0 {
			candado.Lock()
			*elem = Estructuras.ColaPrioridad[0]
			Estructuras.ColaPrioridad = Estructuras.ColaPrioridad[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
			candado.Unlock()
		}
	} else if len(Estructuras.ColaNormal) > 0 {
		candado.Lock()
		*elem = Estructuras.ColaNormal[0]
		Estructuras.ColaNormal = Estructuras.ColaNormal[1:]
		Estructuras.Paquetes[elem.IDPaquete].Estado = 1
		candado.Unlock()
	}

	if item, existe := Estructuras.Tabla[elem.IDPaquete]; existe {
		log.Printf("Entregado paquete id: %d", elem.IDPaquete)
		return &Paquete{
			IDPaquete:   elem.IDPaquete,
			Seguimiento: elem.Seguimiento,
			Tipo:        elem.Tipo,
			Valor:       elem.Valor,
			Intentos:    elem.Intentos,
			Estado:      elem.Estado,
			Origen:      item.Origen,
			Destino:     item.Destino,
		}, nil
	} else {
		return &Paquete{
			IDPaquete:   0,
			Seguimiento: 0,
			Tipo:        "",
			Valor:       0,
			Intentos:    0,
			Estado:      0,
			Origen:      "",
			Destino:     "",
		}, nil
	}

}

func (c *ServerCamion) DevolverPaquete(ctx context.Context, in *Paquete) (*Paquete, error) {
	Estructuras.Paquetes[in.IDPaquete].Intentos = in.GetIntentos()
	Estructuras.Paquetes[in.IDPaquete].Estado = in.GetEstado()
	return in, nil
}

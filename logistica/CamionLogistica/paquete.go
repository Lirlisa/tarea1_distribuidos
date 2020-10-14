package CamionLogistica

import (
	"golang.org/x/net/context"
	"log"
	"../Estructuras"
)

type ServerCamion struct {
	placeholder int
}

func (c *ServerCamion) PedirPaquete(ctx context.Context, in *Tipo) (*Paquete, error) {
	log.Printf("Pedido de paquete de tipo %s", in.getTipo()==1?"Retail":"Normal")
	elem := new(Estructuras.Paquete)

	if x := in.GetClase(); x == (1) {
		if len(Estructuras.ColaRetail) > 0 {
			*elem = Estructuras.ColaRetail[0]
			Estructuras.ColaRetail = Estructuras.ColaRetail[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
		} else if len(Estructuras.ColaPrioridad) > 0 {
			*elem = Estructuras.ColaPrioridad[0]
			Estructuras.ColaPrioridad = Estructuras.ColaPrioridad[1:]
			Estructuras.Paquetes[elem.IDPaquete].Estado = 1
		}
	} else if len(Estructuras.ColaNormal) > 0 {
		*elem = Estructuras.ColaNormal[0]
		Estructuras.ColaNormal = Estructuras.ColaNormal[1:]
		Estructuras.Paquetes[elem.IDPaquete].Estado = 1
	}

	return &Paquete{
		IDPaquete:   elem.IDPaquete,
		Seguimiento: elem.Seguimiento,
		Tipo:        elem.Tipo,
		Valor:       elem.Valor,
		Intentos:    elem.Intentos,
		Estado:      elem.Estado,
	}, nil
}

func (c *ServerCamion) DevolverPaquete(ctx context.Context, in *Paquete) (*Paquete, error) {
	Estructuras.Paquetes[in.IDPaquete].Intentos = in.GetIntentos()
	Estructuras.Paquetes[in.IDPaquete].Estado = in.GetEstado()
	return in, nil
}

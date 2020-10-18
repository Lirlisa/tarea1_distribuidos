package ClienteLogistica

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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
	if in.TipoLocal == "gg" {
		paqueteFinal()
		defer Estructuras.GrpcServerCliente.GracefulStop()
		return &Producto{ID: 0}, nil
	}
	content := fmt.Sprintf("%s,%s,%d,%s,%s\n", in.GetTipoLocal(), in.GetNombreProducto(), in.GetValor(), in.GetOrigen(), in.GetDestino())
	var candado sync.Mutex
	candado.Lock()
	f, err := os.OpenFile("bitacora.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Ha ocurrido un error con el archivo: %s", err)
	}

	if _, err = f.WriteString(content); err != nil {
		f.Close()
		log.Panicf("Error al escribir en archivo: %s", err)
	}
	f.Close()
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
	candado.Lock()
	Estructuras.Tabla[idReservada] = nuevoRegistro
	candado.Unlock()
	pack := Estructuras.Paquete{
		idReservada,
		uint32(test),
		in.GetTipoLocal(),
		in.GetValor(),
		0,
		0,
	}
	candado.Lock()
	Estructuras.Paquetes[idReservada] = &pack
	candado.Unlock()
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
	var candado sync.Mutex
	candado.Lock()
	if valor, existencia := Estructuras.Paquetes[Estructuras.SeguimientoAId[in.ID]]; existencia {
		candado.Unlock()
		return &Estatus{Valor: int32(valor.Estado)}, nil
	}
	candado.Unlock()
	return &Estatus{Valor: 4}, nil //el extraño caso en que no exista el paquete
}

func paqueteFinal() {
	var candado sync.Mutex
	paquete := new(Estructuras.Paquete)
	paquete.Tipo = "gg"
	registro := new(Estructuras.Registro)
	registro.Tipo = "gg"
	candado.Lock()
	Estructuras.Paquetes[paquete.IDPaquete] = paquete
	Estructuras.Tabla[registro.Id] = registro
	for i := 0; i < 3; i++ {
		Estructuras.ColaRetail = append(Estructuras.ColaRetail, *paquete)
		Estructuras.ColaPrioridad = append(Estructuras.ColaPrioridad, *paquete)
		Estructuras.ColaNormal = append(Estructuras.ColaNormal, *paquete)
	}
	candado.Unlock()

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

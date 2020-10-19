package Estructuras

import (
	"time"

	"google.golang.org/grpc"
)

//dato que permite abstraer a los paquetes como dato local
type Paquete struct {
	IDPaquete   uint32
	Seguimiento uint32
	Tipo        string
	Valor       uint32
	Intentos    uint32
	Estado      uint32
}

//dato que permite abstraer los registros
type Registro struct {
	Timestamp   time.Time
	Id          uint32
	Tipo        string
	Nombre      string
	Valor       uint32
	Origen      string
	Destino     string
	Seguimiento int32
}

//inicialización de estructuras globales para compartir información entre los distintos paquetes

var Tabla map[uint32]*Registro = make(map[uint32]*Registro)    //almacacena los registros
var SeguimientoAId map[uint32]uint32 = make(map[uint32]uint32) //permite mapear los valores del seguimiento a las id

var Paquetes map[uint32]*Paquete = make(map[uint32]*Paquete) //mantiene los paquetes

//las colas de los pedidos según tipo
var ColaRetail []Paquete = make([]Paquete, 0, 10)
var ColaPrioridad []Paquete = make([]Paquete, 0, 10)
var ColaNormal []Paquete = make([]Paquete, 0, 10)

//almacena los listener de cada servidor
var GrpcServerCliente *grpc.Server
var GrpcServerCamion *grpc.Server

func main() {
	Paquetes[0] = new(Paquete)
}

package Estructuras

import (
	"time"

	"google.golang.org/grpc"
)

type Paquete struct {
	IDPaquete   uint32
	Seguimiento uint32
	Tipo        string
	Valor       uint32
	Intentos    uint32
	Estado      uint32
}

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

var Tabla map[uint32]*Registro = make(map[uint32]*Registro)
var SeguimientoAId map[uint32]uint32 = make(map[uint32]uint32)

var Paquetes map[uint32]*Paquete = make(map[uint32]*Paquete)
var ColaRetail []Paquete = make([]Paquete, 0, 10)
var ColaPrioridad []Paquete = make([]Paquete, 0, 10)
var ColaNormal []Paquete = make([]Paquete, 0, 10)

var GrpcServerCliente *grpc.Server
var GrpcServerCamion *grpc.Server

func main() {
	Paquetes[0] = new(Paquete)
}

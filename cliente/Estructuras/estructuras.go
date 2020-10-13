package Estructuras

type Paquete struct {
	IDPaquete   uint32
	Seguimiento uint32
	Tipo        string
	Valor       uint32
	Intentos    uint32
	Estado      uint32
}

var Paquetes map[uint32]*Paquete = make(map[uint32]*Paquete)
var ColaRetail []Paquete = make([]Paquete, 0, 10)
var ColaPrioridad []Paquete = make([]Paquete, 0, 10)
var ColaNormal []Paquete = make([]Paquete, 0, 10)

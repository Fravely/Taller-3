package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// Nodo de la lista doblemente enlazada
type Nodo struct {
	clave, valor int
	prev, next   *Nodo
}

// LRU Cache
type LRU struct {
	cap        int
	mapa       map[int]*Nodo
	head, tail *Nodo
}

// NewLRU crea una nueva caché LRU con capacidad dada - O(1)
func NewLRU(cap int) *LRU {
	head := &Nodo{}
	tail := &Nodo{}
	head.next = tail
	tail.prev = head
	return &LRU{
		cap:  cap,
		mapa: make(map[int]*Nodo),
		head: head,
		tail: tail,
	}
}

// moverAlFrente mueve un nodo existente al frente (más reciente) - O(1)
func (l *LRU) moverAlFrente(n *Nodo) {
	// Desconectar nodo
	n.prev.next = n.next
	n.next.prev = n.prev
	// Insertar después del head
	n.next = l.head.next
	n.prev = l.head
	l.head.next.prev = n
	l.head.next = n
}

// insertarAlFrente inserta un nodo nuevo después del head - O(1)
func (l *LRU) insertarAlFrente(n *Nodo) {
	n.next = l.head.next
	n.prev = l.head
	l.head.next.prev = n
	l.head.next = n
}

// eliminarCola elimina el nodo menos reciente (antes del tail) - O(1)
func (l *LRU) eliminarCola() *Nodo {
	lru := l.tail.prev
	lru.prev.next = l.tail
	l.tail.prev = lru.prev
	return lru
}

// Get obtiene un valor de la caché - O(1)
func (l *LRU) Get(clave int) (int, bool) {
	if n, ok := l.mapa[clave]; ok {
		l.moverAlFrente(n)
		return n.valor, true
	}
	return 0, false
}

// Put inserta o actualiza un valor en la caché - O(1)
func (l *LRU) Put(clave, valor int) {
	if n, ok := l.mapa[clave]; ok {
		n.valor = valor
		l.moverAlFrente(n)
		return
	}
	n := &Nodo{clave: clave, valor: valor}
	l.mapa[clave] = n
	l.insertarAlFrente(n)
	if len(l.mapa) > l.cap {
		eliminado := l.eliminarCola()
		delete(l.mapa, eliminado.clave)
	}
}

// Registro representa una fila del CSV
type Registro struct {
	movieID   int
	timestamp int64
}

// CargarSecuencia lee ratings.csv y devuelve movieIDs ordenados por timestamp - O(n log n)
func CargarSecuencia(ruta string) ([]int, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	// Saltar encabezado
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	var registros []Registro
	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		mid, err := strconv.Atoi(row[1])
		if err != nil {
			continue
		}
		ts, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			continue
		}
		registros = append(registros, Registro{movieID: mid, timestamp: ts})
	}

	// Ordenar por timestamp
	sort.Slice(registros, func(i, j int) bool {
		return registros[i].timestamp < registros[j].timestamp
	})

	secuencia := make([]int, len(registros))
	for i, reg := range registros {
		secuencia[i] = reg.movieID
	}
	return secuencia, nil
}

// SimularLRU simula la secuencia y calcula el hit ratio
// SimularLRU simula la secuencia y calcula el hit ratio
func SimularLRU(secuencia []int, capacidad int) float64 {
	if len(secuencia) == 0 {
		return 0
	}
	cache := NewLRU(capacidad)
	hits := 0
	for _, id := range secuencia {
		if _, ok := cache.Get(id); ok {
			hits++
		} else {
			cache.Put(id, id)
		}
	}
	return float64(hits) / float64(len(secuencia))
}

func main() {
	secuencia, err := CargarSecuencia(`C:\Users\Daniel\Downloads\ml-latest-small\ml-latest-small\ratings.csv`)
	if err != nil {
		fmt.Println("Error al cargar el dataset:", err)
		return
	}

	fmt.Printf("Total de accesos cargados: %d\n\n", len(secuencia))

	capacidades := []int{50, 100, 500, 1000}
	fmt.Printf("%-15s %-10s\n", "Tamaño caché", "Hit ratio")
	fmt.Println("---------------------------")
	for _, cap := range capacidades {
		ratio := SimularLRU(secuencia, cap)
		fmt.Printf("%-15d %.4f\n", cap, ratio)
	}

}

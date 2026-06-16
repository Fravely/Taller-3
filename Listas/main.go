package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// ---------- LISTA DOBLEMENTE ENLAZADA + CACHE LRU ----------
type Nodo struct {
	Clave, Valor int
	Prev, Next   *Nodo
}

type LRU struct {
	Cap        int
	Mapa       map[int]*Nodo
	Head, Tail *Nodo
}

func NewLRU(cap int) *LRU {
	head := &Nodo{}
	tail := &Nodo{}
	head.Next = tail
	tail.Prev = head

	return &LRU{
		Cap:  cap,
		Mapa: make(map[int]*Nodo),
		Head: head,
		Tail: tail,
	}
}

func (l *LRU) moverAlFrente(n *Nodo) {
	n.Prev.Next = n.Next
	n.Next.Prev = n.Prev
	l.insertarAlFrente(n)
}

func (l *LRU) insertarAlFrente(n *Nodo) {
	n.Next = l.Head.Next
	n.Prev = l.Head
	l.Head.Next.Prev = n
	l.Head.Next = n
}

func (l *LRU) eliminarCola() *Nodo {
	lru := l.Tail.Prev
	lru.Prev.Next = l.Tail
	l.Tail.Prev = lru.Prev
	return lru
}

func (l *LRU) Get(clave int) (int, bool) {
	if n, ok := l.Mapa[clave]; ok {
		l.moverAlFrente(n)
		return n.Valor, true
	}
	return 0, false
}

func (l *LRU) Put(clave, valor int) {
	if n, ok := l.Mapa[clave]; ok {
		n.Valor = valor
		l.moverAlFrente(n)
		return
	}

	n := &Nodo{Clave: clave, Valor: valor}
	l.Mapa[clave] = n
	l.insertarAlFrente(n)

	if len(l.Mapa) > l.Cap {
		eliminado := l.eliminarCola()
		delete(l.Mapa, eliminado.Clave)
	}
}

// ---------- CARGA DEL DATASET RATINGS.CSV ----------
type Registro struct {
	movieID   int
	timestamp int64
}

func CargarSecuencia(ruta string) ([]int, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	registros := []Registro{}
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if len(row) < 4 {
			continue
		}

		movieID, err := strconv.Atoi(row[1])
		if err != nil {
			continue
		}

		timestamp, err := strconv.ParseInt(row[3], 10, 64)
		if err != nil {
			continue
		}

		registros = append(registros, Registro{
			movieID:   movieID,
			timestamp: timestamp,
		})
	}

	sort.Slice(registros, func(i, j int) bool {
		return registros[i].timestamp < registros[j].timestamp
	})

	secuencia := make([]int, len(registros))
	for i, reg := range registros {
		secuencia[i] = reg.movieID
	}
	return secuencia, nil
}

func SimularLRU(secuencia []int, capacidad int) float64 {
	if len(secuencia) == 0 || capacidad <= 0 {
		return 0
	}

	cache := NewLRU(capacidad)
	hits := 0

	for _, movieID := range secuencia {
		if _, ok := cache.Get(movieID); ok {
			hits++
		} else {
			cache.Put(movieID, movieID)
		}
	}

	return float64(hits) / float64(len(secuencia))
}

// ---------- PROGRAMA PRINCIPAL ----------
func main() {
	rutaCSV := "./ratings.csv"
	if len(os.Args) > 1 {
		rutaCSV = os.Args[1]
	}

	fmt.Printf("Cargando secuencia de accesos desde: %s\n", rutaCSV)
	secuencia, err := CargarSecuencia(rutaCSV)
	if err != nil {
		fmt.Println("Error al cargar el dataset:", err)
		fmt.Println("Asegurate de que ratings.csv este dentro de la carpeta Listas.")
		return
	}

	fmt.Printf("Total de accesos cargados: %d\n\n", len(secuencia))

	capacidades := []int{50, 100, 500, 1000}
	fmt.Printf("%-15s %-10s\n", "Tamano cache", "Hit ratio")
	fmt.Println("---------------------------")

	for _, cap := range capacidades {
		ratio := SimularLRU(secuencia, cap)
		fmt.Printf("%-15d %.4f\n", cap, ratio)
	}
}

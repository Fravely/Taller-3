package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Stack representa una pila genérica implementada sobre un slice de Go.
type Stack[T any] struct {
	elements []T
}

// Push agrega un nuevo elemento al tope de la pila.
// Complejidad Temporal: O(1) amortizado.
func (s *Stack[T]) Push(v T) {
	s.elements = append(s.elements, v)
}

// IsEmpty retorna true si la pila no contiene elementos, false en caso contrario.
// Complejidad Temporal: O(1).
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}

// Pop remueve y retorna el elemento en el tope de la pila.
// Retorna el valor cero del tipo y false si la pila estaba vacía.
// Complejidad Temporal: O(1).
func (s *Stack[T]) Pop() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	index := len(s.elements) - 1
	element := s.elements[index]
	s.elements = s.elements[:index]
	return element, true
}

// Peek retorna el elemento del tope de la pila sin removerlo.
// Retorna false si la pila está vacía.
// Complejidad Temporal: O(1).
func (s *Stack[T]) Peek() (T, bool) {
	if s.IsEmpty() {
		var zero T
		return zero, false
	}
	return s.elements[len(s.elements)-1], true
}

// Registro representa un nodo de datos estructurado del CSV de acciones.
type Registro struct {
	Fecha  string
	Precio float64
}

// LeerPrecios abre un archivo CSV/TXT y parsea las columnas 'Date' y 'Close'.
// Complejidad Temporal: O(n) donde n es el número de líneas del archivo.
// Complejidad Espacial: O(n) para almacenar los registros en memoria.
func LeerPrecios(ruta string) ([]Registro, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil || len(records) < 2 {
		return nil, fmt.Errorf("archivo inválido o vacío")
	}

	dateIdx, closeIdx := 0, 4
	for i, col := range records[0] {
		colClean := strings.TrimSpace(strings.ToLower(col))
		if colClean == "date" {
			dateIdx = i
		} else if colClean == "close" {
			closeIdx = i
		}
	}

	var registros []Registro
	for _, record := range records[1:] {
		if len(record) <= dateIdx || len(record) <= closeIdx {
			continue
		}
		if price, err := strconv.ParseFloat(record[closeIdx], 64); err == nil {
			registros = append(registros, Registro{Fecha: record[dateIdx], Precio: price})
		}
	}
	return registros, nil
}

// CalcularStockSpan calcula el intervalo correlativo de días inferiores o iguales.
// Complejidad Temporal: O(n) lineal gracias a la propiedad de la pila monótona.
// Complejidad Espacial: O(n) para almacenar los resultados del span y la pila.
func CalcularStockSpan(precios []float64) []int {
	n := len(precios)
	span := make([]int, n)
	s := Stack[int]{}

	for i := 0; i < n; i++ {
		for !s.IsEmpty() {
			topIdx, _ := s.Peek()
			if precios[topIdx] <= precios[i] {
				s.Pop()
			} else {
				break
			}
		}

		if s.IsEmpty() {
			span[i] = i + 1
		} else {
			topIdx, _ := s.Peek()
			span[i] = i - topIdx
		}
		s.Push(i)
	}
	return span
}

func main() {
	rutaArchivo := "Stocks/abcd.us.txt"
	if len(os.Args) > 1 {
		rutaArchivo = os.Args[1]
	}

	registros, err := LeerPrecios(rutaArchivo)
	if err != nil || len(registros) == 0 {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	precios := make([]float64, len(registros))
	for i, r := range registros {
		precios[i] = r.Precio
	}

	spans := CalcularStockSpan(precios)
	N := 5

	fmt.Printf("\n=== RESULTADOS %s ===\n", rutaArchivo)
	fmt.Println("Fecha      |Cierre  |Span")
	for i := 0; i < N; i++ {
		fmt.Printf("%s | %.4f | %d\n", registros[i].Fecha, registros[i].Precio, spans[i])
	}
	fmt.Println("========================")
	for i := len(registros) - N; i < len(registros); i++ {
		fmt.Printf("%s | %.4f | %d\n", registros[i].Fecha, registros[i].Precio, spans[i])
	}

	maxIdx := 0
	for i, v := range spans {
		if v > spans[maxIdx] {
			maxIdx = i
		}
	}
	fmt.Printf("\n[MAX] Fecha: %s | Precio: %.4f | Span: %d días\n",
		registros[maxIdx].Fecha, registros[maxIdx].Precio, spans[maxIdx])
}

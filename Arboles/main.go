package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Pelicula struct {
	ID     int
	Titulo string
	Genero string
	año    int
}

type NodoAVL struct {
	Clave    int
	Datos    []Pelicula
	Altura   int
	Izq, Der *NodoAVL
}

type AVL struct {
	Raiz *NodoAVL
	Est  Estadisticas
}

type Estadisticas struct {
	Nodos     int
	Registros int
	RotLL     int
	RotRR     int
	RotLR     int
	RotRL     int
}

func NewAVL() *AVL {
	return &AVL{}
}

func (a *AVL) Insertar(clave int, dato Pelicula) {
	a.Raiz = insertar(a.Raiz, clave, dato, &a.Est)
}

func (a *AVL) ConsultaRango(inicio, fin int) []Pelicula {
	resultado := []Pelicula{}
	consultaRango(a.Raiz, inicio, fin, &resultado)
	return resultado
}

func Altura(n *NodoAVL) int {
	if n == nil {
		return 0
	}
	return n.Altura
}

func actualizarAltura(n *NodoAVL) {
	n.Altura = 1 + max(Altura(n.Izq), Altura(n.Der))
}

func factorBalance(n *NodoAVL) int {
	if n == nil {
		return 0
	}
	return Altura(n.Izq) - Altura(n.Der)
}

func rotarDer(y *NodoAVL) *NodoAVL {
	x := y.Izq
	t2 := x.Der

	x.Der = y
	y.Izq = t2

	actualizarAltura(y)
	actualizarAltura(x)

	return x
}

func rotarIzq(x *NodoAVL) *NodoAVL {
	y := x.Der
	t2 := y.Izq

	y.Izq = x
	x.Der = t2

	actualizarAltura(x)
	actualizarAltura(y)

	return y
}

func insertar(raiz *NodoAVL, clave int, dato Pelicula, est *Estadisticas) *NodoAVL {
	if raiz == nil {
		est.Nodos++
		est.Registros++
		return &NodoAVL{
			Clave:  clave,
			Datos:  []Pelicula{dato},
			Altura: 1,
		}
	}

	if clave < raiz.Clave {
		raiz.Izq = insertar(raiz.Izq, clave, dato, est)
	} else if clave > raiz.Clave {
		raiz.Der = insertar(raiz.Der, clave, dato, est)
	} else {
		raiz.Datos = append(raiz.Datos, dato)
		est.Registros++
		return raiz
	}

	actualizarAltura(raiz)
	balance := factorBalance(raiz)

	// Caso LL: el desbalance esta en izquierda-izquierda.
	if balance > 1 && clave < raiz.Izq.Clave {
		est.RotLL++
		return rotarDer(raiz)
	}

	// Caso RR: el desbalance esta en derecha-derecha.
	if balance < -1 && clave > raiz.Der.Clave {
		est.RotRR++
		return rotarIzq(raiz)
	}

	// Caso LR: primero rotacion izquierda en el hijo, luego derecha.
	if balance > 1 && clave > raiz.Izq.Clave {
		est.RotLR++
		raiz.Izq = rotarIzq(raiz.Izq)
		return rotarDer(raiz)
	}

	// Caso RL: primero rotacion derecha en el hijo, luego izquierda.
	if balance < -1 && clave < raiz.Der.Clave {
		est.RotRL++
		raiz.Der = rotarDer(raiz.Der)
		return rotarIzq(raiz)
	}

	return raiz
}

func consultaRango(n *NodoAVL, inicio, fin int, resultado *[]Pelicula) {
	if n == nil {
		return
	}

	if n.Clave > inicio {
		consultaRango(n.Izq, inicio, fin, resultado)
	}

	if n.Clave >= inicio && n.Clave <= fin {
		*resultado = append(*resultado, n.Datos...)
	}

	if n.Clave < fin {
		consultaRango(n.Der, inicio, fin, resultado)
	}
}

func extraeraño(titulo string) int {
	inicio := strings.LastIndex(titulo, "(")
	fin := strings.LastIndex(titulo, ")")

	if inicio == -1 || fin == -1 || fin <= inicio+1 {
		return 0
	}

	añoStr := titulo[inicio+1 : fin]
	año, err := strconv.Atoi(añoStr)
	if err != nil {
		return 0
	}

	return año
}

func leerPeliculas(ruta string) ([]Pelicula, error) {
	archivo, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer archivo.Close()

	reader := csv.NewReader(archivo)
	reader.FieldsPerRecord = -1

	cabecera, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer la cabecera: %w", err)
	}

	columnas := map[string]int{}
	for i, nombre := range cabecera {
		columnas[strings.TrimSpace(nombre)] = i
	}

	idCol, okID := columnas["movieId"]
	tituloCol, okTitulo := columnas["title"]
	generoCol, okGenero := columnas["genres"]
	if !okID || !okTitulo || !okGenero {
		return nil, fmt.Errorf("el CSV debe tener las columnas movieId,title,genres")
	}

	peliculas := []Pelicula{}
	for {
		fila, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(fila) <= generoCol {
			continue
		}

		id, err := strconv.Atoi(fila[idCol])
		if err != nil {
			continue
		}

		año := extraeraño(fila[tituloCol])
		if año == 0 {
			continue
		}

		peliculas = append(peliculas, Pelicula{
			ID:     id,
			Titulo: fila[tituloCol],
			Genero: fila[generoCol],
			año:    año,
		})
	}

	return peliculas, nil
}

func construirAVL(peliculas []Pelicula) *AVL {
	arbol := NewAVL()
	for _, pelicula := range peliculas {
		arbol.Insertar(pelicula.año, pelicula)
	}
	return arbol
}

func imprimirMetricas(nombre string, arbol *AVL) {
	totalRotaciones := arbol.Est.RotLL + arbol.Est.RotRR + arbol.Est.RotLR + arbol.Est.RotRL

	fmt.Println(nombre)
	fmt.Println("Registros indexados:", arbol.Est.Registros)
	fmt.Println("Nodos AVL (años distintos):", arbol.Est.Nodos)
	fmt.Println("Altura del arbol:", Altura(arbol.Raiz))
	fmt.Printf("Rotaciones LL=%d RR=%d LR=%d RL=%d Total=%d\n",
		arbol.Est.RotLL,
		arbol.Est.RotRR,
		arbol.Est.RotLR,
		arbol.Est.RotRL,
		totalRotaciones,
	)
}

func imprimirResultados(resultados []Pelicula, limite int) {
	if limite <= 0 || limite > len(resultados) {
		limite = len(resultados)
	}

	fmt.Println("movieId | año | titulo | generos")
	for i := 0; i < limite; i++ {
		p := resultados[i]
		fmt.Printf("%d | %d | %s | %s\n", p.ID, p.año, p.Titulo, p.Genero)
	}

	if len(resultados) > limite {
		fmt.Println("... resultados restantes:", len(resultados)-limite)
	}
}

func main() {
	ruta := flag.String("movies", "movies.csv", "ruta del archivo movies.csv")
	inicio := flag.Int("from", 1995, "inicio del rango")
	fin := flag.Int("to", 2000, "fin del rango")
	limite := flag.Int("limit", 20, "cantidad maxima de resultados a mostrar")
	flag.Parse()

	if *inicio > *fin {
		fmt.Println("Error: from no puede ser mayor que to")
		return
	}

	peliculas, err := leerPeliculas(*ruta)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	arbol := construirAVL(peliculas)
	imprimirMetricas("Métricas del árbol:", arbol)

	resultados := arbol.ConsultaRango(*inicio, *fin)
	fmt.Printf("\nPeliculas entre %d y %d: %d\n", *inicio, *fin, len(resultados))
	imprimirResultados(resultados, *limite)
}

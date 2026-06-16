package main

import "testing"

func TestExtraeraño(t *testing.T) {
	año := extraeraño("Toy Story (1995)")
	if año != 1995 {
		t.Fatalf("se esperaba 1995, se obtuvo %d", año)
	}
}

func TestConsultaRangoOrdenAscendente(t *testing.T) {
	peliculas := []Pelicula{
		{ID: 1, Titulo: "A (2000)", año: 2000},
		{ID: 2, Titulo: "B (1990)", año: 1990},
		{ID: 3, Titulo: "C (2010)", año: 2010},
		{ID: 4, Titulo: "D (2000)", año: 2000},
	}

	arbol := construirAVL(peliculas)
	resultado := arbol.ConsultaRango(1995, 2005)

	if arbol.Est.Registros != 4 {
		t.Fatalf("registros esperados 4, obtenidos %d", arbol.Est.Registros)
	}
	if arbol.Est.Nodos != 3 {
		t.Fatalf("nodos esperados 3, obtenidos %d", arbol.Est.Nodos)
	}
	if len(resultado) != 2 {
		t.Fatalf("resultados esperados 2, obtenidos %d", len(resultado))
	}
	for _, pelicula := range resultado {
		if pelicula.año != 2000 {
			t.Fatalf("resultado fuera de rango: %+v", pelicula)
		}
	}
}

func TestAVLSeBalanceaConInsercionesOrdenadas(t *testing.T) {
	arbol := NewAVL()

	for i := 1; i <= 100; i++ {
		arbol.Insertar(i, Pelicula{ID: i, año: i})
	}

	if Altura(arbol.Raiz) > 8 {
		t.Fatalf("altura demasiado grande para AVL con 100 claves ordenadas: %d", Altura(arbol.Raiz))
	}
	if arbol.Est.RotRR == 0 {
		t.Fatalf("se esperaban rotaciones RR con inserciones ascendentes")
	}
}

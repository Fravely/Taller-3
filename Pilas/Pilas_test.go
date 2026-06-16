package main

import (
	"testing"
)

// 5.5 TEST: Caso Límite y Normal de la Estructura Stack
func TestStack_CasosLimite(t *testing.T) {
	s := Stack[int]{}

	// Caso límite específico: Pop en una pila vacía
	_, ok := s.Pop()
	if ok {
		t.Error("Se esperaba ok=false al hacer Pop en una pila vacía")
	}

	// Caso límite específico: Peek en una pila vacía
	_, ok = s.Peek()
	if ok {
		t.Error("Se esperaba ok=false al hacer Peek en una pila vacía")
	}

	// Caso Normal de la pila
	s.Push(100)
	if s.IsEmpty() {
		t.Error("La pila no debería estar vacía tras un Push")
	}

	val, ok := s.Peek()
	if !ok || val != 100 {
		t.Errorf("Peek falló. Obtenido: %d", val)
	}
}

// 5.5 TEST: Caso de negocio normal del Stock Span
func TestCalcularStockSpan(t *testing.T) {
	// Datos de prueba basados en el ejemplo cronológico analizado
	precios := []float64{3.65, 4.09, 4.80, 4.76, 4.51}
	esperado := []int{1, 2, 3, 1, 1}

	// Llamamos a la función una sola vez
	resultado := CalcularStockSpan(precios)

	// Validamos las longitudes del arreglo devuelto
	if len(resultado) != len(esperado) {
		t.Fatalf("Longitudes diferentes. Esperado: %d, Obtenido: %d", len(esperado), len(resultado))
	}

	// Validamos cada valor del span elemento por elemento
	for i, v := range resultado {
		if v != esperado[i] {
			t.Errorf("Error en índice %d. Esperado: %d, Obtenido: %d", i, esperado[i], v)
		}
	}
}

// 5.5 TEST: Caso de Error en Lectura de Archivos
func TestLeerPrecios_Error(t *testing.T) {
	_, err := LeerPrecios("ruta/inexistente/archivo_fantasma.txt")
	if err == nil {
		t.Error("Se esperaba un error al ingresar una ruta de archivo falsa")
	}
}

// ============================================================================
// 5.4 BENCHMARKS: Colección de datos empíricos para el Informe de Performance
// ============================================================================

func ejecutarBenchmark(b *testing.B, n int) {
	precios := make([]float64, n)
	for i := 0; i < n; i++ {
		precios[i] = float64(i % 100) // Rampa repetitiva controlada
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalcularStockSpan(precios)
	}
}

func BenchmarkCalcularStockSpan_1K(b *testing.B)   { ejecutarBenchmark(b, 1000) }
func BenchmarkCalcularStockSpan_10K(b *testing.B)  { ejecutarBenchmark(b, 10000) }
func BenchmarkCalcularStockSpan_100K(b *testing.B) { ejecutarBenchmark(b, 100000) }

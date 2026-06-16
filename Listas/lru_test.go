package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Test caso normal: Get y Put básicos
func TestLRU_Normal(t *testing.T) {
	cache := NewLRU(3)
	cache.Put(1, 10)
	cache.Put(2, 20)
	cache.Put(3, 30)

	if val, ok := cache.Get(1); !ok || val != 10 {
		t.Errorf("esperado 10, obtenido %d", val)
	}
}

// Test caso límite: caché con capacidad 1
func TestLRU_CapacidadUno(t *testing.T) {
	cache := NewLRU(1)
	cache.Put(1, 10)
	cache.Put(2, 20)

	if _, ok := cache.Get(1); ok {
		t.Error("el nodo 1 debería haber sido expulsado")
	}
	if val, ok := cache.Get(2); !ok || val != 20 {
		t.Errorf("esperado 20, obtenido %d", val)
	}
}

// Test expulsión correcta del LRU
func TestLRU_ExpulsionCorrecta(t *testing.T) {
	cache := NewLRU(2)
	cache.Put(1, 10)
	cache.Put(2, 20)
	cache.Get(1)
	cache.Put(3, 30)

	if _, ok := cache.Get(2); ok {
		t.Error("el nodo 2 debería haber sido expulsado")
	}
	if _, ok := cache.Get(1); !ok {
		t.Error("el nodo 1 debería seguir en caché")
	}
	if _, ok := cache.Get(3); !ok {
		t.Error("el nodo 3 debería estar en caché")
	}
}

// Test caso límite: Get en caché vacía
func TestLRU_CacheVacia(t *testing.T) {
	cache := NewLRU(5)
	if _, ok := cache.Get(99); ok {
		t.Error("no debería encontrar nada en caché vacía")
	}
}

// Test actualizar valor existente
func TestLRU_Actualizar(t *testing.T) {
	cache := NewLRU(2)
	cache.Put(1, 10)
	cache.Put(1, 99)
	if val, ok := cache.Get(1); !ok || val != 99 {
		t.Errorf("esperado 99, obtenido %d", val)
	}
}

// Test CargarSecuencia con ruta inválida
func TestCargarSecuencia_RutaInvalida(t *testing.T) {
	_, err := CargarSecuencia("no_existe.csv")
	if err == nil {
		t.Error("debería retornar error con ruta inválida")
	}
}

// Test SimularLRU con secuencia vacía
func TestSimularLRU_Vacia(t *testing.T) {
	ratio := SimularLRU([]int{}, 10)
	if ratio != 0 {
		t.Errorf("esperado 0, obtenido %f", ratio)
	}
}

// Test muchas inserciones para verificar expulsión continua
func TestLRU_MuchasInserciones(t *testing.T) {
	cache := NewLRU(3)
	for i := 0; i < 100; i++ {
		cache.Put(i, i*10)
	}
	if _, ok := cache.Get(97); !ok {
		t.Error("97 debería estar en caché")
	}
	if _, ok := cache.Get(98); !ok {
		t.Error("98 debería estar en caché")
	}
	if _, ok := cache.Get(99); !ok {
		t.Error("99 debería estar en caché")
	}
}

// Test que Put no excede la capacidad
func TestLRU_Capacidad(t *testing.T) {
	cap := 5
	cache := NewLRU(cap)
	for i := 0; i < 20; i++ {
		cache.Put(i, i)
	}
	if len(cache.Mapa) > cap {
		t.Errorf("caché excedió capacidad: tiene %d elementos", len(cache.Mapa))
	}
}

// Benchmark de Put (Prueba de rendimiento requerida por la rúbrica)
func BenchmarkLRU_Put(b *testing.B) {
	cache := NewLRU(1000)
	for i := 0; i < b.N; i++ {
		cache.Put(i%2000, i)
	}
}

// Benchmark de Get (Prueba de rendimiento requerida por la rúbrica)
func BenchmarkLRU_Get(b *testing.B) {
	cache := NewLRU(1000)
	for i := 0; i < 1000; i++ {
		cache.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i % 1000)
	}
}

// Test CargarSecuencia usando un archivo virtual temporal dinámico
func TestCargarSecuencia_Real(t *testing.T) {
	// Creamos un directorio temporal único del sistema operativo
	dir := t.TempDir()
	tmpFile := filepath.Join(dir, "ratings_mock.csv")

	// Escribimos filas falsas simulando MovieLens con timestamps desordenados
	// Esto valida si el sort.Slice de tu lru.go funciona correctamente
	csvContent := []byte("userId,movieId,rating,timestamp\n" +
		"1,10,4.0,1500000000\n" + // Segundo en la línea de tiempo
		"1,20,5.0,1400000000\n" + // Primero en la línea de tiempo
		"1,30,3.5,1600000000\n") // Tercero en la línea de tiempo

	if err := os.WriteFile(tmpFile, csvContent, 0666); err != nil {
		t.Fatalf("error creando archivo temporal de pruebas: %v", err)
	}

	// Probamos tu lector enviando la ruta dinámica
	secuencia, err := CargarSecuencia(tmpFile)
	if err != nil {
		t.Fatalf("error cargando dataset: %v", err)
	}

	if len(secuencia) != 3 {
		t.Errorf("se esperaban 3 elementos parseados, obtenido %d", len(secuencia))
	}

	// Como tu algoritmo ordena por tiempo, la película 20 (timestamp menor) debe ir al inicio
	if secuencia[0] != 20 {
		t.Errorf("se esperaba que el primer elemento ordenado sea la película 20, obtenido %d", secuencia[0])
	}
}

// Test SimularLRU usando un archivo simulado autogenerado
func TestSimularLRU_Real(t *testing.T) {
	dir := t.TempDir()
	tmpFile := filepath.Join(dir, "ratings_mock.csv")

	csvContent := []byte("userId,movieId,rating,timestamp\n" +
		"1,10,4.0,1400000000\n" +
		"1,10,4.5,1450000000\n" + // Repetido de inmediato -> Debe causar Cache Hit
		"1,20,5.0,1500000000\n")

	if err := os.WriteFile(tmpFile, csvContent, 0666); err != nil {
		t.Fatalf("error creando archivo: %v", err)
	}

	secuencia, err := CargarSecuencia(tmpFile)
	if err != nil {
		t.Fatalf("error leyendo archivo: %v", err)
	}

	// Simulamos con una capacidad ajustada
	ratio := SimularLRU(secuencia, 2)
	if ratio < 0 || ratio > 1 {
		t.Errorf("hit ratio fuera de rango porcentual válido: %f", ratio)
	}
}

// Test Get mueve el nodo al frente correctamente
func TestLRU_OrdenUso(t *testing.T) {
	cache := NewLRU(3)
	cache.Put(1, 10)
	cache.Put(2, 20)
	cache.Put(3, 30)
	cache.Get(1)     // El elemento 1 pasa a ser el usado más recientemente
	cache.Put(4, 40) // El elemento menos usado era el 2, por ende se expulsa
	if _, ok := cache.Get(2); ok {
		t.Error("nodo 2 debería haber sido expulsado de la memoria")
	}
	if _, ok := cache.Get(1); !ok {
		t.Error("nodo 1 debería seguir disponible")
	}
}

// Test Put actualiza sin duplicar nodos en el mapa
func TestLRU_PutActualizaSinDuplicar(t *testing.T) {
	cache := NewLRU(3)
	cache.Put(1, 10)
	cache.Put(1, 20)
	cache.Put(1, 30)
	if len(cache.Mapa) != 1 {
		t.Errorf("debería haber solo 1 clave en el mapa, se encontraron %d", len(cache.Mapa))
	}
	if val, ok := cache.Get(1); !ok || val != 30 {
		t.Errorf("esperado el último valor actualizado (30), obtenido %d", val)
	}
}

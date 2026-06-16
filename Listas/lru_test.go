package main

import "testing"

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
	if len(cache.mapa) > cap {
		t.Errorf("caché excedió capacidad: tiene %d elementos", len(cache.mapa))
	}
}

// Benchmark de Put
func BenchmarkLRU_Put(b *testing.B) {
	cache := NewLRU(1000)
	for i := 0; i < b.N; i++ {
		cache.Put(i%2000, i)
	}
}

// Benchmark de Get
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

// Test CargarSecuencia con archivo real
func TestCargarSecuencia_Real(t *testing.T) {
	secuencia, err := CargarSecuencia(`C:\Users\Daniel\Downloads\ml-latest-small\ml-latest-small\ratings.csv`)
	if err != nil {
		t.Fatalf("error cargando dataset: %v", err)
	}
	if len(secuencia) == 0 {
		t.Error("la secuencia no debería estar vacía")
	}
}

// Test SimularLRU con secuencia real
func TestSimularLRU_Real(t *testing.T) {
	secuencia, _ := CargarSecuencia(`C:\Users\Daniel\Downloads\ml-latest-small\ml-latest-small\ratings.csv`)
	ratio := SimularLRU(secuencia, 100)
	if ratio < 0 || ratio > 1 {
		t.Errorf("hit ratio fuera de rango: %f", ratio)
	}
}

// Test Get mueve el nodo al frente correctamente
func TestLRU_OrdenUso(t *testing.T) {
	cache := NewLRU(3)
	cache.Put(1, 10)
	cache.Put(2, 20)
	cache.Put(3, 30)
	cache.Get(1)     // 1 pasa a ser el más reciente
	cache.Put(4, 40) // debe expulsar al 2
	if _, ok := cache.Get(2); ok {
		t.Error("nodo 2 debería haber sido expulsado")
	}
	if _, ok := cache.Get(1); !ok {
		t.Error("nodo 1 debería seguir en caché")
	}
}

// Test Put actualiza sin duplicar nodos
func TestLRU_PutActualizaSinDuplicar(t *testing.T) {
	cache := NewLRU(3)
	cache.Put(1, 10)
	cache.Put(1, 20)
	cache.Put(1, 30)
	if len(cache.mapa) != 1 {
		t.Errorf("debería haber solo 1 elemento, hay %d", len(cache.mapa))
	}
	if val, ok := cache.Get(1); !ok || val != 30 {
		t.Errorf("esperado 30, obtenido %d", val)
	}
}

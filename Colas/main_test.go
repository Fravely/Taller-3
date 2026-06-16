package main

import "testing"

func TestColaFIFO(t *testing.T) {
	var cola Cola

	if cola.Largo() != 0 {
		t.Fatalf("largo inicial esperado 0, obtenido %d", cola.Largo())
	}

	cola.Encolar(10)
	cola.Encolar(20)
	cola.Encolar(30)

	if cola.Largo() != 3 {
		t.Fatalf("largo esperado 3, obtenido %d", cola.Largo())
	}

	frente, ok := cola.Frente()
	if !ok || frente != 10 {
		t.Fatalf("frente esperado 10, obtenido %d, ok=%v", frente, ok)
	}

	valor, ok := cola.Desencolar()
	if !ok || valor != 10 {
		t.Fatalf("primer desencolar esperado 10, obtenido %d, ok=%v", valor, ok)
	}

	valor, ok = cola.Desencolar()
	if !ok || valor != 20 {
		t.Fatalf("segundo desencolar esperado 20, obtenido %d, ok=%v", valor, ok)
	}

	valor, ok = cola.Desencolar()
	if !ok || valor != 30 {
		t.Fatalf("tercer desencolar esperado 30, obtenido %d, ok=%v", valor, ok)
	}

	_, ok = cola.Desencolar()
	if ok {
		t.Fatalf("no se esperaba desencolar desde una cola vacia")
	}
}

func TestParsearLineaApache(t *testing.T) {
	linea := `127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] "GET /index.html HTTP/1.0" 200 2326`

	ip, ts, err := ParsearLinea(linea)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if ip != "127.0.0.1" {
		t.Fatalf("ip esperada 127.0.0.1, obtenida %s", ip)
	}
	if ts <= 0 {
		t.Fatalf("timestamp esperado positivo, obtenido %d", ts)
	}
}

func TestParsearLineaInvalida(t *testing.T) {
	_, _, err := ParsearLinea("linea sin formato apache")
	if err == nil {
		t.Fatalf("se esperaba error para una linea invalida")
	}
}

func TestLimitadorAceptaYRechazaPorVentana(t *testing.T) {
	limitador := NuevoLimitador(2, 60)
	ip := "10.0.0.1"

	if !limitador.Permitir(ip, 100) {
		t.Fatalf("la primera peticion debio aceptarse")
	}
	if !limitador.Permitir(ip, 120) {
		t.Fatalf("la segunda peticion debio aceptarse")
	}
	if limitador.Permitir(ip, 130) {
		t.Fatalf("la tercera peticion dentro de la ventana debio rechazarse")
	}
	if limitador.TotalRechazos() != 1 {
		t.Fatalf("rechazos esperados 1, obtenidos %d", limitador.TotalRechazos())
	}

	if !limitador.Permitir(ip, 161) {
		t.Fatalf("la peticion despues de vencer la primera marca debio aceptarse")
	}
}

func TestTopRechazos(t *testing.T) {
	limitador := NuevoLimitador(1, 60)

	limitador.Permitir("10.0.0.1", 100)
	limitador.Permitir("10.0.0.1", 110)
	limitador.Permitir("10.0.0.1", 120)

	limitador.Permitir("10.0.0.2", 100)
	limitador.Permitir("10.0.0.2", 110)

	top := limitador.TopRechazos(2)
	if len(top) != 2 {
		t.Fatalf("se esperaban 2 IPs en el top, se obtuvo %d", len(top))
	}
	if top[0].IP != "10.0.0.1" || top[0].Cnt != 2 {
		t.Fatalf("primer top incorrecto: %+v", top[0])
	}
	if top[1].IP != "10.0.0.2" || top[1].Cnt != 1 {
		t.Fatalf("segundo top incorrecto: %+v", top[1])
	}
}

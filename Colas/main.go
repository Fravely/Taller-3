package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ---------- COLA FIFO CON LISTA ENLAZADA ----------
type Nodo struct {
	ts  int64
	sig *Nodo
}

type Cola struct {
	cabeza *Nodo
	cola   *Nodo
	largo  int
}

func (c *Cola) Encolar(ts int64) {
	n := &Nodo{ts: ts}
	if c.cola == nil {
		c.cabeza = n
		c.cola = n
	} else {
		c.cola.sig = n
		c.cola = n
	}
	c.largo++
}

func (c *Cola) Desencolar() (int64, bool) {
	if c.cabeza == nil {
		return 0, false
	}
	ts := c.cabeza.ts
	c.cabeza = c.cabeza.sig
	if c.cabeza == nil {
		c.cola = nil
	}
	c.largo--
	return ts, true
}

func (c *Cola) Frente() (int64, bool) {
	if c.cabeza == nil {
		return 0, false
	}
	return c.cabeza.ts, true
}

func (c *Cola) Largo() int {
	return c.largo
}

// ---------- PARSEAR LÍNEA DE LOG (FORMATO APACHE) ----------
func ParsearLinea(linea string) (ip string, ts int64, err error) {
	idx := strings.Index(linea, " ")
	if idx == -1 {
		return "", 0, fmt.Errorf("no se encontró IP")
	}
	ip = linea[:idx]

	ini := strings.Index(linea, "[")
	fin := strings.Index(linea, "]")
	if ini == -1 || fin == -1 {
		return "", 0, fmt.Errorf("no se encontró timestamp")
	}
	fechaStr := linea[ini+1 : fin]

	layout := "02/Jan/2006:15:04:05 -0700"
	t, err := time.Parse(layout, fechaStr)
	if err != nil {
		return "", 0, fmt.Errorf("fecha inválida: %s", fechaStr)
	}
	return ip, t.Unix(), nil
}

// ---------- LIMITADOR DE TASA ----------
type Limitador struct {
	colas    map[string]*Cola
	M        int
	T        int64
	rechazos map[string]int
}

func NuevoLimitador(M int, T int64) *Limitador {
	return &Limitador{
		colas:    make(map[string]*Cola),
		M:        M,
		T:        T,
		rechazos: make(map[string]int),
	}
}

// Permitir decide y además retorna si fue aceptado (true) o rechazado (false)
func (l *Limitador) Permitir(ip string, ahora int64) bool {
	q, ok := l.colas[ip]
	if !ok {
		q = &Cola{}
		l.colas[ip] = q
	}

	// Limpiar vencidas
	for q.Largo() > 0 {
		front, _ := q.Frente()
		if front < ahora-l.T {
			q.Desencolar()
		} else {
			break
		}
	}

	if q.Largo() < l.M {
		q.Encolar(ahora)
		return true
	}
	l.rechazos[ip]++
	return false
}
func (l *Limitador) TotalRechazos() int {
	total := 0
	for _, v := range l.rechazos {
		total += v
	}
	return total
}
func (l *Limitador) TopRechazos(n int) []struct {
	IP  string
	Cnt int
} {
	type par struct {
		ip  string
		cnt int
	}
	lista := make([]par, 0, len(l.rechazos))
	for ip, cnt := range l.rechazos {
		lista = append(lista, par{ip, cnt})
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].cnt > lista[j].cnt })
	res := make([]struct {
		IP  string
		Cnt int
	}, 0, n)
	for i := 0; i < n && i < len(lista); i++ {
		res = append(res, struct {
			IP  string
			Cnt int
		}{lista[i].ip, lista[i].cnt})
	}
	return res
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Uso: %s <archivo_log> <M> <T_segundos>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Ejemplo: %s access.log 5 60\n", os.Args[0])
		os.Exit(1)
	}

	ruta := os.Args[1]
	M, err := strconv.Atoi(os.Args[2])
	if err != nil || M <= 0 {
		fmt.Fprintf(os.Stderr, "M debe ser un entero positivo\n")
		os.Exit(1)
	}
	T, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil || T <= 0 {
		fmt.Fprintf(os.Stderr, "T debe ser un entero positivo (segundos)\n")
		os.Exit(1)
	}

	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error abriendo archivo: %v\n", err)
		os.Exit(1)
	}
	defer archivo.Close()

	limitador := NuevoLimitador(M, T)
	scanner := bufio.NewScanner(archivo)
	lineaNum := 0

	// (Opcional) mensaje de inicio por stderr para no mezclar con la salida normal
	fmt.Fprintf(os.Stderr, "Procesando %s (M=%d, T=%d s)...\n", ruta, M, T)

	for scanner.Scan() {
		lineaNum++
		linea := scanner.Text()
		if linea == "" {
			continue
		}
		ip, ts, err := ParsearLinea(linea)
		if err != nil {
			// Si hay error, no contamos la línea; podemos mostrar el error por stderr si se desea
			// pero para no saturar, solo lo mostramos si es una línea con error real (opcional)
			// Aquí decidimos ignorarla silenciosamente.
			continue
		}
		aceptado := limitador.Permitir(ip, ts)
		// Mostrar resultado por cada petición (formato compacto)
		hora := time.Unix(ts, 0).Format("15:04:05")
		if aceptado {
			fmt.Printf("ACEPTA   %s - %s\n", ip, hora)
		} else {
			fmt.Printf("RECHAZA  %s - %s\n", ip, hora)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error leyendo archivo: %v\n", err)
	}

	// Resumen final
	fmt.Println("\n°°°°°°Resumen de Peticiones Rechazadas°°°°°°")
	fmt.Printf("Total de peticiones rechazadas: %d\n", limitador.TotalRechazos())
	fmt.Println("Top 5 IPs con más rechazos:")
	for i, r := range limitador.TopRechazos(5) {
		fmt.Printf("%d. %s - %d rechazos\n", i+1, r.IP, r.Cnt)
	}
}

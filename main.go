package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

var t_0 float64 = 0.5
var y_0 float64 = 2.5
var y = make(map[float64]float64)

var filename = flag.String("o", "values.csv", "Output File Name")
var deltaT = flag.Float64("t", .001, "Time Delta")
var iterations = flag.Int("i", 1000, "Number of iterations")

func main() {
	// Commandline Parameter
	flag.Parse()

	// Anfangswerte
	y[t_0] = y_0

	for i := 1; i < *iterations; i++ {
		calcY(i)
		saveValues(y)
	}
}

// y'(t)=phi(t, y(t))
func phi(t, y_t float64) float64 {
	// Beispielfunktion phi(t, y(t)) = y/t
	return y_t / t
}

// Speichern der berechneten Werte
func saveValues(v map[float64]float64) {
	var data []byte
	for t, y := range v {
		data = append(data, []byte(fmt.Sprintf("%f, %f\n", t, y))...)
	}
	ioutil.WriteFile(*filename, data, 0644)
}

// Berechnet k_1 bis k_4 für Zeitschritt i
func k(i int) (k1 float64, k2 float64, k3 float64, k4 float64) {
	k1 = phi(t(i), y[t(i)])
	k2 = phi(t(i)+*deltaT/2, y[t(i)]+*deltaT*k1/2)
	k3 = phi(t(i)+*deltaT/2, y[t(i)]+*deltaT*k2/2)
	k3 = phi(t(i)+*deltaT, y[t(i)]+*deltaT*k3)
	return
}

// Berechnet y am i-ten Zeitschritt, wenn y für Zeitschritte < i bereits berechnet ist.
func calcY(i int) {
	k1, k2, k3, k4 := k(i - 1)
	y[t(i)] = y[t(i-1)] + (k1/6+k2/3+k3/3+k4/6)*(*deltaT)
}

// Zeitpunkt nach i Zeitschritten
func t(i int) float64 {
	return t_0 + *deltaT*float64(i)
}

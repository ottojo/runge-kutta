package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

var Sonnensystem []Planet
var deltaT float64 = 3600 * 24 * 1

type Planet struct {
	position map[float64]Vector
	velocity Vector
	mass     float64
	name     string
}

func main() {
	planetdata, err := ioutil.ReadFile("Sonnensystem.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	planets := strings.Split(string(planetdata), "\r\n")
	for _, planet := range planets {
		data := strings.Split(planet, "\t")
		name := data[0]
		mass, _ := strconv.ParseFloat(data[1], 64)
		r_x, _ := strconv.ParseFloat(data[2], 64)
		r_y, _ := strconv.ParseFloat(data[3], 64)
		r_z, _ := strconv.ParseFloat(data[4], 64)
		v_x, _ := strconv.ParseFloat(data[5], 64)
		v_y, _ := strconv.ParseFloat(data[6], 64)
		v_z, _ := strconv.ParseFloat(data[7], 64)

		position := make(map[float64]Vector)
		position[0] = Vector{r_x, r_y, r_z}

		Sonnensystem = append(Sonnensystem, Planet{name: name,
			mass:     mass,
			position: position,
			velocity: Vector{v_x, v_y, v_z}})
	}
	fmt.Println(Sonnensystem[3].position[t(0)])

	i := 1
	var data []byte

	for (i < 50000) {
		for _, p := range Sonnensystem {
			p.calcTimestep(i)
		}
		fmt.Println(Sonnensystem[3].velocity)
		data = append(data, []byte(fmt.Sprintf("%f, %f, %f\n",
			Sonnensystem[3].position[t(i)].X,
			Sonnensystem[3].position[t(i)].Y,
			Sonnensystem[3].position[t(i)].Z))...)

		i++
	}
	ioutil.WriteFile("sonnensystem.csv", data, 0644)
}

func (p *Planet) phi(t float64, r Vector) (velocity Vector) {
	var impulseChange Vector
	for _, otherPlanet := range Sonnensystem {
		if otherPlanet.name != p.name {
			f := F(p.mass, otherPlanet.mass, r, otherPlanet.position[lastValidTime(t)])
			impulseChange = impulseChange.plus(f)
		}
	}
	//a := impulseChange.scale(deltaT / p.mass)
	//fmt.Println(a)
	velocity = impulseChange.scale(deltaT / p.mass)
	return
}

func (p *Planet) k(i int) (k1, k2, k3, k4 Vector) {
	k1 = p.phi(t(i), p.position[t(i)])
	k2 = p.phi(t(i)+deltaT/2, p.position[t(i)].plus(k1.scale(deltaT/2)))
	k3 = p.phi(t(i)+deltaT/2, p.position[t(i)].plus(k2.scale(deltaT/2)))
	k4 = p.phi(t(i)+deltaT/2, p.position[t(i)].plus(k3.scale(deltaT)))
	return
}

func (p *Planet) calcTimestep(i int) {
	//fmt.Printf("calculating step %d for %s\n", i, p.name)
	k1, k2, k3, k4 := p.k(i - 1);
	positionChange := k1.scale(deltaT / 6).
		plus(k2.scale(deltaT / 3)).
		plus(k3.scale(deltaT / 3)).
		plus(k4.scale(deltaT / 6))
	p.velocity = (positionChange)
	p.position[t(i)] = p.position[t(i-1)].plus(positionChange.scale(deltaT))
}

// Kraft von 2 auf 1
func F(m1, m2 float64, r1, r2 Vector) Vector {
	G := 6.674e-20
	return r2.minus(r1).scale(G * m1 * m2 * (1 / math.Pow(r2.minus(r1).length(), 3)))
}

type Vector struct {
	X, Y, Z float64
}

func (v Vector) minus(v2 Vector) Vector {
	return v.plus(v2.scale(-1))
}

func (v Vector) length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) scale(a float64) Vector {
	return Vector{a * v.X, a * v.Y, a * v.Z}
}

func (v Vector) plus(v2 Vector) Vector {
	return Vector{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

// Zeitpunkt nach i Zeitschritten
func t(i int) float64 {
	return deltaT * float64(i)
}

func lastValidTime(t float64) float64 {
	return float64(int(t/deltaT)) * deltaT
}

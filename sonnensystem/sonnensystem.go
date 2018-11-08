package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var Sonnensystem []Planet
var deltaT float64 = 3600 * 24 * 5

var realTimeYears = flag.Float64("realTime", 1, "Time to simulate in years")
var simulationTimeSeconds = flag.Float64("simulationDuration", 30, "Simulation duration in seconds")
var frameRate = flag.Int("fps", 30, "Number of frames per second")
var outputFileName = flag.String("o", "sonnensystem.csv", "Output file name")
var render = flag.Bool("r", false, "Render animation using gnuplot and ffmpeg")
var inputFileName = flag.String("i", "Sonnensystem.dat", "Input file name")

type Planet struct {
	position map[float64]Vector
	velocity map[float64]Vector
	mass     float64
	name     string
}

func main() {
	flag.Parse()
	planetdata, err := ioutil.ReadFile(*inputFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	numberOfFrames := int(float64(*frameRate) * *simulationTimeSeconds)
	deltaT = *realTimeYears * 365 * 24 * 3600 / float64(numberOfFrames)

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

		velocity := make(map[float64]Vector)
		velocity[0] = Vector{v_x, v_y, v_z}

		Sonnensystem = append(Sonnensystem, Planet{name: name,
			mass:     mass,
			position: position,
			velocity: velocity})
	}
	fmt.Println(Sonnensystem[3].position[t(0)])

	i := 1
	var data []byte

	for i < numberOfFrames {
		for _, p := range Sonnensystem {
			p.calcTimestep(i)
		}

		for _, planet := range Sonnensystem {
			data = append(data, []byte(fmt.Sprintf("%f, %f, %f,",
				planet.position[t(i)].X,
				planet.position[t(i)].Y,
				planet.position[t(i)].Z))...)
		}
		data[len(data)-1] = '\n'
		i++
	}
	ioutil.WriteFile(*outputFileName, data, 0644)
	if *render {
		exec.Command("mkdir", "animation").Run()
		gnuplot := exec.Command("gnuplot",
			"-e", "filename='"+*outputFileName+"'",
			"-e", "framenumber='"+strconv.Itoa(numberOfFrames)+"'",
			"plot.gp")
		gnuplot.Stderr = os.Stderr
		gnuplot.Stdout = os.Stdout
		gnuplot.Run()
		ffmpeg := exec.Command("ffmpeg",
			"-f", "image2",
			"-framerate", strconv.Itoa(*frameRate),
			"-i", "animation/%04d.png",
			*outputFileName+".mp4")
		ffmpeg.Stderr = os.Stderr
		ffmpeg.Stdout = os.Stdout
		ffmpeg.Run()
	}
}

func (p *Planet) a(r Vector, newTimeStep int) Vector {
	var vPunktM Vector
	for _, otherPlanet := range Sonnensystem {
		if otherPlanet.name != p.name {
			f := F(p.mass, otherPlanet.mass, r, otherPlanet.position[t(newTimeStep-1)])
			vPunktM = vPunktM.plus(f)
		}
	}
	return vPunktM.scale(1 / p.mass)
}

func (p *Planet) calcTimestep(i int) {

	// Werte am Anfang des Zeitschritts
	r1 := p.position[t(i-1)]
	v1 := p.velocity[t(i-1)]
	a1 := p.a(r1, i)

	// Werte in der Mitte des Zeitschritts, basierend auf Anfangswerten
	r2 := r1.plus(v1.scale(deltaT / 2))
	v2 := v1.plus(a1.scale(deltaT / 2))
	a2 := p.a(r2, i)

	// Werte in der Mitte des Zeitschritts, basierend auf soeben berechneten Werten
	r3 := r1.plus(v2.scale(deltaT / 2))
	v3 := v1.plus(a2.scale(deltaT / 2))
	a3 := p.a(r3, i)

	// Werte am Ende des Zeitschritts, basierend auf soeben berechneten Wertem
	r4 := r1.plus(v3.scale(deltaT))
	v4 := v1.plus(a3.scale(deltaT))
	a4 := p.a(r4, i)

	// Mittelwerte nutzen, um neue Position zu bestimmen
	rEnde := r1.plus(
		v1.plus(v2.scale(2)).
			plus(v3.scale(2)).
			plus(v4).
			scale(deltaT / 6))

	vEnde := v1.plus(
		a1.plus(a2.scale(2)).
			plus(a3.scale(2)).
			plus(a4).
			scale(deltaT / 6))

	p.position[t(i)] = rEnde
	p.velocity[t(i)] = vEnde
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

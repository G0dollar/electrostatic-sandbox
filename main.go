package main

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	K                 = 1000.0
	dt                = 1.0 / 60.0
	physicsStep       = 1
	maxTrail          = 400
	fieldSpacing      = 10
	softening         = 1.0
	doubleClickWindow = 400 * time.Millisecond
	doubleClickRadius = 6.0
	particleHitRadius = 8.0
)

var C30 float64 = math.Cos(math.Pi / 6)
var S30 float64 = math.Sin(math.Pi / 6)
var C_30 float64 = math.Cos(-math.Pi / 6)
var S_30 float64 = math.Sin(-math.Pi / 6)

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Electrostatic Simulation")
	game := Game{
		particles: []Particle{},
		menu: ParticleMenu{
			charge:    1,
			mass:      1,
			radius:    5,
			editIndex: -1,
		},
		isRunning:     false,
		draggingIndex: -1,
	}

	ebiten.RunGame(&game)
}

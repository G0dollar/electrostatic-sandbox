package main

import (
	"math"
)

type Particle struct {
	charge float64
	mass   float64
	radius float64

	position     Vec2
	velocity     Vec2
	acceleration Vec2

	fixed bool

	trail []Vec2
}

type ParticleMenu struct {
	charge float64
	mass   float64
	radius float64

	position Vec2
	velocity Vec2

	fixed bool

	activeField int
	input       string

	editIndex int // -1 = creating a new particle, >=0 = editing g.particles[editIndex]
}

func electricFieldAt(point Vec2, particles []Particle) Vec2 {
	var field Vec2

	for _, particle := range particles {
		direction := subtract(point, particle.position)

		distanceSquared := direction.x*direction.x + direction.y*direction.y

		if distanceSquared < softening {
			continue
		}

		distance := math.Sqrt(distanceSquared)

		unitDirection := multiply(direction, 1/distance)

		magnitude := K * particle.charge / distanceSquared

		contribution := multiply(unitDirection, magnitude)

		field = add(field, contribution)
	}

	return field
}

func coulombForce(p1, p2 Particle) Vec2 {
	r := subtract(p1.position, p2.position)
	distSq := r.x*r.x + r.y*r.y
	if distSq < softening {
		distSq = softening
	}
	distance := math.Sqrt(distSq)
	forceMagnitude := K * p1.charge * p2.charge / distSq
	return multiply(r, forceMagnitude/distance)
}

package main

import "math"

type Vec2 struct {
	x float64
	y float64
}

func add(v1, v2 Vec2) Vec2 {
	return Vec2{v1.x + v2.x, v1.y + v2.y}
}

func subtract(v1, v2 Vec2) Vec2 {
	return Vec2{v1.x - v2.x, v1.y - v2.y}
}

func multiply(v Vec2, m float64) Vec2 {
	return Vec2{m * v.x, m * v.y}
}

func magnitude(v Vec2) float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

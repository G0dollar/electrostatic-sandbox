package main

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) updateMenu() {
	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) &&
		x >= 300 && x <= 520 && y >= 150 && y <= 170 {
		ratio := float64(x-300) / 220.0
		g.menu.charge = -10 + ratio*20
		if g.menu.charge < -10 {
			g.menu.charge = -10
		}
		if g.menu.charge > 10 {
			g.menu.charge = 10
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		matched := false

		for _, f := range menuFields {
			if float32(x) >= f.x1 && float32(x) <= f.x2 &&
				float32(y) >= f.y1 && float32(y) <= f.y2 {
				g.selectField(f.id)
				matched = true
				break
			}
		}

		if !matched {
			switch {
			case x >= 300 && x <= 318 && y >= 448 && y <= 466: // shifted 408→448
				g.menu.fixed = !g.menu.fixed

			case x >= 225 && x <= 325 && y >= 490 && y <= 525: // shifted 450→490
				g.closeMenu()

			case g.menu.editIndex >= 0 && x >= 350 && x <= 450 && y >= 490 && y <= 525:
				g.removeParticleFromMenu()

			case x >= 475 && x <= 575 && y >= 490 && y <= 525:
				g.saveParticleFromMenu()
			}
		}
	}

	if g.menu.activeField != 0 {
		chars := ebiten.InputChars()
		for _, char := range chars {
			g.menu.input += string(char)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(g.menu.input) > 0 {
				g.menu.input = g.menu.input[:len(g.menu.input)-1]
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.applyInput()
		}
	}
}

func (g *Game) applyInput() {
	if g.menu.activeField == 0 || g.menu.input == "" {
		return
	}
	value, err := strconv.ParseFloat(g.menu.input, 64)
	if err != nil {
		return
	}
	switch g.menu.activeField {
	case 1:
		g.menu.mass = value
	case 2:
		g.menu.velocity.x = value
	case 3:
		g.menu.velocity.y = value
	case 4:
		g.menu.position.x = value
	case 5:
		g.menu.position.y = value
	case 6:
		g.menu.radius = value
	}
	g.menu.input = ""
	g.menu.activeField = 0
}

func (g *Game) openEditMenu(index int) {
	p := g.particles[index]
	g.menu = ParticleMenu{
		charge:    p.charge,
		mass:      p.mass,
		radius:    p.radius,
		position:  p.position,
		velocity:  p.velocity,
		fixed:     p.fixed,
		editIndex: index,
	}
	g.menuOpen = true
}

func (g *Game) closeMenu() {
	g.menuOpen = false
	g.menu = ParticleMenu{charge: 1, mass: 1, radius: 5, editIndex: -1}
}

func (g *Game) saveParticleFromMenu() {
	g.applyInput()

	if g.menu.mass <= 0 {
		g.menu.mass = 1
	}
	if g.menu.radius <= 0 {
		g.menu.radius = 5
	}

	p := Particle{
		charge:   g.menu.charge,
		mass:     g.menu.mass,
		radius:   g.menu.radius,
		position: g.menu.position,
		velocity: g.menu.velocity,
		fixed:    g.menu.fixed,
	}

	if g.menu.editIndex >= 0 && g.menu.editIndex < len(g.particles) {
		g.particles[g.menu.editIndex] = p
	} else {
		g.particles = append(g.particles, p)
	}

	g.closeMenu()
}

func (g *Game) removeParticleFromMenu() {
	if g.menu.editIndex >= 0 && g.menu.editIndex < len(g.particles) {
		g.particles = append(g.particles[:g.menu.editIndex], g.particles[g.menu.editIndex+1:]...)
	}
	g.closeMenu()
}

func (g *Game) selectField(field int) {
	g.applyInput()

	g.menu.activeField = field
	g.menu.input = fmt.Sprintf("%g", g.menuFieldRawValue(field))
}

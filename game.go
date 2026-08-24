package main

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	particles []Particle
	field     [80][60]Vec2

	menuOpen bool
	menu     ParticleMenu

	isRunning bool

	lastClickTime time.Time
	lastClickX    int
	lastClickY    int

	draggingIndex int
}

type fieldRect struct {
	id             int
	x1, y1, x2, y2 float32
	label          string
}

var menuFields = []fieldRect{
	{1, 230, 200, 490, 225, "Mass"},
	{6, 230, 240, 490, 265, "Radius"},
	{4, 230, 280, 490, 305, "Position X"},
	{5, 230, 320, 490, 345, "Position Y"},
	{2, 230, 360, 490, 385, "Velocity X"},
	{3, 230, 400, 490, 425, "Velocity Y"},
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.isRunning = !g.isRunning
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.closeMenu()
	}

	if g.menuOpen {
		g.updateMenu()
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) &&
		(ebiten.IsKeyPressed(ebiten.KeyControlLeft) ||
			ebiten.IsKeyPressed(ebiten.KeyControlRight)) {
		g.menuOpen = true
		g.menu = ParticleMenu{charge: 1, mass: 1, radius: 5, editIndex: -1}
	}

	// --- Click handling: one hit-test, then branch into
	// double-click-edit / drag-start / empty-space-tracking ---
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		now := time.Now()

		dx := float64(x - g.lastClickX)
		dy := float64(y - g.lastClickY)
		sameSpot := dx*dx+dy*dy <= doubleClickRadius*doubleClickRadius
		isDoubleClick := sameSpot && now.Sub(g.lastClickTime) <= doubleClickWindow

		if idx, ok := g.particleAt(x, y); ok {
			if isDoubleClick {
				g.openEditMenu(idx)
				g.lastClickTime = time.Time{} // consume it, so a 3rd click doesn't chain
			} else {
				g.draggingIndex = idx
				g.lastClickTime = now
				g.lastClickX = x
				g.lastClickY = y
			}
		} else {
			g.lastClickTime = now
			g.lastClickX = x
			g.lastClickY = y
		}
	}

	if g.draggingIndex >= 0 && g.draggingIndex < len(g.particles) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			p := &g.particles[g.draggingIndex]
			p.position = Vec2{float64(x), float64(y)}
			p.velocity = Vec2{} // prevent a velocity spike on release
			p.trail = nil
		} else {
			g.draggingIndex = -1 // released
		}
	}

	if g.isRunning {
		for range physicsStep {
			g.physicsUpdate()
		}
	}

	g.calculateField()

	if !g.isRunning {
		return nil
	}

	for i := range g.particles {
		particle := &g.particles[i]

		particle.trail = append(particle.trail, particle.position)

		if len(particle.trail) > maxTrail {
			particle.trail = particle.trail[1:]
		}
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func (g *Game) calculateField() {
	var wg sync.WaitGroup
	for x := range g.field {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for y := range g.field[x] {
				point := Vec2{float64(x * fieldSpacing), float64(y * fieldSpacing)}
				g.field[x][y] = electricFieldAt(point, g.particles)
			}
		}(x)
	}
	wg.Wait()
}

func (g *Game) particleAt(x, y int) (int, bool) {
	px, py := float64(x), float64(y)
	for i := len(g.particles) - 1; i >= 0; i-- {
		p := g.particles[i]
		dx := px - p.position.x
		dy := py - p.position.y
		hitR := p.radius
		if hitR < particleHitRadius {
			hitR = particleHitRadius
		}
		if dx*dx+dy*dy <= hitR*hitR {
			return i, true
		}
	}
	return -1, false
}

func (g *Game) physicsUpdate() {
	if !g.isRunning {
		return
	}

	for i := range g.particles {
		particle := &g.particles[i]

		if particle.fixed {
			continue
		}
		particle.acceleration = Vec2{}

		for j := range g.particles {
			if i == j {
				continue
			}

			other := &g.particles[j]

			force := coulombForce(*particle, *other)

			particle.acceleration = add(
				particle.acceleration,
				multiply(force, 1/particle.mass),
			)
		}
	}

	for i := range g.particles {
		particle := &g.particles[i]

		if particle.fixed {
			continue
		}

		particle.velocity = add(particle.velocity, multiply(particle.acceleration, dt))

		particle.position = add(particle.position, multiply(particle.velocity, dt))
	}
}

func (g *Game) menuFieldRawValue(id int) float64 {
	switch id {
	case 1:
		return g.menu.mass
	case 2:
		return g.menu.velocity.x
	case 3:
		return g.menu.velocity.y
	case 4:
		return g.menu.position.x
	case 5:
		return g.menu.position.y
	case 6:
		return g.menu.radius
	}
	return 0
}

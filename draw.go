package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawField(screen)

	for _, particle := range g.particles {
		particleColor := color.RGBA{0, 0, 255, 255}

		if particle.charge > 0 {
			particleColor = color.RGBA{255, 0, 0, 255}
		}

		vector.FillCircle(
			screen,
			float32(particle.position.x),
			float32(particle.position.y),
			float32(particle.radius),
			particleColor,
			false,
		)

		for _, point := range particle.trail {
			vector.FillCircle(
				screen,
				float32(point.x),
				float32(point.y),
				1,
				particleColor,
				false,
			)
		}

	}

	if !g.menuOpen {
		mx, my := ebiten.CursorPosition()
		if idx, ok := g.particleAt(mx, my); ok {
			p := g.particles[idx]
			info := fmt.Sprintf("charge %.2f  mass %.2f  r %.2f", p.charge, p.mass, p.radius)
			ebitenutil.DebugPrintAt(screen, info, mx+10, my-15)
		}
	}

	g.drawMenu(screen)

	status := "RUNNING (P to pause)"
	if !g.isRunning {
		status = "PAUSED (P to resume)"
	}
	ebitenutil.DebugPrintAt(screen, status, 10, 10)
}

func (g *Game) drawField(screen *ebiten.Image) {
	for x := range g.field {
		for y := range g.field[x] {
			field := g.field[x][y]

			magnitude := magnitude(field)

			if magnitude == 0 {
				continue
			}

			if magnitude < 1e-9 {
				continue
			}

			startX := float32(x * fieldSpacing)
			startY := float32(y * fieldSpacing)

			// Normalize the field vector
			direction := multiply(field, 1/magnitude)

			arrowLength := 2.0 + 3.0*math.Log1p(magnitude)
			if arrowLength > 10 {
				arrowLength = 10
			}

			endX := startX + float32(direction.x*arrowLength)
			endY := startY + float32(direction.y*arrowLength)

			vector.StrokeLine(
				screen,
				startX,
				startY,
				endX,
				endY,
				1,
				color.RGBA{100, 100, 100, 255},
				false,
			)

			drawArrowHead(
				screen,
				Vec2{float64(startX), float64(startY)},
				Vec2{float64(endX), float64(endY)},
			)
		}
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	if !g.menuOpen {
		return
	}

	vector.FillRect(screen, 200, 100, 400, 440, color.RGBA{30, 30, 30, 245}, false)
	vector.StrokeRect(screen, 200, 100, 400, 440, 2, color.RGBA{150, 150, 150, 255}, false)

	title := "CREATE PARTICLE"
	if g.menu.editIndex >= 0 {
		title = "EDIT PARTICLE"
	}
	ebitenutil.DebugPrintAt(screen, title, 320, 120)

	vector.FillRect(screen, 225, 490, 100, 35, color.RGBA{120, 50, 50, 255}, false)
	ebitenutil.DebugPrintAt(screen, "CANCEL", 245, 502)

	if g.menu.editIndex >= 0 {
		vector.FillRect(screen, 350, 490, 100, 35, color.RGBA{150, 60, 30, 255}, false)
		ebitenutil.DebugPrintAt(screen, "REMOVE", 368, 502)
	}

	saveLabel := "ADD"
	if g.menu.editIndex >= 0 {
		saveLabel = "SAVE"
	}
	vector.FillRect(screen, 475, 490, 100, 35, color.RGBA{50, 120, 50, 255}, false)
	ebitenutil.DebugPrintAt(screen, saveLabel, 505, 502)

	ebitenutil.DebugPrintAt(screen, "Charge", 230, 155)
	vector.FillRect(screen, 300, 158, 220, 6, color.RGBA{100, 100, 100, 255}, false)

	sliderMin := 300.0
	sliderMax := 520.0
	sliderPosition := sliderMin + ((g.menu.charge+10)/20)*(sliderMax-sliderMin)

	vector.FillCircle(screen, float32(sliderPosition), 161, 7, color.RGBA{255, 255, 255, 255}, false)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.2f", g.menu.charge), 535, 155)

	for _, f := range menuFields {
		value := fmt.Sprintf("%.3f", g.menuFieldRawValue(f.id))
		g.drawInputField(screen, f.label, value, f.x1, f.y1, f.x2, f.y2, f.id)
	}

	ebitenutil.DebugPrintAt(screen, "Fixed", 230, 450)
	checkboxColor := color.RGBA{70, 70, 70, 255}
	if g.menu.fixed {
		checkboxColor = color.RGBA{80, 200, 80, 255}
	}
	vector.FillRect(screen, 300, 448, 18, 18, checkboxColor, false)
	if g.menu.fixed {
		ebitenutil.DebugPrintAt(screen, "X", 304, 448)
	}
}

func (g *Game) drawInputField(
	screen *ebiten.Image,
	label string,
	value string,
	x1, y1, x2, y2 float32,
	field int,
) {
	ebitenutil.DebugPrintAt(
		screen,
		label,
		int(x1),
		int(y1-18),
	)

	fieldColor := color.RGBA{60, 60, 60, 255}

	if g.menu.activeField == field {
		fieldColor = color.RGBA{80, 80, 120, 255}

		// Show what the user is currently typing
		value = g.menu.input
	}

	vector.FillRect(
		screen,
		x1,
		y1,
		x2-x1,
		y2-y1,
		fieldColor,
		false,
	)

	vector.StrokeRect(
		screen,
		x1,
		y1,
		x2-x1,
		y2-y1,
		1,
		color.RGBA{150, 150, 150, 255},
		false,
	)

	ebitenutil.DebugPrintAt(
		screen,
		value,
		int(x1+8),
		int(y1+5),
	)
}

func drawArrowHead(
	screen *ebiten.Image,
	start Vec2,
	end Vec2,
) {
	// Direction from start to end
	direction := subtract(end, start)

	length := magnitude(direction)

	if length == 0 {
		return
	}

	direction = multiply(direction, 1/length)

	arrowSize := 3.0

	// Rotate direction by +30°
	left := Vec2{
		x: direction.x*C30 - direction.y*S30,
		y: direction.x*S30 + direction.y*C30,
	}

	// Rotate direction by -30°
	right := Vec2{
		x: direction.x*C_30 - direction.y*S_30,
		y: direction.x*S_30 + direction.y*C_30,
	}

	left = multiply(left, -arrowSize)
	right = multiply(right, -arrowSize)

	leftPoint := add(end, left)
	rightPoint := add(end, right)

	vector.StrokeLine(
		screen,
		float32(end.x),
		float32(end.y),
		float32(leftPoint.x),
		float32(leftPoint.y),
		1,
		color.RGBA{100, 100, 100, 255},
		false,
	)

	vector.StrokeLine(
		screen,
		float32(end.x),
		float32(end.y),
		float32(rightPoint.x),
		float32(rightPoint.y),
		1,
		color.RGBA{100, 100, 100, 255},
		false,
	)
}

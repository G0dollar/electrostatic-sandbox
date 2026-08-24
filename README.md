# Electrostats

An interactive electrostatic particle simulation and game written in Go. **Electrostats** models particle physics, charge interactions, and real-time vector math through a graphical interface.

## Project Architecture

The codebase is organized into modular Go source components:

* **`main.go`**: Application entry point responsible for initialization and startup.

* **`game.go`**: Core game loop managing state transitions, updates, and frame logic.

* **`menu.go`**: Navigation, menu UI, and user selection handling.

* **`particle.go`**: Particle object definitions, physical properties (e.g., charge, mass), and force calculations.

* **`vector.go`**: 2D vector mathematics and spatial mechanics (velocity, acceleration, displacement).

* **`draw.go`**: Visual rendering system for graphics, menus, and particles.

## Prerequisites

* **Go 1.20+** installed on your system.

## Getting Started

1; **Clone the repository**

```bash
git clone https://github.com/<your-username>/Electrostats.git
cd Electrostats

```

2; **Download dependencies**

```bash
go mod download

```

3; **Run the application**

```bash
go run .

```

## Key Features

* Real-time particle physics and electrostatic interactions.

* 2D vector calculation engine.

* Interactive GUI menu system.

## Electrostatic Simulation

A 2D electrostatic particle simulator built in Go with [Ebiten](https://ebitengine.org/). Place charged particles on a canvas, watch them attract and repel each other under Coulomb's law, and visualize the electric field they generate in real time.

## Features

* **Coulomb-force physics** — every particle exerts an inverse-square force on every other particle, with softening to avoid singularities at close range

* **Live vector field visualization** — a grid of arrows shows field direction and (log-scaled) magnitude across the whole canvas, recalculated every frame using goroutines (one per column)

* **Particle creation & editing menu** — add particles with a specific charge, mass, radius, position, and velocity; double-click an existing particle to edit or remove it

* **Drag-and-drop** — click and drag any particle to reposition it (this resets its velocity and trail)

* **Motion trails** — moving particles leave a fading trail showing their recent path

* **Fixed particles** — pin a particle in place so it exerts force but never moves (useful for building static field configurations)

* **Play/pause** — freeze physics at any time to inspect the current state

## Controls

|         Input           |                            Action                                |
|------------------------ |------------------------------------------------------------------|
|        `Ctrl+C`         |                Open the "create particle" menu                   |
| Double-click a particle |              Open the "edit particle" menu for it                |
| Click + drag a particle | Move it (only while paused or running — velocity resets on drop) |
|          `P`            |           Toggle play/pause of the physics simulation            |
|         `Esc`           |                     Close the open menu                          |

Inside the particle menu:

* Drag the **Charge** slider (range −10 to +10)

* Click **Mass**, **Radius**, **Position X/Y**, or **Velocity X/Y** to type a value, then press `Enter` to apply it

* Toggle **Fixed** to pin the particle in place

* **ADD**/**SAVE** commits the particle; **REMOVE** deletes it (edit mode only); **CANCEL** discards changes

Positive charges are drawn in red, negative charges in blue.

Getting Started

Prerequisites

* [Go](https://go.dev/dl/) 1.26 or later

* A platform Ebiten supports (Windows, macOS, Linux, or WebAssembly) with the usual OpenGL/graphics dependencies for your OS — see the [Ebiten install guide](https://ebitengine.org/en/documents/install.html) if you hit build errors on Linux (you may need `libgl1-mesa-dev`, `libxrandr-dev`, `libxcursor-dev`, `libxinerama-dev``libxi-dev`, and `libasound2-dev` or similar).

### Run it

```bash
go mod tidy
go run .
```

This opens an 800×600 window titled "Electrostatic Simulation."

### Build a binary

```bash
go build -o electrostatic-sim .
```

## Project Structure

```text
.
├── main.go      # Entry point: window setup and initial game state
├── game.go      # Core game loop — input handling, physics stepping, field calculation
├── particle.go  # Particle struct, Coulomb force, and electric field math
├── vector.go    # Minimal 2D vector type and operations
├── menu.go      # Particle creation/edit menu logic (input handling, save/remove)
├── draw.go      # All rendering: field arrows, particles, trails, menu UI
├── go.mod
└── go.sum
```

## How It Works

Each frame, when running:

1. **Force accumulation** — for every particle, the Coulomb force from every other particle is summed to get acceleration (`F = K·q₁·q₂/r²`, softened at short range to avoid blow-ups)
2. **Integration** — velocity and position are updated via simple explicit Euler integration, run `physicsStep` (50) times per frame for stability
3. **Field sampling** — the electric field is sampled on an 80×60 grid across the canvas, computed in parallel (one goroutine per column) and rendered as directional arrows scaled by log-magnitude
4. **Trails** — each particle's last 400 positions are recorded and drawn as a fading path

Key constants (in `main.go`) you can tune:

|    Constant    |                           Meaning                          |
|----------------|------------------------------------------------------------|
|      `K`       | Coulomb's constant (scaled for visual scale, not SI units) |
|      `dt`      |                      Physics timestep                      |
|  `physicsStep` |            Physics substeps per rendered frame             |
|   `softening`  |      Minimum squared distance to prevent singular forces   |
| `fieldSpacing` |          Pixel spacing between field-arrow samples         |
|   `maxTrail`   |               Max trail length per particle                |

## Known Limitations / Ideas for Later

* Field grid resolution and canvas size are hardcoded to 800×600

* No save/load of particle configurations

* No collision handling between particles

* Explicit Euler integration is simple but will drift over long/close-orbit runs — a symplectic or RK4 integrator would be more accurate

## License

No license specified yet — add one (MIT is a common default) if you plan to share this publicly.

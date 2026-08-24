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

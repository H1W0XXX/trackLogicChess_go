# Track Logic Chess

**4x4 Rotating Chess (GUI Supported)**

This project implements a simplified version of Track Logic Chess on a 4×4 board with rotating rings, and supports enabling a graphical interface and AI opponent via startup parameters.

---

## Game Rules Overview

* **Board**: 4×4 grid (16 cells), divided into Outer Ring (12 cells) and Inner Ring (4 cells)
* **Win Condition**: Four in a row (horizontal, vertical, or diagonal)
* **Turn Order**: Black (●) moves first, followed by White (○)
* **Rotation Mechanics**:

  * **Outer Ring**: After each move, the outer ring shifts one cell in a fixed direction (clockwise or counterclockwise)
  * **Inner Ring**: Similarly rotates one step per move
  * Both ring directions are specified at startup and cannot be changed during the game

---

## Command-Line Options

Control rotation directions, AI usage, and UI mode through command-line flags:

| Flag     | Type   | Default      | Description                                                   |
| -------- | ------ | ------------ | ------------------------------------------------------------- |
| `-outer` | int    | `0`          | Outer ring direction: `0` = clockwise, `1` = counterclockwise |
| `-inner` | int    | `0`          | Inner ring direction: `0` = clockwise, `1` = counterclockwise |
| `-ai`    | bool   | `true`       | Enable AI for the White player                                |
| `-ui`    | string | `"terminal"` | UI mode: `"terminal"` or `"gui"`                              |

---

## Example Command (GUI Mode)

```bash
./tracklogicchess -outer 0 -inner 1 -ai=true -ui=gui
```

* Outer ring rotates clockwise (`-outer 0`) and inner ring counterclockwise (`-inner 1`)
* White (second player) is controlled by AI
* Launches the game with the graphical interface

---

## GUI Notes

* Built with [Ebiten](https://ebiten.org) for basic graphics and input handling
* AI is enabled by default, with the human player controlling Black (first player)
* After each move, the board automatically performs the preset ring rotations
* Initial version focuses on core functionality; future updates may add animations, interactive buttons, and enhanced UI

---

## Development & Build Requirements

* **Go** 1.18 or higher
* Supported platforms: Windows, macOS, Linux
* Recommended: Use Go Modules

```bash
go build -o tracklogicchess ./cmd/tracklogicchess
```

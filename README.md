# Go Arcade (Terminal Game Suite)

A production-quality, offline, terminal-native arcade game suite written in Go. Features classic **Snake (Nokia-style)**, **Pacman (Arcade Maze Chase)**, and **Ball & Plate (Breakout)** running inside your terminal with pure cell buffers, box-drawing chrome, multiple visual themes, concrete mechanical difficulty levels, and atomic local persistence.

---

## Features

- **Shared Game Engine**: Unified rendering, non-blocking tick-based event loop, pause management, and input mapping via direct `github.com/gdamore/tcell/v2`.
- **Three Complete Arcade Titles**:
  - **Snake**: Nokia-style grid action with wrap-around boundaries on Easy, deadly arena walls on Medium, and deadly obstacle pellets on Hard.
  - **Pacman**: Maze chase featuring dots, power pellets, vulnerable frightened states, and AI archetypes (Random Wander, Manhattan Greedy Chase, and Directional Ambush).
  - **Ball & Plate**: Breakout-style brick buster with sub-cell floating-point physics, paddle impact angle deflection, tough 2-hit bricks, and life tracking.
- **Three Visual Themes**:
  - **Retro Green**: Nokia LCD / Game Boy green-on-black aesthetic.
  - **Neon**: Cyberpunk arcade cabinet palette (cyan, yellow, magenta).
  - **Monochrome**: Pure high-contrast ASCII (`#`, `o`, `*`) with zero color reliance.
- **Three Concrete Difficulty Modes**: Each mode introduces genuine mechanical differences (arena size, speeds, lives, AI, obstacles), not cosmetic labels.
- **Persistent Play History & Statistics**: Logs every session to disk atomically. Includes a built-in History & Stats browser with personal bests, total playtime, and per-game breakdowns.
- **Robust Terminal Safety**: Clean shutdown and guaranteed raw mode teardown (`Fini()`) on normal quit, menu exit, or unexpected panic recovery.
- **100% Offline & Headless Testable**: Zero external network dependencies, telemetry, or remote calls.

---

## Installation & Running

### Requirements
- Go 1.22 or higher
- A terminal emulator supporting ANSI or Unicode escape sequences (standard Terminal, iTerm2, Alacritty, Kitty, WezTerm, Windows Terminal, etc.)
- Minimum terminal size: 80 columns by 24 rows

### Quick Start

```bash
# Clone the repository
git clone https://github.com/srinjoypramanick/Golang_games.git
cd Golang_games

# Build the suite
go build -o arcade cmd/arcade/main.go

# Launch the arcade suite
./arcade
```

Or run directly without manual binary compilation:

```bash
go run cmd/arcade/main.go
```

### Running Tests

All game simulation logic and storage aggregations are decoupled from terminal display hardware and can be tested headlessly:

```bash
go test -v ./...
go vet ./...
```

---

## Controls & Keybinding Reference

The suite supports intuitive dual-control schemes (Arrow keys and WASD) everywhere:

### Main Menu & Configuration Picker

| Key | Action |
| --- | --- |
| `↑` / `w` / `W` | Move selection up |
| `↓` / `s` / `S` | Move selection down |
| `Enter` / `Space` | Confirm selection / Launch game |
| `1` – `5` | Instant jump to menu option (in Main Menu) |
| `Esc` / `q` / `Q` | Exit back / Return to previous screen / Quit suite |

### In-Game (All Games)

| Key | Snake | Pacman | Ball & Plate |
| --- | --- | --- | --- |
| `↑` / `w` / `W` | Move Up | Move Up | Launch Ball |
| `↓` / `s` / `S` | Move Down | Move Down | - |
| `←` / `a` / `A` | Move Left | Move Left | Move Paddle Left |
| `→` / `d` / `D` | Move Right | Move Right | Move Paddle Right |
| `Space` / `Enter` | Start game / Advance screen | Start game / Advance screen | Launch Ball / Start game |
| `p` / `P` | Pause / Resume | Pause / Resume | Pause / Resume |
| `q` / `Q` / `Esc` | Quit to Main Menu | Quit to Main Menu | Quit to Main Menu |

### History & Statistics Screen

| Key | Action |
| --- | --- |
| `←` / `→`, `Tab`, `1`–`4` | Switch game filter tab (`All`, `Snake`, `Pacman`, `Ball & Plate`) |
| `↑` / `↓`, `w` / `s` | Scroll through past session log rows |
| `Esc` / `q` / `Q` / `Enter` | Return to Main Menu |

---

## Difficulty Parameters

Each game has three distinct difficulty tiers with concrete mechanical parameter adjustments:

### Snake

| Parameter | Easy | Medium | Hard |
| --- | --- | --- | --- |
| **Tick Interval** | 180 ms (relaxed) | 130 ms (standard) | 90 ms (fast) |
| **Arena Size** | 40 × 20 cells | 34 × 18 cells | 28 × 16 cells |
| **Wall Collisions** | **Wraps around edges** (safe) | **Deadly** (game over on impact) | **Deadly** (game over on impact) |
| **Deadly Obstacles** | None | None | **Spawns lethal obstacle pellets** periodically |

### Pacman

| Parameter | Easy | Medium | Hard |
| --- | --- | --- | --- |
| **Active Ghosts** | 2 ghosts (Blinky, Pinky) | 3 ghosts (Blinky, Pinky, Inky) | 4 ghosts (Blinky, Pinky, Inky, Clyde) |
| **Ghost AI** | Random wander (casual) | Greedy Manhattan chase | Greedy chase + **Directional Ambush AI** (targets 4 tiles ahead) |
| **Power Pellet Duration** | 50 ticks (~7.0 seconds) | 35 ticks (~4.4 seconds) | 20 ticks (~2.2 seconds) |
| **Player Lives** | 3 lives | 3 lives | 2 lives |
| **Tick Interval** | 140 ms | 125 ms | 110 ms |

### Ball & Plate (Breakout)

| Parameter | Easy | Medium | Hard |
| --- | --- | --- | --- |
| **Paddle Width** | 7 cells (wide) | 5 cells (standard) | 3 cells (narrow) |
| **Base Ball Speed** | 0.32 cells/tick | 0.40 cells/tick | 0.48 cells/tick |
| **Speed Scaling** | None (constant speed) | +0.002 cells/tick per brick cleared | +0.004 cells/tick per brick cleared |
| **Brick Grid Layout** | Gapped checkerboard (fewer bricks) | Full solid grid | Full solid grid + **Tough 2-hit brick row** |
| **Player Lives** | 4 lives | 3 lives | 2 lives |
| **Tick Interval** | 40 ms | 38 ms | 35 ms |

---

## Visual Themes

All games support three distinct themes switchable at launch:

| Element | Retro Green | Neon | Monochrome |
| --- | --- | --- | --- |
| **Aesthetic** | Nokia LCD Green-on-Black | Cyberpunk Arcade Cabinet | High-Contrast Pure ASCII |
| **Background** | Black | Black | Black |
| **Wall Style** | `█` (Dark Green) | `║` (Blue Violet) | `#` (White) |
| **Snake Head/Body** | `█` (Lime / Green) | `█` (Aqua head, Fuchsia body) | `o` head, `#` body |
| **Snake Food** | `◆` (Lime) | `●` (Yellow) | `*` (White) |
| **Pacman Player** | `█` (Green) | `>` `<` `^` `v` (Yellow) | `C` (White) |
| **Pacman Ghosts** | `G` (Forest Green) | `M` (Red, Pink, Cyan, Orange) | `M` (White) |
| **Frightened Ghost** | `w` (Olive Green) | `w` (Deep Sky Blue / White flash) | `w` (Gray) |
| **Breakout Plate** | `▬` (Green) | `━` (Yellow) | `-` (White) |
| **Breakout Ball** | `●` (Lime) | `●` (Aqua) | `o` (White) |
| **Breakout Bricks** | `▀` (Green gradient) | `█` (Rainbow row gradient) | `=` (White, tough `#`) |

---

## Persistence & Storage

### Storage Location
Session records are saved to an OS-standard local configuration path:
- **macOS**: `~/Library/Application Support/goarcade/history.json`
- **Linux**: `~/.config/goarcade/history.json`
- **Windows**: `%AppData%\goarcade\history.json`

### Record Schema
```json
{
  "game": "snake",
  "difficulty": "medium",
  "theme": "Neon",
  "score": 140,
  "outcome": "wall",
  "played_at": "2026-09-21T22:30:00Z",
  "duration": 48000000000
}
```

### Atomic Writes & Corruption Recovery
- **Atomic Save**: When a session finishes, history is serialized and written to `history.json.tmp`. Once safely flushed to disk, it is atomically renamed to `history.json` via `os.Rename`. A power failure or crash mid-write can never corrupt existing history.
- **Self-Healing Corruption Handling**: If `history.json` is modified externally with invalid JSON or damaged, the application does not crash. Instead, it renames the corrupted file to `history.json.corrupt-<timestamp>` for inspection and initializes a clean history log.
- **How to Reset History**:
  Simply delete `history.json` from the directory path listed above, or run:
  ```bash
  rm -f "$(go run -e 'p, _ := os.UserConfigDir(); println(p)' 2>/dev/null)/goarcade/history.json"
  ```

---

## Architecture & Design Notes

```
cmd/arcade/main.go          Application entry point, menu dispatcher, panic recovery
internal/
├── engine/                 Game engine core (zero game-specific rules)
│   ├── difficulty.go       Difficulty enum (Easy/Medium/Hard) and profiles
│   ├── entity.go           Position, Sprite, Rect (AABB collision)
│   ├── game.go             Universal Game interface contract
│   ├── input.go            Hardware-to-Action translation (Arrows + WASD)
│   ├── loop.go             Fixed-interval tick loop, non-blocking channel drain
│   ├── screen.go           tcell Screen wrapper with box-drawing & text helpers
│   ├── theme.go            Theme struct & ThemeRegistry
│   └── engine_test.go      Headless unit tests
├── history/                Decoupled persistence layer (no tcell imports)
│   ├── record.go           Record struct & aggregation queries
│   ├── store.go            Atomic JSON save/load & corruption recovery
│   └── history_test.go     Store round-trip, atomic rename, query tests
├── menu/                   Menu and statistics screens (built on engine.Screen)
│   ├── mainmenu.go         Title banner & game selector
│   ├── picker.go           Interactive Difficulty & Theme picker
│   └── history.go          Browsable/scrollable history table & summary panel
└── games/
    ├── snake/              Nokia-style Snake implementation & tests
    ├── pacman/             Maze parser, ghost AI archetypes, game logic & tests
    └── ballplate/          Breakout with sub-cell physics, deflection math & tests
```

### Dependency Rules Enforced
1. `internal/engine` and `internal/history` depend on no other internal packages.
2. `internal/games/*` depend only on `internal/engine` and `internal/history`.
3. `internal/menu` depends only on `internal/engine` and `internal/history`.
4. No game package imports another game package or `internal/menu`.
5. `internal/history` does not import `tcell` or anything rendering-related, ensuring 100% headless testing.

### Sub-Cell Motion & Paddle Deflection Physics
In `internal/games/ballplate`:
- Ball coordinates and velocities are tracked as double-precision floating-point numbers (`float64` X, Y, Vx, Vy).
- Terminal grid coordinates are derived by rounding (`math.Round`) at collision and render time, creating smooth motion without grid locking.
- When the ball strikes the paddle, the horizontal deflection angle is determined by the relative impact point from the paddle's center:
  $$\text{offset} = \frac{X_{\text{hit}} - X_{\text{center}}}{W_{\text{paddle}} / 2}$$
  $$V_x = \text{speed} \cdot \sin(\text{offset} \cdot 60^\circ)$$
  $$V_y = -\left|\text{speed} \cdot \cos(\text{offset} \cdot 60^\circ)\right|$$
  This allows skilled players to aim ball bounces with precision.

# Go Arcade — Terminal Game Suite (Snake, Pacman, Ball & Plate)

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Offline](https://img.shields.io/badge/network-100%25%20offline-success)](README.md)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](README.md)

A production-quality, fully offline, terminal-native arcade game suite written in modern Go. Launches into a game-selection menu and runs classic arcade games — **Snake (Nokia-style)**, **Pacman (Arcade Maze Chase)**, and **Ball & Plate (Breakout)** — entirely inside your terminal using direct raw-mode screen buffers, box-drawing chrome, multiple visual themes, concrete mechanical difficulty tiers, and atomic local persistence with deep statistics tracking.

---

## Table of Contents

- [Features & Highlights](#features--highlights)
- [Installation & Distribution](#installation--distribution)
  - [1. Install via `go install`](#1-install-via-go-install-all-platforms)
  - [2. Build from Source](#2-build-from-source-all-platforms)
  - [3. Pre-Built Release Binaries & Cross-Compilation](#3-pre-built-release-binaries--cross-compilation)
  - [4. Uninstallation & Data Reset](#4-uninstallation--data-reset)
- [Keybinding Reference](#keybinding-reference)
- [Difficulty Parameters](#difficulty-parameters)
- [Visual Themes](#visual-themes)
- [Deep Stats, Streaks & Achievement Badges](#deep-stats-streaks--achievement-badges)
- [Export & Reset CLI Commands](#export--reset-cli-commands)
- [Architecture & Design Notes](#architecture--design-notes)
- [Running Headless Tests](#running-headless-tests)

---

## Features & Highlights

1. **Shared Game Engine Core**: One unified rendering buffer, fixed-ticker loop, non-blocking input queue, and pause management system built directly on `github.com/gdamore/tcell/v2`.
2. **Three Fully-Featured Arcade Titles**:
   - **Snake (Nokia-style)**: Grid movement on bordered arenas with edge-wrapping on Easy, deadly boundary walls on Medium, and deadly spawning obstacle pellets on Hard.
   - **Pacman**: Authentic maze chase featuring dots, power pellets with frightened ghost states, and distinct ghost AI personalities (Random Wander, Manhattan Greedy Chase, and Directional Ambush).
   - **Ball & Plate (Breakout)**: Sub-cell floating-point ball physics, continuous angle-deflection paddle dynamics, tough 2-hit bricks, and life tracking.
3. **Three Mechanical Difficulty Modes (Easy, Medium, Hard)**: Concretely distinct gameplay tuning parameters (arena dimensions, tick intervals, ghost count and AI, paddle size, speeds, obstacles), not superficial cosmetic labels.
4. **Three Built-in Visual Themes**:
   - **Retro Green**: Nokia LCD / classic Game Boy green-on-black aesthetic.
   - **Neon**: Vibrant cyberpunk arcade cabinet colors (cyan, magenta, yellow, orange).
   - **Monochrome / High-Contrast**: Pure ASCII glyphs (`#`, `o`, `*`, `=`) designed for guaranteed compatibility on legacy consoles and accessibility without color reliance.
5. **Persistent History & Deep Stats Tracking**:
   - Every completed or quit session is atomically saved to local JSON (`os.UserConfigDir()/goarcade/history.json`).
   - Per-game metrics: near-misses, food eaten, ghosts eaten, bricks broken, longest rally.
   - Consecutive calendar-day play streaks (`current_streak`, `longest_streak`) displayed on the menu and history screens.
   - 10+ badge achievement system stored in `achievements.json` with dedicated in-game achievement browser and unlock hints.
   - Lightweight 30-day ASCII play activity heatmap (`·`, `▪`, `█`).
   - Local JSON and CSV export commands.
6. **100% Offline & Private**: Zero network dependencies, telemetry, update checkers, or remote calls.
7. **Terminal Safety**: Guaranteed raw-mode restoration (`Fini()`) on clean exit, menu quit, Ctrl-C, or panic recovery.

---

## Installation & Distribution

### 1. Install via `go install` (All Platforms)

Requires Go 1.22 or higher installed on your system.

```bash
go install github.com/srinjoypramanick/Golang_games/cmd/arcade@latest
```

This places the `arcade` binary in `$(go env GOPATH)/bin` (or `%USERPROFILE%\go\bin` on Windows).

#### PATH Verification & Fix:

- **macOS / Linux** (`bash` / `zsh`):
  ```bash
  # Check if GOPATH/bin is on PATH:
  echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "GOPATH/bin is in PATH" || echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc
  ```
- **Windows** (PowerShell):
  ```powershell
  # Check and add Go bin to User PATH:
  $goBin = "$(go env GOPATH)\bin"
  if ($env:Path -notlike "*$goBin*") { [Environment]::SetEnvironmentVariable("Path", $env:Path + ";$goBin", "User") }
  ```

---

### 2. Build from Source (All Platforms)

Clone the repository and build the binary:

#### Linux & macOS:
```bash
git clone https://github.com/srinjoypramanick/Golang_games.git
cd Golang_games
go build -buildvcs=false -o arcade ./cmd/arcade
./arcade
```

#### Windows (Command Prompt / PowerShell):
```cmd
git clone https://github.com/srinjoypramanick/Golang_games.git
cd Golang_games
go build -buildvcs=false -o arcade.exe .\cmd\arcade
arcade.exe
```

---

### 3. Pre-Built Release Binaries & Cross-Compilation

To generate all cross-compiled release binaries from a single machine into `dist/`, run the cross-compilation script:

#### Linux / macOS:
```bash
bash scripts/build_all.sh
```

#### Manual Cross-Compilation Commands:
```bash
# Linux (64-bit AMD & ARM)
GOOS=linux   GOARCH=amd64 go build -buildvcs=false -o dist/arcade_linux_amd64       ./cmd/arcade
GOOS=linux   GOARCH=arm64 go build -buildvcs=false -o dist/arcade_linux_arm64       ./cmd/arcade

# macOS (Intel & Apple Silicon)
GOOS=darwin  GOARCH=amd64 go build -buildvcs=false -o dist/arcade_darwin_amd64      ./cmd/arcade
GOOS=darwin  GOARCH=arm64 go build -buildvcs=false -o dist/arcade_darwin_arm64      ./cmd/arcade

# Windows (64-bit AMD)
GOOS=windows GOARCH=amd64 go build -buildvcs=false -o dist/arcade_windows_amd64.exe  ./cmd/arcade
```

#### Running Downloaded Pre-Built Binaries:

- **Linux**:
  ```bash
  chmod +x arcade_linux_amd64
  ./arcade_linux_amd64
  ```
- **macOS (Gatekeeper Workaround)**:
  Unsigned macOS binaries downloaded from browsers or GitHub releases trigger an "unidentified developer" warning. Remove the quarantine flag or right-click Open:
  ```bash
  chmod +x arcade_darwin_arm64
  xattr -d com.apple.quarantine arcade_darwin_arm64
  ./arcade_darwin_arm64
  ```
  *(Alternatively: Right-click the file in Finder → Open → Click "Open" on the prompt).*
- **Windows**:
  Run `arcade_windows_amd64.exe` inside Windows Terminal, PowerShell, or Command Prompt.

---

### 4. Uninstallation & Data Reset

#### Removing the Binary:
- **macOS / Linux**:
  ```bash
  rm -f "$(go env GOPATH)/bin/arcade" ./arcade
  ```
- **Windows**:
  ```powershell
  del "$env:USERPROFILE\go\bin\arcade.exe"
  del .\arcade.exe
  ```

#### Wiping Saved History & Achievements:
History and badges are stored under `os.UserConfigDir()/goarcade/`:
- **macOS**: `~/Library/Application Support/goarcade/`
- **Linux**: `~/.config/goarcade/`
- **Windows**: `%APPDATA%\goarcade\`

To completely wipe all local data:
- **macOS**:
  ```bash
  rm -rf "$HOME/Library/Application Support/goarcade"
  ```
- **Linux**:
  ```bash
  rm -rf "$HOME/.config/goarcade"
  ```
- **Windows**:
  ```powershell
  Remove-Item -Recurse -Force "$env:APPDATA\goarcade"
  ```
- **In-App / CLI**:
  ```bash
  arcade --reset-data
  ```

---

## Keybinding Reference

Dual keyboard mappings (Arrow keys and WASD) are fully supported across all screens:

### Main Menu & Configuration Pickers

| Key | Action |
| --- | --- |
| `↑` / `w` / `W` | Move selection up |
| `↓` / `s` / `S` | Move selection down |
| `Enter` / `Space` | Select / Confirm choice / Launch game |
| `1` – `7` | Numeric quick selection (in Main Menu) |
| `q` / `Q` / `Esc` | Return to previous menu / Exit arcade suite |

### In-Game Controls (All Games)

| Key | Snake | Pacman | Ball & Plate |
| --- | --- | --- | --- |
| `↑` / `w` / `W` | Steer Up | Steer Up | Launch Ball |
| `↓` / `s` / `S` | Steer Down | Steer Down | - |
| `←` / `a` / `A` | Steer Left | Steer Left | Move Paddle Left |
| `→` / `d` / `D` | Steer Right | Steer Right | Move Paddle Right |
| `Space` / `Enter` | Start game / Advance | Start game / Advance | Launch Ball / Start |
| `p` / `P` | Pause / Resume | Pause / Resume | Pause / Resume |
| `q` / `Q` / `Esc` | Quit to Main Menu | Quit to Main Menu | Quit to Main Menu |

### History, Badges & Stats Browser

| Key | Action |
| --- | --- |
| `←` / `→`, `a` / `d`, `Tab` | Cycle filter tab (`All`, `Snake`, `Pacman`, `Ball & Plate`, `Achievements`) |
| `1` – `5` | Jump directly to tab index |
| `↑` / `↓`, `w` / `s` | Scroll past session rows or badges |
| `q` / `Q` / `Esc` / `Enter` | Return to Main Menu |

---

## Difficulty Parameters

Each game features three concrete, mechanically distinct difficulty modes:

### 1. Snake (Nokia-style)

| Parameter | Easy | Medium | Hard |
| --- | --- | --- | --- |
| **Tick Interval** | 180 ms | 130 ms | 90 ms |
| **Arena Dimensions** | 40 × 20 cells | 34 × 18 cells | 28 × 16 cells |
| **Wall Collisions** | **Wraps around edges** (safe) | **Lethal** (game over) | **Lethal** (game over) |
| **Deadly Obstacles** | None | None | **Spawns lethal obstacle pellets** periodically |

### 2. Pacman (Arcade Maze Chase)

| Parameter | Easy | Medium | Hard |
| --- | --- | --- | --- |
| **Ghost Count** | 2 ghosts (Blinky, Pinky) | 3 ghosts (Blinky, Pinky, Inky) | 4 ghosts (Blinky, Pinky, Inky, Clyde) |
| **Ghost AI** | Random wander (casual) | Greedy Manhattan chase | Greedy chase + **Directional Ambush AI** (targets 4 tiles ahead) |
| **Pellet Duration** | 50 ticks (~7.0 s) | 35 ticks (~4.4 s) | 20 ticks (~2.2 s) |
| **Player Lives** | 3 lives | 3 lives | 2 lives |
| **Tick Interval** | 140 ms | 125 ms | 110 ms |

### 3. Ball & Plate (Breakout)

| Parameter | Easy | Medium | Hard |
| --- | --- | --- | --- |
| **Paddle Width** | 7 cells (wide) | 5 cells (standard) | 3 cells (narrow) |
| **Base Ball Speed** | 0.32 cells/tick | 0.40 cells/tick | 0.48 cells/tick |
| **Speed Scaling** | Constant speed | +0.002 cells/tick per brick cleared | +0.004 cells/tick per brick cleared |
| **Brick Grid** | Checkerboard gaps (fewer bricks) | Solid full grid | Solid grid + **Tough 2-hit brick row** |
| **Player Lives** | 4 lives | 3 lives | 2 lives |
| **Tick Interval** | 40 ms | 38 ms | 35 ms |

---

## Visual Themes

All three titles support three visual themes, selectable prior to game launch:

| Element | Retro Green | Neon | Monochrome / High-Contrast |
| --- | --- | --- | --- |
| **Aesthetic** | Nokia LCD Green-on-Black | Cyberpunk Arcade Cabinet | High-Contrast Pure ASCII Baseline |
| **Background** | Black | Black | Black |
| **Wall** | `█` (Dark Green) | `║` (Blue Violet) | `#` (White) |
| **Snake Body / Head** | `█` (Green / Lime) | `█` (Aqua head, Fuchsia body) | `o` head, `#` body (White) |
| **Snake Food** | `◆` (Lime) | `●` (Yellow) | `*` (White) |
| **Pacman** | `█` (Green) | `>` `<` `^` `v` (Yellow) | `C` (White) |
| **Ghosts** | `G` (Forest Green) | `M` (Red, Pink, Cyan, Orange) | `M` (White) |
| **Frightened Ghost** | `w` (Olive Green) | `w` (Deep Sky Blue / White flash) | `w` (Gray) |
| **Breakout Plate** | `▬` (Green) | `━` (Yellow) | `-` (White) |
| **Breakout Ball** | `●` (Lime) | `●` (Aqua) | `o` (White) |
| **Breakout Bricks** | `▀` (Green gradient) | `█` (Rainbow row gradient) | `=` (White, tough `#`) |

> [!NOTE]
> **Monochrome / High-Contrast Theme**: Uses exclusively plain 7-bit ASCII characters (`#`, `o`, `*`, `=`, `-`). It is guaranteed to display cleanly on legacy terminals (e.g., Windows `cmd.exe` or minimal serial consoles) without requiring Unicode box-drawing support.

---

## Deep Stats, Streaks & Achievement Badges

Go Arcade records rich telemetry for every game played:

### 1. Per-Game Detailed Metrics
- **Snake**: `food_eaten`, `max_length_reached`, `ticks_survived`, `near_misses` (counts turns made 1 cell before fatal collisions).
- **Pacman**: `dots_eaten`, `power_pellets_used`, `ghosts_eaten`, `levels_cleared`, `lives_lost`.
- **Ball & Plate**: `bricks_broken`, `max_ball_speed_reached`, `plate_hits`, `lives_lost`, `longest_rally` (consecutive plate bounces without losing a ball).

### 2. Calendar-Day Play Streaks
Tracks consecutive calendar days with at least one game played. Evaluated using local timezone day boundaries.
- Displayed prominently in the Main Menu banner: `🔥 3-Day Streak! (Personal Best: 7 days)`.
- Fallback under Monochrome: `[Streak: 3 days | Best: 7 days]`.

### 3. Achievement Badges System
A set of 10 badges evaluated automatically upon game completion:

| Badge | Title | Requirement | Hint |
| --- | --- | --- | --- |
| `first_blood` | **First Blood** | Complete your first session of any game | Play and finish or exit a session in any game. |
| `century_club` | **Century Club** | Score 100+ in any single Snake session | Eat at least 10 food pellets in one Snake game. |
| `ghost_hunter` | **Ghost Hunter** | Eat 20 ghosts total across all Pacman sessions | Consume frightened ghosts during power pellet periods. |
| `brick_breaker` | **Brick Breaker** | Clear a full Ball & Plate level on Hard | Destroy every brick including tough bricks on Hard. |
| `marathon` | **Marathon** | Single session lasting 5+ minutes | Keep any session actively playing for at least 300s. |
| `dedicated` | **Dedicated** | Achieve a 7-day play streak | Play at least one round each day for 7 consecutive days. |
| `completionist`| **Completionist** | Play all three games at least once | Play at least one round of Snake, Pacman, and Ball & Plate. |
| `perfectionist`| **Perfectionist** | Clear a Pacman level without losing a life | Eat all dots and win a Pacman game with 0 deaths. |
| `speed_demon`  | **Speed Demon** | Survive 60+ seconds on Snake Hard | Dodge walls and obstacles at top speed for over 60s. |
| `night_owl`    | **Night Owl** | Play a session between midnight and 4:00 AM | Launch and play a round late at night (00:00 - 04:00). |

Badges are browsable in the **Achievements** tab (`[5]`) on the History screen with unlocked dates and unlock hints for locked badges.

### 4. 30-Day Activity Heatmap
Rendered on the History screen, visually displaying session density over the last 30 calendar days:
```
── 30-DAY HEATMAP ──
· · ▪ · · █ · · ▪ ▪ · · · · · ▪ █ █ ·
· none  ▪ 1-2  █ 3+
```

---

## Export & Reset CLI Commands

### Export History & Achievements
Export your offline data to a human-readable file:

```bash
# Export all sessions and badges to JSON:
./arcade --export --format=json --output=my_arcade_backup.json

# Export sessions to CSV:
./arcade --export --format=csv --output=my_arcade_sessions.csv
```

You can also trigger an export directly inside the game from the Main Menu (**Option 5: Export Data**).

### Reset Data
Wipe all history and achievement records with interactive confirmation:

```bash
./arcade --reset-data
# Prompts: Warning: This will permanently delete all session history and achievements.
# Are you sure? [y/N]:
```

---

## Architecture & Design Notes

```
cmd/arcade/main.go          Application entry point, CLI flag dispatcher, raw mode teardown, signals
internal/
├── engine/                 Shared engine core (zero game rules)
│   ├── difficulty.go       Difficulty enum (Easy/Medium/Hard) and profile interface
│   ├── entity.go           Position, Sprite, Rect (AABB collision primitives)
│   ├── game.go             Universal Game & MetricsProvider interface contracts
│   ├── input.go            Abstract Action translations (Arrows + WASD)
│   ├── loop.go             Fixed-ticker loop, input buffering, drop-and-continue pacing
│   ├── screen.go           tcell Screen wrapper with box-drawing & text centering
│   ├── theme.go            Theme struct & ThemeRegistry (Retro Green, Neon, Monochrome)
│   └── engine_test.go      Engine unit tests
├── history/                Decoupled persistence layer (zero tcell/rendering imports)
│   ├── achievements.go     10+ badge definitions, evaluator, atomic achievements.json store
│   ├── export.go           JSON and CSV data serialisation
│   ├── heatmap.go          30-day activity density calculation & ASCII rendering
│   ├── record.go           Record struct, personal best queries, relative time formatting
│   ├── reset.go            Safe deletion of history and badges files
│   ├── store.go            Atomic JSON persistence & self-healing corruption recovery
│   ├── streaks.go          Consecutive day streak calculator with clock injection
│   └── history_test.go     Comprehensive headless persistence & aggregation tests
├── menu/                   Menu and statistics screens (built on engine.Screen)
│   ├── history.go          Session log table, achievements tab, heatmap & streaks panel
│   ├── mainmenu.go         Title screen, streak banner, game selector
│   └── picker.go           Difficulty & Theme picker with personal best indicators
└── games/
    ├── snake/              Nokia Snake with near-miss tracking & tests
    ├── pacman/             Maze chase with dots, pellets, ghost AI & tests
    └── ballplate/          Breakout with angular deflection physics, rallies & tests
```

### Key Engineering Decisions:

1. **Engine / Game Decoupling**:
   - `internal/engine` owns zero game rules.
   - Every game implements the identical `engine.Game` contract:
     ```go
     type Game interface {
         Init(cfg GameConfig)
         Tick() TickResult
         HandleInput(a Action)
         Render(s *Screen)
         IsOver() bool
         Score() int
         Name() string
     }
     ```
   - Telemetry is supplied via the optional `engine.MetricsProvider` (`Metrics() map[string]int`), preserving contract cleanliness while enabling deep per-game stats.
2. **Atomic Writes & Corruption Recovery**:
   - History and achievement writes serialize data and write to `<file>.tmp` with file permissions `0600`.
   - Files are then atomically replaced via `os.Rename`. A power failure or crash mid-write cannot corrupt existing history.
   - If `history.json` or `achievements.json` is corrupted externally, the loader backs up the file to `<file>.corrupt-<timestamp>` and initializes a clean store rather than crashing the application.
3. **Cross-Platform Compatibility**:
   - **Path Separators**: All disk paths use `filepath.Join` and resolve against standard `os.UserConfigDir()`.
   - **Signals**: Windows does not support POSIX `SIGTERM`. The suite listens exclusively to `os.Interrupt` (portable across Windows, Linux, and macOS) to ensure `screen.Raw.Fini()` is invoked prior to termination.
   - **CGO Avoidance**: Zero CGO dependencies. Standard Go cross-compilation (`GOOS=... GOARCH=... go build`) works cleanly from any host OS.
4. **Smooth Frame Pacing & Terminal Timing**:
   - Simulation steps run on a strict `time.Ticker` rather than naive `time.Sleep` loops.
   - **Input Buffering**: Directional inputs are drained per tick and collapsed to at most the latest player intent, preventing multi-turn stalls.
   - **Drop-and-Continue Strategy**: If rendering takes longer than the tick interval (e.g. on very slow or loaded terminal emulators), lagged ticks are drained and skipped rather than fast-forwarding jerkily.
   - **No Redraw Flicker**: Incremental frames use `tcell.Screen.Show()`. Full repaint (`Sync()`) is reserved strictly for terminal resize events.
5. **Sub-Cell Motion in Ball & Plate**:
   - Ball coordinates ($X, Y$) and velocities ($V_x, V_y$) are double-precision floats (`float64`).
   - Impact angle reflection off the paddle uses relative deflection:
     $$\text{offset} = \frac{X_{\text{hit}} - X_{\text{center}}}{W_{\text{paddle}} / 2} \in [-1.0, 1.0]$$
     $$V_x = \text{speed} \cdot \sin(\text{offset} \cdot 60^\circ)$$
     $$V_y = -\left|\text{speed} \cdot \cos(\text{offset} \cdot 60^\circ)\right|$$

---

## Running Headless Tests

All game simulation logic, ghost AI decisions, physics math, and persistence systems are 100% testable without a terminal emulator or GUI:

```bash
go test -v ./...
go vet ./...
```

# Go Arcade — Terminal Game Suite (Snake, Pacman, Ball & Plate)

[![Latest Release](https://img.shields.io/github/v/release/Codexia-afk/Terminal-Arcade?logo=github&color=00ADD8)](https://github.com/Codexia-afk/Terminal-Arcade/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/Codexia-afk/Terminal-Arcade.svg)](https://pkg.go.dev/github.com/Codexia-afk/Terminal-Arcade)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Downloads](https://img.shields.io/github/downloads/Codexia-afk/Terminal-Arcade/total?color=green&logo=github)](https://github.com/Codexia-afk/Terminal-Arcade/releases)
[![Offline](https://img.shields.io/badge/network-100%25%20offline-success)](README.md)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)](README.md)

A production-quality, fully offline, terminal-native arcade game suite written in modern Go. Launches into an interactive selection menu and runs classic arcade games — **Snake (Nokia-style)**, **Pacman (Arcade Maze Chase)**, and **Ball & Plate (Breakout)** — entirely inside your terminal using direct raw-mode screen buffers, box-drawing chrome, multiple visual themes, concrete mechanical difficulty tiers, and atomic local persistence with deep statistics tracking.

---

## Table of Contents

- [⚡ Quick Start at a Glance](#-quick-start-at-a-glance)
- [🍎 macOS Guide](#-macos-guide)
- [🐧 Linux Guide](#-linux-guide)
- [🪟 Windows Guide](#-windows-guide)
- [📦 Go Developers (Universal `go install`)](#-go-developers-universal-go-install)
- [🎮 How to Play & Controls Reference](#-how-to-play--controls-reference)
  - [Main Menu & Navigation](#main-menu--navigation)
  - [Game Picker & Configuration](#game-picker--configuration)
  - [🐍 Snake Gameplay & Controls](#-snake-gameplay--controls)
  - [🟡 Pacman Gameplay & Controls](#-pacman-gameplay--controls)
  - [🧱 Ball & Plate (Breakout) Gameplay & Controls](#-ball--plate-breakout-gameplay--controls)
  - [📊 History, Badges & Activity Heatmap](#-history-badges--activity-heatmap)
- [⚙️ CLI Flags & Tools](#️-cli-flags--tools)
- [🎯 Difficulty Parameters](#-difficulty-parameters)
- [🎨 Visual Themes](#-visual-themes)
- [🏆 Deep Stats, Streaks & Achievement Badges](#-deep-stats-streaks--achievement-badges)
- [🧪 Headless Tests & Cross-Compilation](#-headless-tests--cross-compilation)
- [🏛️ Architecture & Design Notes](#️-architecture--design-notes)
- [🧹 Uninstallation & Data Locations](#-uninstallation--data-locations)

---

## ⚡ Quick Start at a Glance

Choose your platform to install with a single command, then simply type **`arcade`** to play:

| Platform | One-Line Installation Command | Launch Command |
| :--- | :--- | :--- |
| **macOS** | `curl -fsSL https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.sh \| bash` | `arcade` |
| **Linux** | `curl -fsSL https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.sh \| bash` | `arcade` |
| **Windows** | `irm https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.ps1 \| iex` | `arcade` |
| **Any OS (Go)** | `go install github.com/Codexia-afk/Terminal-Arcade/cmd/arcade@latest` | `arcade` |

---

## 🍎 macOS Guide

Supports both **Apple Silicon (M1, M2, M3, M4)** and **Intel-based Macs**.

### 1. One-Line Auto Install (Recommended)
Open **Terminal**, **iTerm2**, **Alacritty**, or **kitty** and run:

```bash
curl -fsSL https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.sh | bash
```

This script automatically:
1. Detects whether your Mac is Apple Silicon (`arm64`) or Intel (`amd64`).
2. Installs the executable binary into `/usr/local/bin` or `~/.local/bin`.
3. Adds `~/.local/bin` to your `~/.zshrc` / `~/.bash_profile` if not already present.
4. Verifies the installation with a test check.

### 2. Launch & Play on macOS
Open any terminal window and type:

```bash
arcade
```

*(If you just installed it, either run `source ~/.zshrc` or open a new Terminal tab).*

### 3. Build from Source on macOS
If you have Git and Go (1.22+) installed:

```bash
git clone https://github.com/Codexia-afk/Terminal-Arcade.git
cd Terminal-Arcade
make install   # Compiles and installs to /usr/local/bin/arcade
arcade         # Launch game
```

To run locally without system installation:
```bash
go build -buildvcs=false -o arcade ./cmd/arcade
./arcade
```

### 4. Running Pre-Compiled Release Binaries (Gatekeeper Note)
If you download a binary directly from GitHub Releases or browser:
```bash
# For Apple Silicon (M1/M2/M3/M4):
chmod +x arcade_darwin_arm64
xattr -d com.apple.quarantine arcade_darwin_arm64
./arcade_darwin_arm64

# For Intel Mac:
chmod +x arcade_darwin_amd64
xattr -d com.apple.quarantine arcade_darwin_amd64
./arcade_darwin_amd64
```
*(Removing the `com.apple.quarantine` attribute avoids the macOS "unidentified developer" pop-up).*

### 5. macOS Terminal Settings & Tips
- Recommended terminal: Default macOS **Terminal.app**, **iTerm2**, **Ghostty**, or **kitty**.
- Encoding: Ensure terminal encoding is set to **UTF-8** (Terminal Preferences → Profiles → Advanced → Text encoding: Unicode UTF-8).
- Minimum recommended terminal window size: **80 columns × 24 rows**.

---

## 🐧 Linux Guide

Supports all major Linux distributions (**Ubuntu, Debian, Fedora, Arch Linux, openSUSE, Alpine, CentOS/RHEL**).

### 1. One-Line Auto Install (Recommended)
Open your terminal and execute:

```bash
curl -fsSL https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.sh | bash
```

This will:
1. Detect your CPU architecture (`x86_64` / `amd64` or `aarch64` / `arm64`).
2. Install the binary into `/usr/local/bin` (or `~/.local/bin` if running without sudo).
3. Automatically configure your shell profile (`~/.bashrc`, `~/.zshrc`, or `~/.profile`).

### 2. Launch & Play on Linux
```bash
arcade
```

### 3. Build from Source on Linux
```bash
# Clone the repository
git clone https://github.com/Codexia-afk/Terminal-Arcade.git
cd Terminal-Arcade

# Compile and install to /usr/local/bin
make install

# Launch
arcade
```

### 4. Running Pre-Built Binary Manually
```bash
# 64-bit x86 / AMD:
chmod +x dist/arcade_linux_amd64
./dist/arcade_linux_amd64

# 64-bit ARM (Raspberry Pi 4/5, AWS Graviton):
chmod +x dist/arcade_linux_arm64
./dist/arcade_linux_arm64
```

### 5. Linux Terminal Compatibility
- Works seamlessly across `GNOME Terminal`, `Konsole`, `Alacritty`, `kitty`, `st`, `xterm`, and Linux virtual consoles (`tty1`–`tty6`).
- Under minimalist terminal environments or serial consoles with limited UTF-8 support, select the in-game **Monochrome / High-Contrast** theme which uses pure 7-bit ASCII characters.

---

## 🪟 Windows Guide

Supports **Windows 10** and **Windows 11** on 64-bit systems.

### 1. One-Line PowerShell Install (Recommended)
Open **PowerShell** (no Administrator rights required) and paste:

```powershell
irm https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.ps1 | iex
```

This will:
1. Download the latest `arcade_windows_amd64.exe`.
2. Save it to `$HOME\.local\bin\arcade.exe`.
3. Add `$HOME\.local\bin` to your User `PATH` environment variable permanently.

### 2. Launch & Play on Windows
Open a new **PowerShell** or **Windows Terminal** window and type:

```powershell
arcade
```

### 3. Build from Source on Windows
If you have Git and Go installed:

```cmd
git clone https://github.com/Codexia-afk/Terminal-Arcade.git
cd Terminal-Arcade
go build -buildvcs=false -o arcade.exe .\cmd\arcade
.\arcade.exe
```

### 4. Windows Terminal Recommendations
- **Recommended**: Run inside **[Windows Terminal](https://aka.ms/terminal)** with PowerShell for optimal box-drawing character rendering and truecolor arcade themes.
- **Legacy CMD (`cmd.exe`)**:
  If using classic `cmd.exe`, enable UTF-8 character encoding before launching:
  ```cmd
  chcp 65001
  arcade.exe
  ```
  *(Or choose the **Monochrome** theme in the game menu, which relies exclusively on standard ASCII characters like `#`, `=`, `o`, and `*`).*

---

## 📦 Go Developers (Universal `go install`)

If you have Go 1.22+ installed on any operating system, install directly using the standard Go package manager:

```bash
go install github.com/Codexia-afk/Terminal-Arcade/cmd/arcade@latest
```

This downloads, compiles, and installs `arcade` to your `$(go env GOPATH)/bin` directory.

Ensure `$(go env GOPATH)/bin` is in your `PATH`:
- **macOS / Linux**:
  ```bash
  echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.zshrc && source ~/.zshrc
  ```
- **Windows (PowerShell)**:
  ```powershell
  $goBin = "$(go env GOPATH)\bin"
  if ($env:Path -notlike "*$goBin*") { [Environment]::SetEnvironmentVariable("Path", $env:Path + ";$goBin", "User") }
  ```

---

## 🎮 How to Play & Controls Reference

### Main Menu & Navigation

When you launch `arcade`, you are greeted by the arcade title screen and daily streak banner:

| Key | Action |
| :--- | :--- |
| `↑` / `↓` or `W` / `S` | Navigate menu items (*Play Game*, *History & Badges*, *Theme*, *Export Data*, *Reset Data*, *Quit*) |
| `Enter` or `Space` | Select / confirm highlighted item |
| `1` – `6` | Quick select item by number |
| `Q` or `Esc` | Quit the arcade suite |

---

### Game Picker & Configuration

Upon selecting **Play Game**, configure your session:

| Key | Action |
| :--- | :--- |
| `←` / `→` or `A` / `D` | Cycle game: **Snake 🐍** ↔ **Pacman 🟡** ↔ **Ball & Plate 🧱** |
| `↑` / `↓` or `W` / `S` | Change difficulty: **Easy** ↔ **Medium** ↔ **Hard** |
| `Enter` | Launch selected game |
| `Esc` | Return to Main Menu |

---

### 🐍 Snake Gameplay & Controls

Control the classic hungry serpent, eat pellets to grow, and build personal best lengths.

```
┌──────────────────────────────────────┐
│  Score: 80   Length: 9   Best: 140   │
│                                      │
│               ████                   │
│                  █      ◆            │
│                  █████>              │
│                                      │
└──────────────────────────────────────┘
```

| Key | Action |
| :--- | :--- |
| `↑` `↓` `←` `→` or `W` `A` `S` `D` | Steer Snake up, down, left, right |
| `P` | Pause / Resume session |
| `Esc` or `Q` | End game & atomically save score and metrics |

#### Mechanics & Difficulty:
- **Easy**: Slower pace; **wraps around arena borders** safely. Great for casual play.
- **Medium**: Classic arcade speed; arena border walls are **fatal**.
- **Hard**: Blazing speed, fatal walls, and **lethal obstacles** periodically materialize on the field.
- **Telemetry**: Tracks near-misses (turns made 1 cell before fatal collision), food eaten, and survival ticks.

---

### 🟡 Pacman Gameplay & Controls

Navigate the maze corridors, gobble pellets, and evade or hunt the ghosts.

```
┌──────────────────────────────────────┐
│  Score: 320   Lives: ♥ ♥ ♥           │
│  ╔════════════════════════════════╗  │
│  ║ · · · · · · ║    ║ · · · · · · ║  │
│  ║ · ╔══════╗ · ║    ║ · ╔══════╗ · ║  │
│  ║ O ║      ║ · ╚════╝ · ║      ║ O ║  │
│  ║ · ╚══════╝ · · M    · ╚══════╝ · ║  │
│  ║ · · · · · · · · > · · · · · · · ║  │
│  ╚════════════════════════════════╝  │
└──────────────────────────────────────┘
```

| Key | Action |
| :--- | :--- |
| `↑` `↓` `←` `→` or `W` `A` `S` `D` | Buffer movement direction (Pacman glides continuously along paths) |
| `P` | Pause / Resume session |
| `Esc` or `Q` | End game & save session record |

#### Mechanics & Ghost AI:
- **Dots (`·`)**: 10 points each. Clear all dots to win the round.
- **Power Pellets (`O`)**: 50 points. Turns ghosts vulnerable (`w`) for a limited duration.
- **Ghost Personalities**:
  - 🔴 **Blinky**: Direct greedy chaser (Manhattan path targeting).
  - 🌸 **Pinky**: Ambush strategist (targets 4 steps ahead of Pacman).
  - 🔷 **Inky**: Flanker and unpredictable pursuer.
  - 🟠 **Clyde**: Casual wanderer.
- **Telemetry**: Tracks dots eaten, power pellets used, ghosts eaten, and levels cleared.

---

### 🧱 Ball & Plate (Breakout) Gameplay & Controls

Deflect the bouncing ball, demolish rows of bricks, and maintain long rallies.

```
┌──────────────────────────────────────┐
│  Score: 240   Lives: ♥ ♥   Rally: 7  │
│                                      │
│   ▀▀▀▀ ▀▀▀▀ ▀▀▀▀ ▀▀▀▀ ▀▀▀▀ ▀▀▀▀      │
│   ▀▀▀▀ ▀▀▀▀ ▀▀▀▀ ▀▀▀▀ ▀▀▀▀ ▀▀▀▀      │
│   #### #### #### #### #### ####      │
│                                      │
│                  ●                   │
│                                      │
│              ━━━━━━━                 │
└──────────────────────────────────────┘
```

| Key | Action |
| :--- | :--- |
| `←` / `→` or `A` / `D` | Move plate left / right |
| `Space` or `↑` / `W` | Launch ball from plate at start of life |
| `P` | Pause / Resume session |
| `Esc` or `Q` | End game & save session record |

#### Mechanics & Physics:
- **Sub-Cell Continuous Physics**: Ball trajectory is computed with double-precision floating-point precision ($X, Y, V_x, V_y$).
- **Dynamic Deflection Angle**: Hitting the edges of the paddle angles the ball sharply ($\pm 60^\circ$); hitting the center sends it upwards.
- **Bricks**:
  - Regular bricks (`▀` / `=`): 1 hit to break (10 points).
  - Tough bricks (`#`): 2 hits to break (25 points).
- **Telemetry**: Tracks bricks broken, max ball speed reached, plate bounces, and longest unbroken rally.

---

### 📊 History, Badges & Activity Heatmap

Access **History & Badges** from the Main Menu to review your arcade career:

| Key | Action |
| :--- | :--- |
| `←` / `→` or `A` / `D` | Switch tabs: **[1] History & Heatmap** ↔ **[2] Badges & Achievements** |
| `1` / `2` | Jump directly to Tab 1 or Tab 2 |
| `↑` / `↓` or `W` / `S` | Scroll session log rows or badges list |
| `Esc` or `Q` | Return to Main Menu |

- **Tab 1 (History & Heatmap)**: Displays consecutive play streaks, personal best scores per game, recent session table, and an ASCII **30-day activity density heatmap** (`·` = 0, `▪` = 1-2, `█` = 3+ sessions).
- **Tab 2 (Badges & Achievements)**: Lists all 10 unlockable badges with live percentage progress bars (`[██████░░░░] 60%`) and unlock timestamps.

---

## ⚙️ CLI Flags & Tools

The `arcade` binary includes built-in command line flags for data export and management without opening the GUI:

```bash
# Display full CLI help
arcade -h

# Export all gameplay history, stats, streaks, and badges to JSON
arcade --export --format=json --output=arcade_stats.json

# Export all session records to CSV (compatible with Excel, Sheets, Numbers)
arcade --export --format=csv --output=arcade_sessions.csv

# View exported files
cat arcade_stats.json

# Reset all persistent history, streak data, and unlocked achievements (with safety confirmation prompt)
arcade --reset-data
```

---

## 🎯 Difficulty Parameters

Each game features mechanically distinct parameters:

### Snake (Nokia-style)
| Parameter | Easy | Medium | Hard |
| :--- | :--- | :--- | :--- |
| **Tick Interval** | 180 ms | 130 ms | 90 ms |
| **Arena Dimensions** | 40 × 20 cells | 34 × 18 cells | 28 × 16 cells |
| **Wall Collisions** | **Wraps around edges** | **Fatal** | **Fatal** |
| **Deadly Obstacles** | None | None | **Spawns lethal obstacles** |

### Pacman (Arcade Maze Chase)
| Parameter | Easy | Medium | Hard |
| :--- | :--- | :--- | :--- |
| **Ghost Count** | 2 ghosts (Blinky, Pinky) | 3 ghosts (Blinky, Pinky, Inky) | 4 ghosts (All 4 ghosts) |
| **Ghost AI** | Random wander | Greedy Manhattan chase | Greedy + **Directional Ambush** |
| **Pellet Duration**| 50 ticks (~7.0s) | 35 ticks (~4.4s) | 20 ticks (~2.2s) |
| **Lives** | 3 lives | 3 lives | 2 lives |
| **Tick Interval** | 140 ms | 125 ms | 110 ms |

### Ball & Plate (Breakout)
| Parameter | Easy | Medium | Hard |
| :--- | :--- | :--- | :--- |
| **Paddle Width** | 7 cells (wide) | 5 cells (standard) | 3 cells (narrow) |
| **Base Speed** | 0.32 cells/tick | 0.40 cells/tick | 0.48 cells/tick |
| **Speed Scaling** | Constant | +0.002 cells/tick per brick | +0.004 cells/tick per brick |
| **Brick Grid** | Checkerboard gaps | Solid full grid | Solid grid + **Tough 2-hit bricks** |
| **Lives** | 4 lives | 3 lives | 2 lives |
| **Tick Interval** | 40 ms | 38 ms | 35 ms |

---

## 🎨 Visual Themes

Customize your visual experience from the Main Menu (**Option 3: Visual Theme**):

| Theme | Aesthetic | Foreground / Accents | Terminal Compatibility |
| :--- | :--- | :--- | :--- |
| **Retro Green** | Nokia LCD / Classic Game Boy | Lime / Olive on Deep Black | All ANSI terminals |
| **Neon** | Cyberpunk Arcade Cabinet | Cyan, Magenta, Yellow, Orange | Modern 256-color & Truecolor terminals |
| **Monochrome** | High-Contrast Pure ASCII | High-contrast White on Black | Guaranteed on legacy & serial consoles |

> [!TIP]
> Under Windows Command Prompt (`cmd.exe`) or minimal SSH sessions, the **Monochrome** theme uses only 7-bit ASCII characters (`#`, `=`, `o`, `*`) for maximum clarity and zero character glitches.

---

## 🏆 Deep Stats, Streaks & Achievement Badges

### 10 Unlockable Achievements

| ID | Title | Requirement | Unlock Hint |
| :--- | :--- | :--- | :--- |
| `first_blood` | **First Blood** | Complete 1 session of any game | Finish or exit your very first arcade round. |
| `century_club` | **Century Club** | Score 100+ in Snake | Eat at least 10 pellets in a single Snake game. |
| `ghost_hunter` | **Ghost Hunter** | Eat 20 ghosts across Pacman games | Hunt frightened ghosts during power pellet periods. |
| `brick_breaker` | **Brick Breaker** | Clear a full Breakout level on Hard | Demolish every brick including tough bricks on Hard. |
| `marathon` | **Marathon** | Single session lasting 5+ minutes | Keep playing continuously in any single round for 300s. |
| `dedicated` | **Dedicated** | Maintain a 7-day play streak | Play at least one game every day for 7 consecutive days. |
| `completionist`| **Completionist** | Play all 3 games at least once | Play Snake, Pacman, and Ball & Plate. |
| `perfectionist`| **Perfectionist** | Clear Pacman without losing a life | Eat every dot and clear the maze with 0 deaths. |
| `speed_demon`  | **Speed Demon** | Survive 60+ seconds on Snake Hard | Navigate fatal walls and obstacles on Hard for >60s. |
| `night_owl`    | **Night Owl** | Play between midnight and 4:00 AM | Play an arcade round between 00:00 and 04:00. |

---

## 🧪 Headless Tests & Cross-Compilation

### Running Automated Tests
All game rules, collisions, ghost AI, physics, and persistence engines are verified with headless unit tests:

```bash
# Run all unit tests headlessly
go test -v -count=1 ./...

# Run static analysis
go vet ./...
```

### Cross-Compiling for All Platforms
Build pre-compiled release binaries for macOS, Linux, and Windows from any machine:

```bash
make release
# or
./scripts/build_all.sh

# Inspect built artifacts
ls -lh dist/
# ./dist/arcade_darwin_arm64      (macOS Apple Silicon)
# ./dist/arcade_darwin_amd64      (macOS Intel)
# ./dist/arcade_linux_amd64       (Linux x86_64)
# ./dist/arcade_linux_arm64       (Linux ARM64)
# ./dist/arcade_windows_amd64.exe (Windows 64-bit)
```

---

## 🏛️ Architecture & Design Notes

```
cmd/arcade/main.go          Application entry point, CLI flags, signal handling, panic recovery
install.sh                  Universal one-line installer for macOS & Linux
install.ps1                 Universal one-line installer for Windows PowerShell
internal/
├── engine/                 Shared game engine (zero game-specific rules)
│   ├── difficulty.go       Difficulty tier types and contracts
│   ├── entity.go           Position, Sprite, Rect (AABB collision primitives)
│   ├── game.go             Universal Game & MetricsProvider interfaces
│   ├── input.go            Action mappings (Arrow keys + WASD)
│   ├── loop.go             Fixed-ticker loop, drop-and-continue pacing, input buffering
│   ├── screen.go           tcell Screen abstraction, border drawing, centered text
│   └── theme.go            Theme palettes (Retro Green, Neon, Monochrome)
├── history/                Offline persistence layer (independent of engine/tcell)
│   ├── achievements.go     10 badge definitions, evaluator, atomic JSON persistence
│   ├── export.go           JSON and CSV export serializers
│   ├── heatmap.go          30-day activity density calculation & ASCII renderer
│   ├── record.go           Session records, personal best aggregation, relative time formatting
│   ├── reset.go            Safe history wiping utilities
│   ├── store.go            Atomic file writes (.tmp + rename) & corruption recovery
│   └── streaks.go          Consecutive day streak engine with clock injection
├── menu/                   Terminal-native menu presentation
│   ├── history.go          Session history browser, badges progress tab, heatmap panel
│   ├── mainmenu.go         Interactive title menu and streak banner
│   └── picker.go           Game, difficulty, and theme selector with personal best display
└── games/
    ├── snake/              Nokia-style Snake with near-miss calculations
    ├── pacman/             Pacman maze chase with 4 ghost AI personalities
    └── ballplate/          Ball & Plate (Breakout) with sub-cell deflection physics
```

---

## 🧹 Uninstallation & Data Locations

### Storage Directories
Your data is stored locally in standard OS configuration directories:
- **macOS**: `~/Library/Application Support/goarcade/`
- **Linux**: `~/.config/goarcade/`
- **Windows**: `%APPDATA%\goarcade\`

### Remove Data Only
```bash
# Using CLI
arcade --reset-data

# Or manually:
# macOS:
rm -rf "$HOME/Library/Application Support/goarcade"
# Linux:
rm -rf "$HOME/.config/goarcade"
# Windows (PowerShell):
Remove-Item -Recurse -Force "$env:APPDATA\goarcade"
```

### Complete Uninstallation (Remove Binary)
```bash
# macOS / Linux (if installed via install.sh or make install):
sudo rm -f /usr/local/bin/arcade ~/.local/bin/arcade

# If installed via go install:
rm -f "$(go env GOPATH)/bin/arcade"

# Windows (PowerShell):
Remove-Item -Force "$HOME\.local\bin\arcade.exe"
```

---

## 📜 License

MIT License. Crafted with ❤️ for terminal gamers and Go developers.

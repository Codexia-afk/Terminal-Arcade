# Go Arcade — Terminal Arcade Game Suite (Enhanced Edition)

[![Latest Release](https://img.shields.io/github/v/release/Codexia-afk/Terminal-Arcade?logo=github&color=00ADD8)](https://github.com/Codexia-afk/Terminal-Arcade/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/Codexia-afk/Terminal-Arcade.svg)](https://pkg.go.dev/github.com/Codexia-afk/Terminal-Arcade)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Downloads](https://img.shields.io/github/downloads/Codexia-afk/Terminal-Arcade/total?color=green&logo=github)](https://github.com/Codexia-afk/Terminal-Arcade/releases)
[![Sponsor on Ko-fi](https://img.shields.io/badge/Ko--fi-Support-FF5E5B?logo=ko-fi&logoColor=white)](https://ko-fi.com/srinjoypramanick)
[![Offline](https://img.shields.io/badge/network-100%25%20offline-success)](README.md)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)](README.md)

```
  _______                  _             _                         _      
 |__   __|                (_)           | |      /\               | |     
    | | ___ _ __ _ __ ___  _ _ __   __ _| |     /  \   _ __ ___ __| | ___ 
    | |/ _ \ '__| '_ ` _ \| | '_ \ / _` | |    / /\ \ | '__/ __/ _` |/ _ \
    | |  __/ |  | | | | | | | | | | (_| | |   / ____ \| | | (_| (_| |  __/
    |_|\___|_|  |_| |_| |_|_|_| |_|\__,_|_|  /_/    \_\_|  \___\__,_|\___|
```

> **A production-quality, visually polished, truly engaging offline terminal arcade experience.**
> Features **5 distinct Snake gameplay modes** (Classic, Zen, Survival, Time Attack, Obstacle Challenge) × **3 difficulty tiers**, **Pacman** with 4 ghost AI personalities, **Ball & Plate** with continuous float64 physics, **5 visual themes**, **15 unlockable achievements**, **60 FPS decoupled simulation**, mouse navigation, and zero external runtime dependencies.

---

## ⚡ Quick Install (One-Line Commands)

Install with a single command and immediately launch by typing **`arcade`**:

### macOS & Linux
```bash
curl -fsSL https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.ps1 | iex
```

### Go Toolchain (`go install`)
```bash
go install github.com/Codexia-afk/Terminal-Arcade/cmd/arcade@latest
```

---

## 📑 Table of Contents

- [⚡ Quick Install](#-quick-install-one-line-commands)
- [✨ Design Philosophy](#-design-philosophy)
- [🎮 Complete Snake Modes Guide](#-complete-snake-modes-guide)
- [🟡 Pacman & 🧱 Ball & Plate Refinements](#-pacman---ball--plate-refinements)
- [🎨 Visual Theme Gallery](#-visual-theme-gallery)
- [🖥️ UI/UX Walkthrough](#️-uiux-walkthrough)
- [⌨️ Complete Keybindings & Mouse Reference](#️-complete-keybindings--mouse-reference)
- [🏆 15 Unlockable Achievements](#-15-unlockable-achievements)
- [🍎 macOS Guide](#-macos-guide)
- [🐧 Linux Guide](#-linux-guide)
- [🪟 Windows Guide](#-windows-guide)
- [⚙️ CLI Flags & Data Tools](#️-cli-flags--data-tools)
- [🧪 Headless Tests & Cross-Compilation](#-headless-tests--cross-compilation)
- [🏛️ Architecture & Code Layout](#️-architecture--code-layout)
- [🧹 Uninstallation & Data Locations](#-uninstallation--data-locations)
- [💖 Support & Sponsoring](#-support--sponsoring)
- [📜 License](#-license)

---

## ✨ Design Philosophy

Terminal Arcade is built around the ethos of **"Premium Terminal Gaming"** — combining the retro mechanical purity of Nokia 3310 classics with the responsive, zero-friction feel of modern arcade indie hits:

1. **Pixel-Perfect Centering & Breathing Room**: Every menu, pre-game confirmation screen, pause overlay, and HUD strip is dynamically centered within the terminal viewport. Double-box (`╔═╗║╚═╝`) borders, uniform padding, and responsive resizing ensure zero visual clutter or clipped elements.
2. **60 FPS Decoupled Game Loop**: The rendering pipeline runs at a steady 60 FPS (16.6ms intervals) with dirty cell tracking for tear-free rendering, while game logic advances on discrete, configurable tick intervals (80–180ms). Timer countdowns, particle pops, and animations remain buttery smooth without affecting gameplay tick pacing.
3. **2-Step Input Queue (Anti-Suicide Protection)**: Rapid key taps (e.g. `Right` then `Down` in a single tick) are buffered in a 2-step FIFO queue. The snake never reverses into its own neck due to fast-finger execution.
4. **Sub-Cell Collision & Smart Spawning**: Food and power pellets use Manhattan distance and forward trajectory prediction (`PredictNextCell`) to guarantee that food never spawns inside the snake's body or directly on the immediate next head cell.
5. **Zero CGO, Pure Go Standard Library**: No external C libraries, SDL, or ncurses. Pure Go via `tcell/v2` delivers instant compilation, identical cross-platform behavior, and zero dynamic linking issues.
6. **100% Offline & Private**: Zero network telemetry, zero phone-home pings, zero cloud tracking. Everything is saved locally via atomic JSON writes with automatic corruption recovery.

---

## 🎮 Complete Snake Modes Guide

Terminal Arcade features **5 distinct Snake gameplay modes**, each with unique win/loss conditions, mechanics, and 3 difficulty tiers (15 mode-difficulty combinations):

```
╔══════════════════════════════════════════════════════════════════════════════╗
║  Score: 140   Lives: ♥ ♥   Length: 12   Mode: Classic (Med)   Best: 280      ║
║                                                                              ║
║                 ████                                                         ║
║                    █         ★                                               ║
║                    ████████►                                                 ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

| Mode | Difficulty | Arena Size | Tick Interval | Lives | Core Rules & Special Mechanics |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Classic** | Easy | 40 × 20 | 120 ms | 3 | Nokia baseline. Screen edges **wrap around** safely. 10 pts/food. |
| | Medium | 34 × 18 | 100 ms | 2 | Border walls are **fatal**. Respawn at center on life loss. |
| | Hard | 28 × 16 | 80 ms | 1 | Fatal walls. 1 life (single collision ends session). Fast pacing. |
| **Zen** | Easy | 60 × 25 | 150 ms | ∞ | Relaxed, endless growth. Borders wrap continuously. **No death on self-collision**. |
| | Medium | 50 × 20 | 120 ms | ∞ | Wrap walls. Unbroken flow bonus: **+5 pts every 10s** without turning. |
| | Hard | 40 × 15 | 100 ms | ∞ | Tight arena, wrap walls. Live personal longest length tracker in HUD. |
| **Survival** | Easy | 38 × 18 | 130 ms | 3 | **5 Waves**. 1 hunter enemy chasing snake head via greedy Manhattan AI. |
| | Medium | 32 × 16 | 110 ms | 2 | 2 hunter enemies. 3s grace period on wave transition. +50 wave bonus. |
| | Hard | 28 × 14 | 90 ms | 1 | 3 relentless hunter enemies. Wave clear triggers victory overlay. |
| **Time Attack**| Easy | 40 × 20 | 110 ms | 1 | **90s countdown**. Timer HUD: Green (>30s) → Yellow (10–30s) → Red (<10s). |
| | Medium | 34 × 18 | 95 ms | 1 | **60s countdown**. Final 10s triggers **2-food frenzy** + terminal audio beep. |
| | Hard | 28 × 16 | 80 ms | 1 | **45s countdown**. Remaining time converts to speed bonus (`timeRemaining × 2`). |
| **Obstacle** | Easy | 38 × 18 | 120 ms | 2 | ~20% maze density. BFS connectivity check guarantees open path. |
| | Medium | 34 × 16 | 100 ms | 2 | ~35% maze density. 3-second pre-level planning countdown. |
| | Hard | 30 × 14 | 85 ms | 1 | ~50% maze density. Clear all food pellets for +100 maze clear bonus. |

---

## 🟡 Pacman & 🧱 Ball & Plate Refinements

### 🟡 Pacman (Arcade Maze Chase)
- **4 Distinct Ghost Personalities**:
  - 🔴 **Blinky (Chase)**: Direct greedy pursuer targeting Pacman's current coordinate.
  - 🌸 **Pinky (Ambush)**: Anticipates trajectory, targeting 4 steps ahead of Pacman's heading.
  - 🔷 **Inky (Roam)**: Strategic patrol AI that covers maze quadrants and flanks Pacman.
  - 🟠 **Clyde (Random)**: Independent wanderer who breaks off pursuit when within 8 cells.
- **In-HUD Ghost Legend**: Displays real-time ghost behavior (`Red: Chase | Pink: Ambush | Blue: Roam`).
- **Hollow Diamond (`◇`) Vulnerable State**: When a Power Pellet (`O`) is eaten, ghosts turn into flashing hollow diamonds.
- **Score Popups & Dot Flash**: Eating a ghost triggers a screen-wide white flash and a floating inline `+200` popup. Dot consumption triggers subtle micro-flashes.

### 🧱 Ball & Plate (Breakout)
- **Sub-Cell Continuous Physics**: Ball coordinates ($X, Y$) and velocity vectors ($V_x, V_y$) are modeled with 64-bit floating-point precision for smooth, continuous trajectories.
- **3 Paddle Deflection Zones**:
  - **Left Zone (33%)**: Deflects ball sharply leftward ($\approx -45^\circ$).
  - **Center Zone (34%)**: Deflects ball straight upward ($\approx 90^\circ$).
  - **Right Zone (33%)**: Deflects ball sharply rightward ($\approx +45^\circ$).
- **Color Inversion on Plate Contact**: Plate flashes inverse colors for 50ms upon ball contact.
- **Destruction Particles & In-Cell Score Popups**: Destroyed bricks flash white and display an in-cell score value (`+10` or `+25`) that fades over 300ms.
- **Live HUD Speed Gauge**: Displays real-time ball velocity scaling (`Speed: ▂▃▄`).

---

## 🎨 Visual Theme Gallery

Terminal Arcade includes **5 bespoke visual themes** switchable from the main menu or settings. Each theme features directional snake head runes (`▲`, `▼`, `◄`, `►`), distinct brick colors, and custom border chrome:

```
┌──────────────┬──────────────┬──────────────┬──────────────┬──────────────┐
│ Retro Green  │     Neon     │  Monochrome  │  Cyberpunk   │    Ocean     │
├──────────────┼──────────────┼──────────────┼──────────────┼──────────────┤
│ Wall: + - |  │ Wall: ╔ ═ ║  │ Wall: ▓ ▓ ▓  │ Wall: ┌ ─ │  │ Wall: ╭ ─ │  │
│ Head: ^ v < >│ Head: ▲ ▼ ◄ ►│ Head: ^ v < >│ Head: ▲ ▼ ◄ ►│ Head: ▲ ▼ ◄ ►│
│ Body: o      │ Body: O      │ Body: ■      │ Body: █      │ Body: ■      │
│ Food: *      │ Food: ♥      │ Food: ★      │ Food: ◆      │ Food: ●      │
│ Nokia LCD    │ Cyber Synth  │ Pure 7-Bit   │ Hot Neon     │ Deep Blue    │
└──────────────┴──────────────┴──────────────┴──────────────┴──────────────┘
```

| Theme | Aesthetic & Tone | Wall Style | Snake Glyphs (Head/Body/Food) | Terminal Compatibility |
| :--- | :--- | :--- | :--- | :--- |
| **Retro Green** | Nokia LCD / Classic Game Boy | `+ - \|` | `^ v < >` / `o` / `*` | Any ANSI / monochrome terminal |
| **Neon** | Synthwave 80s Arcade | `╔ ═ ║` | `▲ ▼ ◄ ►` / `O` / `♥` | 256-color & Truecolor terminals |
| **Monochrome** | High-Contrast Pure 7-Bit ASCII | `▓ ▓ ▓` | `^ v < >` / `■` / `★` | Guaranteed on serial & legacy consoles |
| **Cyberpunk** | High-Tech Dystopian Terminal | `┌ ─ │` | `▲ ▼ ◄ ►` / `█` / `◆` | Truecolor terminals (iTerm2, Alacritty) |
| **Ocean** | Calm Deep Aquatic Palette | `╭ ─ │` | `▲ ▼ ◄ ►` / `■` / `●` | Rounded UTF-8 terminals |

> [!TIP]
> Use the in-game **Theme Live Preview** (**Main Menu → 4: Theme Settings**) to view live animated swatches of each theme before playing!

---

## 🖥️ UI/UX Walkthrough

The interface is organized into 10 structured, keyboard- and mouse-accessible screens:

```
               ┌────────────────────────────────────────────────────────┐
               │                  1. Title Splash Screen                │
               └───────────────────────────┬────────────────────────────┘
                                           │ Press any key
                                           ▼
               ┌────────────────────────────────────────────────────────┐
               │                  2. Main Menu (8 Options)              │
               └───────┬───────────────────┬────────────────────┬───────┘
                       │ 1: Play           │ 3: Stats/Badges    │ 5: Settings
                       ▼                   ▼                    ▼
          ┌──────────────────────┐ ┌───────────────┐ ┌──────────────────┐
          │ 3. Game & Mode Picker│ │ 8. Dashboard  │ │ 10. Settings     │
          └────────────┬─────────┘ └───────────────┘ └──────────────────┘
                       │ Enter
                       ▼
          ┌──────────────────────┐
          │ 4. Pre-Game Confirm  │
          └────────────┬─────────┘
                       │ Enter
                       ▼
          ┌──────────────────────┐         Pause [P]   ┌────────────────┐
          │ 5. In-Game Arena     │ ──────────────────> │ 6. Pause Modal │
          └────────────┬─────────┘                     └────────────────┘
                       │ Win / Loss
                       ▼
          ┌──────────────────────┐
          │ 7. Victory / Over    │
          └──────────────────────┘
```

1. **Title Screen / Splash**: Clean ASCII "GO ARCADE" banner with daily streak counter, all-time score summary, and "Press any key to enter...".
2. **Main Menu**: Centered 8-row menu supporting arrow keys, numbers `1`–`8`, and **mouse clicks**.
3. **Game & Mode Picker**: Carousel allowing instant selection between Snake (Classic, Zen, Survival, Time Attack, Obstacle Challenge), Pacman, and Ball & Plate, with real-time difficulty parameter tables.
4. **Pre-Game Confirmation Screen**: Centered double-bordered panel displaying selected game mode, difficulty profile, personal best target, controls reference, and `[Enter] Start / [Esc] Back`.
5. **In-Game HUD & Arena**: Dynamic centered playfield with persistent top HUD strip showing streak badge, score, lives, near-miss indicator, and mode-specific gauges (timer / speed / enemy count).
6. **Pause Modal Overlay**: Centered translucent box freezing game simulation without clearing arena content.
7. **Victory / Game Over Overlay**: Displays final score, session duration, near-miss count, personal best comparison, unlocked badges, and `[Enter] Play Again / [Esc] Menu`.
8. **Theme Live Preview**: Dedicated screen rendering a live mini-arena displaying player, food, walls, and color swatches for each of the 5 themes.
9. **History & Badges Dashboard**: Tabbed interface featuring 30-day ASCII activity density heatmap (`·` = 0, `▪` = 1–2, `█` = 3+ games) and live percentage progress bars for all 15 achievements.
10. **Settings Screen**: Toggle audio bell alerts (`\a`), view 60 FPS refresh rate indicator, and browse full keybindings.

---

## ⌨️ Complete Keybindings & Mouse Reference

### Universal Navigation & Menus
| Key | Action |
| :--- | :--- |
| `↑` / `↓` or `W` / `S` | Navigate menu rows |
| `←` / `→` or `A` / `D` | Cycle game mode / difficulty / tabs |
| `Enter` or `Space` | Select / confirm highlighted item |
| `1` – `8` | Direct numeric item selection |
| `Left Mouse Click` | Click directly on any menu row to select |
| `Esc` or `Q` | Back to previous menu / cancel |

### In-Game Controls
| Game | Key | Action |
| :--- | :--- | :--- |
| **Snake (All Modes)** | `↑` `↓` `←` `→` / `W` `A` `S` `D` | Steer snake (buffered with 2-step input queue) |
| | `P` | Pause / resume session |
| | `Esc` or `Q` | End game & atomically persist metrics |
| **Pacman** | `↑` `↓` `←` `→` / `W` `A` `S` `D` | Buffer movement direction along maze corridors |
| | `P` | Pause / resume session |
| | `Esc` or `Q` | End game & save record |
| **Ball & Plate** | `←` / `→` or `A` / `D` | Move plate left / right |
| | `Space` or `↑` / `W` | Launch ball from plate |
| | `P` | Pause / resume session |
| | `Esc` or `Q` | End game & save record |

---

## 🏆 15 Unlockable Achievements

All achievements are evaluated deterministically and saved locally in your user profile:

| ID | Title | Requirement | Category |
| :--- | :--- | :--- | :--- |
| `first_blood` | **First Blood** | Complete 1 session of any game | General |
| `century_club` | **Century Club** | Score 100+ points in Snake | Snake |
| `ghost_hunter` | **Ghost Hunter** | Eat 20 ghosts across Pacman games | Pacman |
| `brick_breaker` | **Brick Breaker** | Clear a full Breakout level on Hard | Ball & Plate |
| `marathon` | **Marathon** | Single session lasting 5+ minutes (300s) | Endurance |
| `dedicated` | **Dedicated** | Maintain a 7-day play streak | Consistency |
| `completionist`| **Completionist** | Play all 3 games (Snake, Pacman, Breakout) | Exploration |
| `perfectionist`| **Perfectionist** | Clear Pacman without losing a single life | Mastery |
| `speed_demon` | **Speed Demon** | Survive 60+ seconds on Snake Hard | Agility |
| `night_owl` | **Night Owl** | Play an arcade round between midnight and 4:00 AM | Secret |
| `zen_master` | **Zen Master** | Reach a snake length of 30+ in Zen mode | Snake |
| `survivor_5` | **Wave Survivor** | Clear all 5 waves in Survival mode | Snake |
| `speed_runner` | **Speed Runner** | Score 150+ points in Time Attack mode | Snake |
| `maze_runner` | **Maze Runner** | Clear an Obstacle Challenge maze without dying | Snake |
| `arcade_legend` | **Arcade Legend** | Unlock 10 or more other achievements | Mastery |

---

## 🍎 macOS Guide

Supports both **Apple Silicon (M1, M2, M3, M4)** and **Intel-based Macs**.

### 1. One-Line Auto Install (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.sh | bash
```

### 2. Launch & Play
```bash
arcade
```

### 3. Build from Source
```bash
git clone https://github.com/Codexia-afk/Terminal-Arcade.git
cd Terminal-Arcade
make install
arcade
```

---

## 🐧 Linux Guide

Supports all major distributions (**Ubuntu, Debian, Fedora, Arch Linux, openSUSE, Alpine, CentOS/RHEL**).

### 1. One-Line Auto Install
```bash
curl -fsSL https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.sh | bash
```

### 2. Launch & Play
```bash
arcade
```

### 3. Build from Source
```bash
git clone https://github.com/Codexia-afk/Terminal-Arcade.git
cd Terminal-Arcade
make install
arcade
```

---

## 🪟 Windows Guide

Supports **Windows 10** and **Windows 11** (64-bit).

### 1. One-Line PowerShell Install (Recommended)
Open **PowerShell** and run:
```powershell
irm https://raw.githubusercontent.com/Codexia-afk/Terminal-Arcade/main/install.ps1 | iex
```

### 2. Launch & Play
Open Windows Terminal or PowerShell and type:
```powershell
arcade
```

### 3. Build from Source
```cmd
git clone https://github.com/Codexia-afk/Terminal-Arcade.git
cd Terminal-Arcade
go build -buildvcs=false -o arcade.exe .\cmd\arcade
.\arcade.exe
```

---

## ⚙️ CLI Flags & Data Tools

Manage your scores, export telemetry, or reset data directly from the command line:

```bash
# Display CLI help
arcade -h

# Export all gameplay history, streaks, and badges to JSON
arcade --export --format=json --output=arcade_stats.json

# Export all session records to CSV
arcade --export --format=csv --output=arcade_sessions.csv

# View exported statistics
cat arcade_stats.json

# Reset all persistent history and badges (with safety confirmation prompt)
arcade --reset-data
```

---

## 🧪 Headless Tests & Cross-Compilation

### Running Automated Tests
All game rules, collision algorithms, ghost AI personalities, physics, and persistence engines are verified with headless unit tests:

```bash
# Run all unit tests across all packages
make test
# or
go test -buildvcs=false -v ./...

# Run static analysis
make vet
```

### Cross-Compilation for All Platforms
Compile release binaries for all 5 target architectures and generate SHA256 checksums:

```bash
make cross-compile
# or
bash scripts/build_all.sh

# Verify output artifacts
ls -lh dist/
cat dist/checksums.txt
```

Generated release targets:
- `dist/arcade_darwin_arm64` (macOS Apple Silicon)
- `dist/arcade_darwin_amd64` (macOS Intel)
- `dist/arcade_linux_amd64` (Linux 64-bit x86)
- `dist/arcade_linux_arm64` (Linux 64-bit ARM)
- `dist/arcade_windows_amd64.exe` (Windows 64-bit)
- `dist/checksums.txt` (SHA256 checksums)

---

## 🏛️ Architecture & Code Layout

```
cmd/arcade/main.go               Application entry point, CLI flags, screen init, mouse activation
install.sh                       Universal one-line installer for macOS & Linux
install.ps1                      Universal one-line installer for Windows PowerShell
Makefile                         Build, test, vet, cross-compilation & checksum generation
scripts/build_all.sh             Standalone multi-target release build script
internal/
├── engine/                      Shared engine interfaces & rendering abstractions
│   ├── difficulty.go            Difficulty tier types and contracts
│   ├── entity.go                Position, Sprite, Rect (AABB collision primitives)
│   ├── game.go                  Universal Game & MetricsProvider interfaces
│   ├── input.go                 Action mappings (Arrows, WASD, Mouse, Pause)
│   ├── loop.go                  60 FPS decoupled loop with dirty tracking
│   ├── screen.go                tcell Screen abstraction, double-borders, audio toggle
│   └── theme.go                 5 themes with directional snake runes & palette definitions
├── history/                     Offline persistence & analytics
│   ├── achievements.go          15 achievement definitions & evaluator
│   ├── export.go                JSON and CSV data export serializers
│   ├── heatmap.go               30-day activity density calculation & ASCII renderer
│   ├── record.go                Session records & personal best aggregation
│   ├── reset.go                 Safe history reset utility
│   ├── store.go                 Atomic JSON file persistence (.tmp + rename)
│   └── streaks.go               Consecutive day streak engine with clock injection
├── menu/                        Menus, screens, and overlays
│   ├── history.go               Session browser, badges progress tab, heatmap panel
│   ├── mainmenu.go              Centered 8-item menu with mouse click support
│   ├── picker.go                Game carousel & Snake mode picker
│   ├── pregame.go               Pre-game double-bordered confirmation screen
│   ├── settings.go              Audio bell toggle, refresh indicator & keybindings
│   ├── themepreview.go          Live animated arcade theme preview screen
│   └── titlescreen.go           ASCII launch splash with streak badge & PB summary
└── games/
    ├── snake/                   Snake orchestrator & mechanics
    │   ├── game.go              Universal game wrapper delegating to active mode
    │   ├── mechanics.go         2-step input queue, Manhattan distance, next-cell prediction
    │   ├── modes.go             Canonical mode metadata & difficulty descriptors
    │   └── modes/               5 concrete Snake mode implementations
    │       ├── common.go        BaseMode shared HUD, food spawner, overlays
    │       ├── classic.go       Nokia baseline with lives & border wrap
    │       ├── zen.go           Relaxed endless growth with flow bonus
    │       ├── survival.go      5-wave hunter AI pursuit with grace periods
    │       ├── timeattack.go    Countdown timer with speed bonus & 2-food frenzy
    │       └── obstacle.go      Procedural maze with BFS connectivity verification
    ├── pacman/                  Pacman with 4 ghost personalities (Chase, Ambush, Roam, Random)
    └── ballplate/               Breakout with continuous float64 physics & 3-zone paddle
```

---

## 🧹 Uninstallation & Data Locations

### Storage Directories
Your data is stored locally in standard OS configuration directories:
- **macOS**: `~/Library/Application Support/goarcade/`
- **Linux**: `~/.config/goarcade/`
- **Windows**: `%APPDATA%\goarcade\`

### Reset Data
```bash
arcade --reset-data
```

### Complete Uninstallation
```bash
# macOS / Linux:
sudo rm -f /usr/local/bin/arcade ~/.local/bin/arcade

# Windows (PowerShell):
Remove-Item -Force "$HOME\.local\bin\arcade.exe"
```

---

## 💖 Support & Sponsoring

If you enjoy playing Terminal Arcade and want to support continued development:
- **☕ Buy a coffee on Ko-fi**: [ko-fi.com/srinjoypramanick](https://ko-fi.com/srinjoypramanick)
- **💖 Sponsor on GitHub**: [github.com/sponsors/Codexia-afk](https://github.com/sponsors/Codexia-afk)

---

## 📜 License

MIT License. Crafted with ❤️ for terminal gamers and Go developers.

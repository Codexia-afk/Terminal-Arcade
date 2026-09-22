// Package pacman implements the arcade Pacman maze chase game.
package pacman

import (
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// TileType identifies the content of a maze cell.
type TileType rune

const (
	TileEmpty  TileType = ' '
	TileWall   TileType = '#'
	TileDot    TileType = '.'
	TilePellet TileType = 'o'
	TileGate   TileType = '-'
)

var defaultMazeLayout = []string{
	"############################",
	"#o..........####..........o#",
	"#.####.####.####.####.####.#",
	"#..........................#",
	"#.####.##.########.##.####.#",
	"#......##....##....##......#",
	"######.##### ## #####.######",
	"     #.##### ## #####.#     ",
	"######.##          ##.######",
	"      .   ###--###   .      ",
	"######.## # G  G # ##.######",
	"     #.## # G  G # ##.#     ",
	"######.## ######## ##.######",
	"#............P.............#",
	"#.####.#####.##.#####.####.#",
	"#o..##................##..o#",
	"###.##.##.########.##.##.###",
	"#......##....##....##......#",
	"#.##########.##.##########.#",
	"#..........................#",
	"############################",
}

// Maze represents the parsed board state and dimensions.
type Maze struct {
	Width       int
	Height      int
	Tiles       [][]TileType
	PlayerStart engine.Position
	GhostStarts []engine.Position
	Remaining   int
}

// NewMaze parses the default layout into a playable maze instance.
func NewMaze() *Maze {
	h := len(defaultMazeLayout)
	w := len(defaultMazeLayout[0])
	m := &Maze{
		Width:       w,
		Height:      h,
		Tiles:       make([][]TileType, h),
		GhostStarts: make([]engine.Position, 0, 4),
	}

	for y := 0; y < h; y++ {
		row := defaultMazeLayout[y]
		m.Tiles[y] = make([]TileType, w)
		for x := 0; x < w; x++ {
			char := rune(row[x])
			switch char {
			case 'P':
				m.PlayerStart = engine.Position{X: x, Y: y}
				m.Tiles[y][x] = TileEmpty
			case 'G':
				m.GhostStarts = append(m.GhostStarts, engine.Position{X: x, Y: y})
				m.Tiles[y][x] = TileEmpty
			case '.':
				m.Tiles[y][x] = TileDot
				m.Remaining++
			case 'o':
				m.Tiles[y][x] = TilePellet
				m.Remaining++
			case '#':
				m.Tiles[y][x] = TileWall
			case '-':
				m.Tiles[y][x] = TileGate
			default:
				m.Tiles[y][x] = TileEmpty
			}
		}
	}
	return m
}

// IsWall reports whether the given cell is impassable by the player.
func (m *Maze) IsWall(pos engine.Position) bool {
	// Wrap-around tunnel check
	if pos.Y == 9 {
		if pos.X < 0 || pos.X >= m.Width {
			return false
		}
	}
	if pos.X < 0 || pos.X >= m.Width || pos.Y < 0 || pos.Y >= m.Height {
		return true
	}
	t := m.Tiles[pos.Y][pos.X]
	return t == TileWall || t == TileGate
}

// IsGhostPassable reports whether a ghost can enter the tile.
func (m *Maze) IsGhostPassable(pos engine.Position) bool {
	// Wrap-around tunnel check
	if pos.Y == 9 {
		if pos.X < 0 || pos.X >= m.Width {
			return true
		}
	}
	if pos.X < 0 || pos.X >= m.Width || pos.Y < 0 || pos.Y >= m.Height {
		return false
	}
	return m.Tiles[pos.Y][pos.X] != TileWall
}

// WrapPosition wraps positions around side tunnels if applicable.
func (m *Maze) WrapPosition(pos engine.Position) engine.Position {
	if pos.Y == 9 {
		if pos.X < 0 {
			pos.X = m.Width - 1
		} else if pos.X >= m.Width {
			pos.X = 0
		}
	}
	return pos
}

package engine

import "github.com/gdamore/tcell/v2"

// Theme defines the visual palette and glyph representations for game elements.
type Theme struct {
	Name        string
	Background  tcell.Color
	Wall        tcell.Color
	Player      tcell.Color
	Enemy       tcell.Color
	Item        tcell.Color
	HUD         tcell.Color
	Accent      tcell.Color
	Text        tcell.Color
	Vulnerable  tcell.Color
	BrickColors []tcell.Color

	WallGlyph          rune
	PlayerGlyph        rune
	EnemyGlyph         rune
	ItemGlyph          rune
	PelletGlyph        rune
	BallGlyph          rune
	PlateGlyph         rune
	BrickGlyph         rune
	ToughBrick         rune
	SnakeHeadGlyph     rune
	SnakeHeadUp        rune
	SnakeHeadDown      rune
	SnakeHeadLeft      rune
	SnakeHeadRight     rune
	SnakeBodyGlyph     rune
	SnakeFoodGlyph     rune
	SnakeObstacleGlyph rune
}

// Themes holds the built-in suite themes.
var Themes = map[string]Theme{
	"Retro Green": {
		Name:        "Retro Green",
		Background:  tcell.ColorBlack,
		Wall:        tcell.ColorDarkGreen,
		Player:      tcell.ColorGreen,
		Enemy:       tcell.ColorForestGreen,
		Item:        tcell.ColorLime,
		HUD:         tcell.ColorGreen,
		Accent:      tcell.ColorLime,
		Text:        tcell.ColorGreen,
		Vulnerable:  tcell.ColorDarkOliveGreen,
		BrickColors: []tcell.Color{tcell.ColorLime, tcell.ColorGreen, tcell.ColorDarkGreen, tcell.ColorForestGreen},

		WallGlyph:          '█',
		PlayerGlyph:        '█',
		EnemyGlyph:         'G',
		ItemGlyph:          '◆',
		PelletGlyph:        '●',
		BallGlyph:          '●',
		PlateGlyph:         '▬',
		BrickGlyph:         '▀',
		ToughBrick:         '▓',
		SnakeHeadGlyph:     '▲',
		SnakeHeadUp:        '▲',
		SnakeHeadDown:      '▼',
		SnakeHeadLeft:      '◄',
		SnakeHeadRight:     '►',
		SnakeBodyGlyph:     '█',
		SnakeFoodGlyph:     '◆',
		SnakeObstacleGlyph: 'X',
	},
	"Neon": {
		Name:       "Neon",
		Background: tcell.ColorBlack,
		Wall:       tcell.ColorBlueViolet,
		Player:     tcell.ColorYellow,
		Enemy:      tcell.ColorRed,
		Item:       tcell.ColorGold,
		HUD:        tcell.ColorAqua,
		Accent:     tcell.ColorFuchsia,
		Text:       tcell.ColorWhite,
		Vulnerable: tcell.ColorDeepSkyBlue,
		BrickColors: []tcell.Color{
			tcell.ColorRed,
			tcell.ColorOrangeRed,
			tcell.ColorYellow,
			tcell.ColorGreen,
			tcell.ColorDeepSkyBlue,
			tcell.ColorDarkMagenta,
		},

		WallGlyph:          '║',
		PlayerGlyph:        'C',
		EnemyGlyph:         'M',
		ItemGlyph:          '·',
		PelletGlyph:        '●',
		BallGlyph:          '●',
		PlateGlyph:         '━',
		BrickGlyph:         '█',
		ToughBrick:         '▒',
		SnakeHeadGlyph:     '▲',
		SnakeHeadUp:        '▲',
		SnakeHeadDown:      '▼',
		SnakeHeadLeft:      '◄',
		SnakeHeadRight:     '►',
		SnakeBodyGlyph:     '▪',
		SnakeFoodGlyph:     '☆',
		SnakeObstacleGlyph: '◆',
	},
	"Monochrome": {
		Name:        "Monochrome",
		Background:  tcell.ColorBlack,
		Wall:        tcell.ColorWhite,
		Player:      tcell.ColorWhite,
		Enemy:       tcell.ColorWhite,
		Item:        tcell.ColorWhite,
		HUD:         tcell.ColorWhite,
		Accent:      tcell.ColorWhite,
		Text:        tcell.ColorWhite,
		Vulnerable:  tcell.ColorGray,
		BrickColors: []tcell.Color{tcell.ColorWhite, tcell.ColorWhite, tcell.ColorWhite, tcell.ColorWhite},

		WallGlyph:          '#',
		PlayerGlyph:        'o',
		EnemyGlyph:         'M',
		ItemGlyph:          '*',
		PelletGlyph:        'o',
		BallGlyph:          'o',
		PlateGlyph:         '-',
		BrickGlyph:         '=',
		ToughBrick:         '#',
		SnakeHeadGlyph:     '^',
		SnakeHeadUp:        '^',
		SnakeHeadDown:      'v',
		SnakeHeadLeft:      '<',
		SnakeHeadRight:     '>',
		SnakeBodyGlyph:     '#',
		SnakeFoodGlyph:     '*',
		SnakeObstacleGlyph: 'X',
	},
	"Cyberpunk": {
		Name:       "Cyberpunk",
		Background: tcell.NewRGBColor(12, 10, 24),
		Wall:       tcell.ColorFuchsia,
		Player:     tcell.ColorAqua,
		Enemy:      tcell.ColorYellow,
		Item:       tcell.ColorDarkViolet,
		HUD:        tcell.ColorAqua,
		Accent:     tcell.ColorFuchsia,
		Text:       tcell.ColorWhite,
		Vulnerable: tcell.ColorTeal,
		BrickColors: []tcell.Color{
			tcell.ColorAqua,
			tcell.ColorFuchsia,
			tcell.ColorYellow,
			tcell.ColorDarkMagenta,
			tcell.ColorDeepPink,
		},

		WallGlyph:          '▓',
		PlayerGlyph:        '▲',
		EnemyGlyph:         'X',
		ItemGlyph:          '★',
		PelletGlyph:        '◆',
		BallGlyph:          '◆',
		PlateGlyph:         '═',
		BrickGlyph:         '█',
		ToughBrick:         '▒',
		SnakeHeadGlyph:     '▲',
		SnakeHeadUp:        '▲',
		SnakeHeadDown:      '▼',
		SnakeHeadLeft:      '≪',
		SnakeHeadRight:     '≫',
		SnakeBodyGlyph:     '▮',
		SnakeFoodGlyph:     '✦',
		SnakeObstacleGlyph: '◆',
	},
	"Ocean": {
		Name:       "Ocean",
		Background: tcell.NewRGBColor(4, 18, 36),
		Wall:       tcell.ColorTeal,
		Player:     tcell.ColorAqua,
		Enemy:      tcell.ColorDeepSkyBlue,
		Item:       tcell.ColorLightCyan,
		HUD:        tcell.ColorAquaMarine,
		Accent:     tcell.ColorTurquoise,
		Text:       tcell.ColorWhite,
		Vulnerable: tcell.ColorDarkSlateGray,
		BrickColors: []tcell.Color{
			tcell.ColorLightCyan,
			tcell.ColorAqua,
			tcell.ColorTeal,
			tcell.ColorSteelBlue,
			tcell.ColorDodgerBlue,
		},

		WallGlyph:          '║',
		PlayerGlyph:        '>',
		EnemyGlyph:         'S',
		ItemGlyph:          '*',
		PelletGlyph:        '○',
		BallGlyph:          '●',
		PlateGlyph:         '~',
		BrickGlyph:         '█',
		ToughBrick:         '▓',
		SnakeHeadGlyph:     '^',
		SnakeHeadUp:        '^',
		SnakeHeadDown:      'v',
		SnakeHeadLeft:      '<',
		SnakeHeadRight:     '>',
		SnakeBodyGlyph:     '~',
		SnakeFoodGlyph:     '✦',
		SnakeObstacleGlyph: '◆',
	},
}

// SnakeHeadForDir returns the directional head glyph based on current movement vector.
func (t Theme) SnakeHeadForDir(dir Position) rune {
	if dir.X == 1 && t.SnakeHeadRight != 0 {
		return t.SnakeHeadRight
	}
	if dir.X == -1 && t.SnakeHeadLeft != 0 {
		return t.SnakeHeadLeft
	}
	if dir.Y == -1 && t.SnakeHeadUp != 0 {
		return t.SnakeHeadUp
	}
	if dir.Y == 1 && t.SnakeHeadDown != 0 {
		return t.SnakeHeadDown
	}
	if t.SnakeHeadGlyph != 0 {
		return t.SnakeHeadGlyph
	}
	return '▲'
}

// ThemeNames returns the registered theme names in standard menu order.
func ThemeNames() []string {
	return []string{"Retro Green", "Neon", "Monochrome", "Cyberpunk", "Ocean"}
}

// GetTheme looks up a theme by name with a sensible fallback.
func GetTheme(name string) Theme {
	if t, ok := Themes[name]; ok {
		return t
	}
	return Themes["Retro Green"]
}

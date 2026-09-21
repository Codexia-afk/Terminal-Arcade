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

	WallGlyph   rune
	PlayerGlyph rune
	EnemyGlyph  rune
	ItemGlyph   rune
	PelletGlyph rune
	BallGlyph   rune
	PlateGlyph  rune
	BrickGlyph  rune
	ToughBrick  rune
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

		WallGlyph:   '█',
		PlayerGlyph: '█',
		EnemyGlyph:  'G',
		ItemGlyph:   '◆',
		PelletGlyph: '●',
		BallGlyph:   '●',
		PlateGlyph:  '▬',
		BrickGlyph:  '▀',
		ToughBrick:  '▓',
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

		WallGlyph:   '║',
		PlayerGlyph: 'C',
		EnemyGlyph:  'M',
		ItemGlyph:   '·',
		PelletGlyph: '●',
		BallGlyph:   '●',
		PlateGlyph:  '━',
		BrickGlyph:  '█',
		ToughBrick:  '▒',
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

		WallGlyph:   '#',
		PlayerGlyph: 'o',
		EnemyGlyph:  'M',
		ItemGlyph:   '*',
		PelletGlyph: 'o',
		BallGlyph:   'o',
		PlateGlyph:  '-',
		BrickGlyph:  '=',
		ToughBrick:  '#',
	},
}

// ThemeNames returns the registered theme names in standard menu order.
func ThemeNames() []string {
	return []string{"Retro Green", "Neon", "Monochrome"}
}

// GetTheme looks up a theme by name with a sensible fallback.
func GetTheme(name string) Theme {
	if t, ok := Themes[name]; ok {
		return t
	}
	return Themes["Retro Green"]
}

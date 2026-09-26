package menu

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// RenderThemePreview renders a live miniature arcade scene showcasing the theme's colors and glyphs.
func RenderThemePreview(s *engine.Screen, th engine.Theme, x, y, w, h int, gameName string) {
	// 1. Draw mini preview box
	title := fmt.Sprintf("PREVIEW: %s", th.Name)
	s.BoxWithTitle(x, y, w, h, title, th.Accent, th.Background, th.Player)

	// Interior boundaries
	innerX := x + 1
	innerY := y + 1
	innerW := w - 2
	innerH := h - 2

	if innerW < 18 || innerH < 6 {
		return
	}

	// 2. Fill interior with theme background
	s.Fill(innerX, innerY, innerW, innerH, ' ', th.Text, th.Background)

	// 3. Mini arena wall borders
	miniArenaX := innerX + 1
	miniArenaY := innerY + 1
	miniArenaW := innerW - 2
	miniArenaH := innerH - 4
	if miniArenaH < 4 {
		miniArenaH = 4
	}

	// Draw mini border
	for bx := 0; bx < miniArenaW; bx++ {
		s.DrawCell(miniArenaX+bx, miniArenaY, th.WallGlyph, th.Wall, th.Background)
		s.DrawCell(miniArenaX+bx, miniArenaY+miniArenaH-1, th.WallGlyph, th.Wall, th.Background)
	}
	for by := 0; by < miniArenaH; by++ {
		s.DrawCell(miniArenaX, miniArenaY+by, th.WallGlyph, th.Wall, th.Background)
		s.DrawCell(miniArenaX+miniArenaW-1, miniArenaY+by, th.WallGlyph, th.Wall, th.Background)
	}

	// 4. Draw Snake elements inside mini arena
	// Head + body
	headX := miniArenaX + miniArenaW/2
	headY := miniArenaY + miniArenaH/2
	s.DrawCell(headX, headY, th.SnakeHeadGlyph, th.Player, th.Background)
	s.DrawCell(headX-1, headY, th.SnakeBodyGlyph, th.Player, th.Background)
	s.DrawCell(headX-2, headY, th.SnakeBodyGlyph, th.Player, th.Background)
	s.DrawCell(headX-3, headY, th.SnakeBodyGlyph, th.Player, th.Background)

	// Food pellet
	foodX := headX + 3
	if foodX < miniArenaX+miniArenaW-1 {
		s.DrawCell(foodX, headY, th.SnakeFoodGlyph, th.Item, th.Background)
	}

	// Obstacle / Wall hazard
	obsX := miniArenaX + 2
	obsY := miniArenaY + 1
	if obsY < miniArenaY+miniArenaH-1 {
		s.DrawCell(obsX, obsY, th.SnakeObstacleGlyph, th.Accent, th.Background)
	}

	// 5. Draw Pacman & Ghost demo or Ball & Plate brick demo
	ghostColor := th.Enemy
	brickColor := tcell.ColorYellow
	if len(th.BrickColors) > 0 {
		brickColor = th.BrickColors[0]
	}

	rightDemoX := miniArenaX + miniArenaW - 3
	if rightDemoX > headX+1 {
		s.DrawCell(rightDemoX, miniArenaY+1, 'ᗣ', ghostColor, th.Background)
		s.DrawCell(rightDemoX, miniArenaY+2, th.BrickGlyph, brickColor, th.Background)
	}

	// 6. Color palette swatches at bottom
	swatchY := innerY + innerH - 2
	s.DrawText(innerX+2, swatchY, "Palette:", th.Text, th.Background)
	s.DrawCell(innerX+11, swatchY, '■', th.Player, th.Background)
	s.DrawCell(innerX+13, swatchY, '■', th.Accent, th.Background)
	s.DrawCell(innerX+15, swatchY, '■', th.Item, th.Background)
	s.DrawCell(innerX+17, swatchY, '■', th.Wall, th.Background)
	s.DrawCell(innerX+19, swatchY, '■', th.HUD, th.Background)


	glyphLabel := fmt.Sprintf("Head:%c Food:%c Obs:%c", th.SnakeHeadGlyph, th.SnakeFoodGlyph, th.SnakeObstacleGlyph)
	if len(glyphLabel) > innerW-2 {
		glyphLabel = glyphLabel[:innerW-2]
	}
	s.DrawText(innerX+2, swatchY+1, glyphLabel, th.HUD, th.Background)
}

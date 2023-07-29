package movement

import (
	"2d_game/helpers"
	"2d_game/player"
	"2d_game/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

func PlayerMovement(x int, y int) (int, int) {
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		x -= 2
		player.Player.Count++
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		y -= 2
		player.Player.Count++
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		x += 2
		player.Player.Count++
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		y += 2
		player.Player.Count++
	}

	// +1/-1 is to stop player before it reaches the border
	if x+player.FrameWidthCurrent >= settings.ScreenWidth-settings.Padding {
		x = settings.ScreenWidth - settings.Padding - player.FrameWidthCurrent
	}

	if x-player.FrameWidthCurrent <= settings.Padding {
		x = settings.Padding + player.FrameWidthCurrent
	}

	if y+player.FrameHeightCurrent >= settings.ScreenHeight-settings.Padding {
		y = settings.ScreenHeight - settings.Padding - player.FrameHeightCurrent
	}

	if y-player.FrameHeightCurrent <= settings.Padding {
		y = settings.Padding + player.FrameHeightCurrent
	}

	for _, p := range player.Player.ChildProjectTile {
		p.Count++
		p.X = p.X - math.Sin(p.CamAngle)*p.Speed
		p.Y = p.Y + math.Cos(p.CamAngle)*p.Speed
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		centerX := float64(settings.ScreenWidth / 2)
		centerY := float64(settings.ScreenHeight / 2)

		var line = helpers.Line{X1: centerX, Y1: centerY, X2: float64(mouseX), Y2: float64(mouseY)}
		var angle = line.Angle()
		var projectTile = &settings.Bullets{
			ShotImage: settings.ShotImage,
			CamAngle:  angle + 1.7,
			Count:     1,
			Speed:     8,
			X:         float64(x) - centerX,
			Y:         float64(y) - centerY,
		}
		player.Player.ChildProjectTile = append(player.Player.ChildProjectTile, projectTile)
	}

	return x, y

}

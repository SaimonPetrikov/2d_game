package player

import (
	"2d_game/helpers"
	"2d_game/project_tiles"
	"2d_game/settings"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

var screenWidth = settings.ScreenWidth
var screenHeight = settings.ScreenHeight

var Image *ebiten.Image

var FrameWidthCurrent = 20
var FrameHeightCurrent = 13

var Player = &project_tiles.ProjectTile{
	Width:      311,
	Height:     400,
	Count:      1,
	FrameCount: 17,
	Type:       1,
}

func DrawPlayer(screen *ebiten.Image) {
	var option = &ebiten.DrawImageOptions{}
	option.GeoM.Translate(-float64(Player.Width)/2, -float64(Player.Height)/2)
	option.GeoM.Scale(0.4, 0.4)

	Player.Angle = getAngle()

	option.GeoM.Rotate(1.7)
	option.GeoM.Rotate(Player.Angle)

	option.GeoM.Translate(float64(screenWidth/2), float64(screenHeight/2))

	i := (Player.Count / 5) % Player.FrameCount
	sx, sy := i*Player.Width, 0

	screen.DrawImage(
		Image.SubImage(image.Rect(sx, sy, sx+Player.Width, sy+Player.Height)).(*ebiten.Image),
		option)
}

func getAngle() float64 {
	mouseX, mouseY := ebiten.CursorPosition()
	centerX := float64(screenWidth / 2)
	centerY := float64(screenHeight / 2)
	var line = helpers.Line{X1: centerX, Y1: centerY, X2: float64(mouseX), Y2: float64(mouseY)}
	return line.Angle()
}

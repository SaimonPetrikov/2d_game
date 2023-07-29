package _map

import (
	"2d_game/settings"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileWidth  = 256
	tileHeight = 256
)

type Matrix struct {
	column int
	raw    int
}

var BackgroundImage *ebiten.Image

func CreateMap(screen *ebiten.Image, x float64, y float64) {
	var matrix = Matrix{
		raw:    settings.ScreenHeight/tileHeight + 1,
		column: settings.ScreenWidth/tileWidth + 1,
	}

	for i := -1; i < matrix.column; i++ {
		for j := -1; j < matrix.raw; j++ {
			bgOpt := &ebiten.DrawImageOptions{}
			bgOpt.GeoM.Translate(tileHeight*float64(i), tileWidth*float64(j))
			bgOpt.GeoM.Translate(x, y)

			screen.DrawImage(BackgroundImage, bgOpt)
		}
	}
}

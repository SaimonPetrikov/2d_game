package settings

import (
	"2d_game/helpers"
	"github.com/hajimehoshi/ebiten/v2"
)

var ShotImage *ebiten.Image

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
	Padding      = 20
)

type Object struct {
	Walls []helpers.Line
}

type Bullets struct {
	ShotImage *ebiten.Image
	CamAngle  float64
	Count     int
	Speed     float64
	X         float64
	Y         float64
}

type GameInterface interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int)
}

type GameObject struct {
	PX, PY       int
	OuterWalls   []Object
	Objects      []Object
	ProjectTiles []*Bullets
	Count        int
	CountShot    int
	ShotMove     int
}

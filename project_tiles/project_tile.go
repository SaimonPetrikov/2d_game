package project_tiles

import (
	"2d_game/settings"
	"github.com/hajimehoshi/ebiten/v2"
)

type TypeProjectTile int

type ChildProjectTile struct {
	X, Y       int
	Image      *ebiten.Image
	Width      int
	Height     int
	Angle      float64
	Count      int
	FrameCount int
	Speed      int
	Type       TypeProjectTile
}

type ProjectTile struct {
	X, Y             int
	Image            *ebiten.Image
	Width            int
	Height           int
	Angle            float64
	Count            int
	FrameCount       int
	Speed            int
	Type             TypeProjectTile
	ChildProjectTile []*settings.Bullets
}

const (
	Enemy  TypeProjectTile = 3
	Bullet                 = 2
	Player                 = 1
	Object                 = 4
)

func (projectTile ProjectTile) getType() string {
	switch projectTile.Type {
	case Enemy:
		return "Enemy"
	case Bullet:
		return "Bullet"
	case Player:
		return "Player"
	case Object:
		return "Object"
	}
	return "Unknown"
}

func (projectTile ProjectTile) getIdByType(typeTile string) int {
	switch typeTile {
	case "Enemy":
		return 3
	case "Bullet":
		return 2
	case "Player":
		return 1
	case "Object":
		return 4
	}
	return 0
}

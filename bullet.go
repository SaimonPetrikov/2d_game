package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"os"
)

type Bullet struct {
	shotWidth  int
	shotHeight int
	shotCount  int
}

// NewLevel returns a new randomly generated Level.
func (g *Game) NewBullet() *ebiten.Image {
	// Create a 108x108 Level.

	sfile, err := os.Open("images/shot.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(sfile *os.File) {
		err := sfile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(sfile)

	simg, _, err := image.Decode(sfile)
	if err != nil {
		log.Fatal(err)
	}

	return ebiten.NewImageFromImage(simg)
}

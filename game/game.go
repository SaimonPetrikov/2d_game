package game

import (
	"2d_game/helpers"
	_map "2d_game/map"
	"2d_game/movement"
	"2d_game/player"
	"2d_game/settings"
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"image/color"
	"log"
)

const (
	screenWidth        = settings.ScreenWidth
	screenHeight       = settings.ScreenHeight
	padding            = settings.Padding
	frameHeightCurrent = 13
)

type Game settings.GameObject

func init() {
	player.Image = helpers.NewImage("images/soldier_walking.png")
	settings.ShotImage = helpers.NewImage("images/shot.png")
	_map.BackgroundImage = helpers.NewImage("images/tile.png")
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("game ended by player")
	}

	g.PX, g.PY = movement.PlayerMovement(g.PX, g.PY)

	return nil
}

func (g *Game) renderSprite(screen *ebiten.Image, count int, x float64, y float64, camAngle float64) {
	opShot := &ebiten.DrawImageOptions{}
	opShot.GeoM.Translate(-float64(32)/2+7, -float64(32)/2-frameHeightCurrent*2-8)
	opShot.GeoM.Rotate(camAngle)
	opShot.GeoM.Translate(float64(g.PX)-x, float64(g.PY)-y)

	iShot := (count / 5) % 4
	shotx, shoty := 0+iShot*32, 0

	screen.DrawImage(settings.ShotImage.SubImage(image.Rect(shotx, shoty, shotx+32, shoty+32)).(*ebiten.Image), opShot)
}

func (g *Game) drawProjectiles(screen *ebiten.Image) {
	for _, p := range player.Player.ChildProjectTile {
		g.renderSprite(screen, p.Count, p.X, p.Y, p.CamAngle)
	}
}

func getRealPosition(gameX int, gameY int, absoluteX int, absoluteY int) (float64, float64) {
	return float64(gameX) - screenWidth/2 + float64(absoluteX),
		float64(gameY) - screenHeight/2 + float64(absoluteY)
}

func (g *Game) drawHouse(screen *ebiten.Image) {
	var houseImg = helpers.NewImage("images/house/roof_0036_Layer-0.png")

	houseImgOpt := &ebiten.DrawImageOptions{}

	houseImgOpt.GeoM.Translate(float64(g.PX)-screenWidth/2+1400, float64(g.PY)-screenHeight/2+20)

	if g.PX < 200 && g.PX > 144 && g.PY < 207 && g.PY > 107 {
		houseImg.Clear()
	} else {
		screen.DrawImage(houseImg, houseImgOpt)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	_map.CreateMap(screen, float64(g.PX)-screenWidth/2, float64(g.PY)-screenHeight/2)

	player.DrawPlayer(screen)

	g.drawProjectiles(screen)

	g.drawHouse(screen)

	// Draw outer walls
	for _, obj := range g.OuterWalls {
		for _, w := range obj.Walls {
			ebitenutil.DrawLine(screen, w.X1, w.Y1, w.X2, w.Y2, color.RGBA{255, 0, 0, 255})
		}
	}
	// Draw walls
	for _, obj := range g.Objects {
		for _, w := range obj.Walls {
			ebitenutil.DrawLine(screen, float64(g.PX)-w.X1, float64(g.PY)-w.Y1, float64(g.PX)-w.X2, float64(g.PY)-w.Y2, color.RGBA{255, 0, 0, 255})
		}
	}

	ebitenutil.DebugPrintAt(screen, "WASD: move", 160, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %d", g.PX), padding, 222)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Y: %d", g.PY), padding+60, 222)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func InitGame() {
	var g = &Game{
		PX: screenWidth / 2,
		PY: screenHeight / 2,
	}

	g.OuterWalls = append(g.OuterWalls, settings.Object{Walls: helpers.Rect(padding, padding, screenWidth-2*padding, screenHeight-2*padding)})

	//Rectangles
	g.Objects = append(g.Objects, settings.Object{Walls: helpers.Rect(32, -7, 68, 107)})

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Ray casting and shadows (Ebiten demo)")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

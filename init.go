// Copyright 2019 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"
	"sort"
)

const (
	screenWidth  = 240
	screenHeight = 240
	padding      = 20

	frameOX            = 0
	frameOY            = 0
	frameWidth         = 311
	frameHeight        = 400
	frameWidthCurrent  = 20
	frameHeightCurrent = 13
	frameCount         = 17
)

type Bullets struct {
	shotImage *ebiten.Image
	camAngle  float64
	move      int
	count     int
	speed     float64
	x         float64
	y         float64
}

var (
	bgImage       *ebiten.Image
	shadowImage   = ebiten.NewImage(screenWidth, screenHeight)
	triangleImage = ebiten.NewImage(screenWidth, screenHeight)
	runnerImage   *ebiten.Image
	shotImage     *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	// Now the byte slice is generated with //go:generate for Go 1.15 or older.
	// If you use Go 1.16 or newer, it is strongly recommended to use //go:embed to embed the image file.
	// See https://pkg.go.dev/embed for more details.
	img, _, err := image.Decode(bytes.NewReader(images.Tile_png))
	if err != nil {
		log.Fatal(err)
	}
	bgImage = ebiten.NewImageFromImage(img)
	triangleImage.Fill(color.White)
}

type line struct {
	X1, Y1, X2, Y2 float64
}

func (l *line) angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

type object struct {
	walls []line
}

func (o object) points() [][2]float64 {
	// Get one of the endpoints for all segments,
	// + the startpoint of the first one, for non-closed paths
	var points [][2]float64
	for _, wall := range o.walls {
		points = append(points, [2]float64{wall.X2, wall.Y2})
	}
	p := [2]float64{o.walls[0].X1, o.walls[0].Y1}
	if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
		points = append(points, [2]float64{o.walls[0].X1, o.walls[0].Y1})
	}
	return points
}

func newRay(x, y, length, angle float64) line {
	return line{
		X1: x,
		Y1: y,
		X2: x + length*math.Cos(angle),
		Y2: y + length*math.Sin(angle),
	}
}

// intersection calculates the intersection of given two lines.
func intersection(l1, l2 line) (float64, float64, bool) {
	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	denom := (l1.X1-l1.X2)*(l2.Y1-l2.Y2) - (l1.Y1-l1.Y2)*(l2.X1-l2.X2)
	tNum := (l1.X1-l2.X1)*(l2.Y1-l2.Y2) - (l1.Y1-l2.Y1)*(l2.X1-l2.X2)
	uNum := -((l1.X1-l1.X2)*(l1.Y1-l2.Y1) - (l1.Y1-l1.Y2)*(l1.X1-l2.X1))

	if denom == 0 {
		return 0, 0, false
	}

	t := tNum / denom
	if t > 1 || t < 0 {
		return 0, 0, false
	}

	u := uNum / denom
	if u > 1 || u < 0 {
		return 0, 0, false
	}

	x := l1.X1 + t*(l1.X2-l1.X1)
	y := l1.Y1 + t*(l1.Y2-l1.Y1)
	return x, y, true
}

// rayCasting returns a slice of line originating from point cx, cy and intersecting with objects
func rayCasting(cx, cy float64, objects []object) []line {
	const rayLength = 10000 // something large enough to reach all objects

	var rays []line
	for _, obj := range objects {
		// Cast two rays per point
		for _, p := range obj.points() {
			l := line{cx, cy, p[0], p[1]}
			angle := l.angle()

			for _, offset := range []float64{-0.005, 0.005} {
				points := [][2]float64{}
				ray := newRay(cx, cy, rayLength, angle+offset)

				// Unpack all objects
				for _, o := range objects {
					for _, wall := range o.walls {
						if px, py, ok := intersection(ray, wall); ok {
							points = append(points, [2]float64{px, py})
						}
					}
				}

				// Find the point closest to start of ray
				min := math.Inf(1)
				minI := -1
				for i, p := range points {
					d2 := (cx-p[0])*(cx-p[0]) + (cy-p[1])*(cy-p[1])
					if d2 < min {
						min = d2
						minI = i
					}
				}
				rays = append(rays, line{cx, cy, points[minI][0], points[minI][1]})
			}
		}
	}

	// Sort rays based on angle, otherwise light triangles will not come out right
	sort.Slice(rays, func(i int, j int) bool {
		return rays[i].angle() < rays[j].angle()
	})
	return rays
}

func (g *Game) handleMovement() {
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.px -= 2
		g.count++
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.py -= 2
		g.count++
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.px += 2
		g.count++
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.py += 2
		g.count++
	}

	// +1/-1 is to stop player before it reaches the border
	if g.px+frameWidthCurrent >= screenWidth-padding {
		g.px = screenWidth - padding - frameWidthCurrent
	}

	if g.px-frameWidthCurrent <= padding {
		g.px = padding + frameWidthCurrent
	}

	if g.py+frameHeightCurrent >= screenHeight-padding {
		g.py = screenHeight - padding - frameHeightCurrent
	}

	if g.py-frameHeightCurrent <= padding {
		g.py = padding + frameHeightCurrent
	}

	for _, p := range g.projectTiles {
		p.move++
		p.count++
		p.x = p.x - math.Sin(p.camAngle)*p.speed
		p.y = p.y + math.Cos(p.camAngle)*p.speed
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		centerX := float64(screenWidth / 2)
		centerY := float64(screenWidth / 2)
		angle := angleBetweenPoints(centerX, centerY, float64(mouseX), float64(mouseY))
		bullet := &Bullets{
			shotImage: shotImage,
			camAngle:  angle,
			move:      1,
			count:     1,
			speed:     5,
			x:         centerX,
			y:         centerY,
		}
		g.projectTiles = append(g.projectTiles, bullet)
	}

}

func rayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x3), DstY: float32(y3), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

type Game struct {
	showRays     bool
	px, py       int
	outerWalls   []object
	objects      []object
	projectTiles []*Bullets
	count        int
	countShot    int
	shotMove     int
}

func (g *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("game ended by player")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.showRays = !g.showRays
	}

	g.handleMovement()

	return nil
}

func angleBetweenPoints(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(x1-x2, y2-y1)
}

func (g *Game) DrawPlayer(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Scale(0.15, 0.15)
	//op.GeoM.Translate((screenWidth)/2, (screenHeight)/2)
	mouseX, mouseY := ebiten.CursorPosition()
	centerX := float64(screenWidth / 2)
	centerY := float64(screenHeight / 2)

	//// Рассчет угла поворота камеры
	camAngle := angleBetweenPoints(centerX, centerY, float64(mouseX), float64(mouseY))
	op.GeoM.Rotate(-3.3)
	op.GeoM.Rotate(camAngle)
	op.GeoM.Translate((screenWidth)/2, (screenHeight)/2)
	//op.GeoM.Translate(float64(g.px)-2, flfloat64oat64(g.py)-2)
	i := (g.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) renderSprite(screen *ebiten.Image, count int, x float64, y float64, camAngle float64) {
	opShot := &ebiten.DrawImageOptions{}
	opShot.GeoM.Translate(-float64(32)/2+7, -float64(32)/2-frameHeightCurrent*2-8)
	opShot.GeoM.Rotate(-3.3)
	opShot.GeoM.Rotate(camAngle)
	opShot.GeoM.Translate(x, y)
	//opShot.GeoM.Translate(-px, -py)
	//opShot.GeoM.Translate((screenWidth)/2, (screenHeight)/2)

	iShot := (count / 5) % 4
	shotx, shoty := 0+iShot*32, 0

	screen.DrawImage(shotImage.SubImage(image.Rect(shotx, shoty, shotx+32, shoty+32)).(*ebiten.Image), opShot)
}

func (g *Game) drawProjectiles(screen *ebiten.Image) {
	for _, p := range g.projectTiles {
		g.renderSprite(screen, p.count, p.x, p.y, p.camAngle)
	}
}

func (g *Game) drawHouse(screen *ebiten.Image) {

	sfile, err := os.Open("images/house/roof_0036_Layer-0.png")
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

	houseImg := ebiten.NewImageFromImage(simg)

	houseImgOpt := &ebiten.DrawImageOptions{}

	houseImgOpt.GeoM.Scale(0.4, 0.4)

	//xPosition, yPosition := float64(g.px)-screenWidth/2+20, float64(g.py)-screenHeight/2+20

	houseImgOpt.GeoM.Translate(float64(g.px)-screenWidth/2+20, float64(g.py)-screenHeight/2+20)

	if g.px < 200 && g.px > 144 && g.py < 207 && g.py > 107 {
		houseImg.Clear()
	} else {
		screen.DrawImage(houseImg, houseImgOpt)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Reset the shadowImage

	shadowImage.Fill(color.Black)
	rays := rayCasting(float64(g.px), float64(g.py), g.outerWalls)

	// Subtract ray triangles from shadow
	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.CompositeMode = ebiten.CompositeModeSourceOut
	for i, line := range rays {
		nextLine := rays[(i+1)%len(rays)]

		// Draw triangle of area between rays
		v := rayVertices(float64(g.px), float64(g.py), nextLine.X2, nextLine.Y2, line.X2, line.Y2)
		shadowImage.DrawTriangles(v, []uint16{0, 1, 2}, triangleImage, opt)
	}
	// Draw background
	bgOpt := &ebiten.DrawImageOptions{}

	bgOpt.GeoM.Translate(float64(g.px)-screenWidth/2, float64(g.py)-screenHeight/2)

	//bgOpt.GeoM.Rotate(camAngle)
	screen.DrawImage(bgImage, bgOpt)

	if g.showRays {
		// Draw rays
		for _, r := range rays {
			ebitenutil.DrawLine(screen, r.X1, r.Y1, r.X2, r.Y2, color.RGBA{255, 255, 0, 150})
		}
	}

	// Draw shadow
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.7)
	screen.DrawImage(shadowImage, op)

	// Draw outer walls
	for _, obj := range g.outerWalls {
		for _, w := range obj.walls {
			ebitenutil.DrawLine(screen, w.X1, w.Y1, w.X2, w.Y2, color.RGBA{255, 0, 0, 255})
		}
	}
	// Draw walls
	for _, obj := range g.objects {
		for _, w := range obj.walls {
			ebitenutil.DrawLine(screen, float64(g.px)-w.X1, float64(g.py)-w.Y1, float64(g.px)-w.X2, float64(g.py)-w.Y2, color.RGBA{255, 0, 0, 255})
		}
	}

	// Draw player as a rect
	g.DrawPlayer(screen)

	g.drawProjectiles(screen)

	g.drawHouse(screen)

	if g.showRays {
		ebitenutil.DebugPrintAt(screen, "R: hide rays", padding, 0)
	} else {
		ebitenutil.DebugPrintAt(screen, "R: show rays", padding, 0)
	}

	ebitenutil.DebugPrintAt(screen, "WASD: move", 160, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
	//ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Rays: 2*%d", len(rays)/2), padding, 222)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %d", g.px), padding, 222)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Y: %d", g.py), padding+20, 222)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func rect(x, y, w, h float64) []line {
	return []line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}

func main() {
	g := &Game{
		px: screenWidth / 2,
		py: screenHeight / 2,
	}
	file, err := os.Open("images/soldier_walking.png") // замените на свой файл
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	runnerImage = ebiten.NewImageFromImage(img)

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

	shotImage = ebiten.NewImageFromImage(simg)

	// Add outer walls
	g.outerWalls = append(g.outerWalls, object{rect(padding, padding, screenWidth-2*padding, screenHeight-2*padding)})

	//Rectangles
	g.objects = append(g.objects, object{rect(32, -7, 68, 107)})
	//g.objects = append(g.objects, object{rect(150, 50, 30, 60)})
	//
	//g.objects = append(g.objects, object{rect(50, 100, 50, -1)})
	//g.objects = append(g.objects, object{rect(50, 200, 50, -1)})
	//
	//g.objects = append(g.objects, object{rect(50, 100, -1, 40)})
	//g.objects = append(g.objects, object{rect(50, 160, -1, 40)})
	//
	//g.objects = append(g.objects, object{rect(100, 100, -1, 100)})

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Ray casting and shadows (Ebiten demo)")
	//ebiten.Cur()
	//ebiten.CursorModeHidden()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

package helpers

import "math"

type Line struct {
	X1, Y1, X2, Y2 float64
}

func (l *Line) Angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

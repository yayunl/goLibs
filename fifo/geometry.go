package fifo

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

func (p *Point) Distance(q Point) float64 {
	return math.Hypot(p.X-q.X, p.Y-q.Y)
}

func (p *Point) String() string {
	return fmt.Sprintf("The point is (%g, %g) \n", p.X, p.Y)
}

type Path []Point

func (pt *Path) Distance() float64 {
	sum := 0.0
	for i := len(*pt) - 1; i >= 0; i-- {
		if i > 0 {
			sum += (*pt)[i].Distance((*pt)[i-1])
		}

	}
	return sum
}

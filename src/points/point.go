package points

import (
	"fmt"
	"math"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (p Point) String() string {
	return fmt.Sprintf("%f,%f", p.Lat, p.Lon)
}

func (p Point) Distance(a Point) float64 {
	return math.Sqrt(math.Pow(p.Lon-a.Lon, 2) + math.Pow(p.Lat-a.Lat, 2))
}

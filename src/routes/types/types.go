package routes

import (
	"fmt"
	"math"
	"objects"
	"points"
)

type RouteBuilder interface {
	Build(objects []objects.Object) []objects.Object
}

type ABRoute struct {
	start  points.Point
	finish points.Point
}

const coefficientCurrentLengthToFinish = 1 / 1.6

func (r ABRoute) selectBest(current points.Point, finish points.Point, allObjects []objects.Object) (*objects.Object, int) {
	var selected *objects.Object = nil
	var selectedIndex int
	minDistance := 999.0
	for index, next := range allObjects {
		distanceNextToFinish := getManhattanDistance(next.Position, finish)
		distanceCurrentToFinish := getManhattanDistance(current, finish)
		distance := getManhattanDistance(current, next.Position) + distanceNextToFinish*coefficientCurrentLengthToFinish
		if distance < minDistance && distanceNextToFinish < distanceCurrentToFinish {
			minDistance = distance
			selected = &allObjects[index]
			selectedIndex = index
			fmt.Print("2")
		}
	}
	return selected, selectedIndex
}

func (r ABRoute) Build(allObjects []objects.Object) []objects.Object {

	allFind := false

	var pathObjects []objects.Object

	current := r.start

	for !allFind {
		selected, selectedIndex := r.selectBest(current, r.finish, allObjects)

		if selected != nil {
			pathObjects = append(pathObjects, *selected)
			current = selected.Position
			allObjects = append(allObjects[:selectedIndex], allObjects[selectedIndex+1:]...)
		} else {
			allFind = true
		}
	}

	return pathObjects
}

type RoundRoute struct {
	start  points.Point
	radius int
}

func (r RoundRoute) Build([]objects.Object) []objects.Object {
	// todo: implement
	return nil
}

func getManhattanDistance(a points.Point, b points.Point) float64 {
	return getPythagorasDistance(a.Lat, a.Lon, b.Lat, a.Lon) + getPythagorasDistance(a.Lat, a.Lon, a.Lat, b.Lon)
}

func getPythagorasDistance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	dLat := (lat1 + lat2) / 2
	lon1P := lon1 * math.Cos(deg2rad(dLat)) // поправка из-за разницы широты и долготы
	lon2P := lon2 * math.Cos(deg2rad(dLat))
	return math.Sqrt(math.Pow(lat1-lat2, 2) + math.Pow(lon1P-lon2P, 2))
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

package routes

import (
	"math"
	"objects"
	"points"
	"sort"
)

type RouteBuilder interface {
	Build(objects []objects.Object) []objects.Object
}

type ABRoute struct {
	Start  points.Point
	Finish points.Point
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
		}
	}
	return selected, selectedIndex
}

func (r ABRoute) Build(allObjects []objects.Object) []objects.Object {

	allFind := false

	var pathObjects []objects.Object

	current := r.Start

	for !allFind {
		selected, selectedIndex := r.selectBest(current, r.Finish, allObjects)

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
	Center points.Point
	Radius int // in meters
}

func (r RoundRoute) Build(allObjects []objects.Object) []objects.Object {
	if len(allObjects) == 0 {
		return allObjects
	}

	var pathObjects []objects.Object

	var firstObject = allObjects[0]
	var minLat = firstObject.Position.Lat
	var minLon = firstObject.Position.Lon

	for _, one := range allObjects {
		if GetMetersDistanceByPoints(one.Position, r.Center) < GetMetersDistanceByPoints(firstObject.Position, r.Center) {
			firstObject = one
		}
		if one.Position.Lat < minLat {
			minLat = one.Position.Lat
		}
		if one.Position.Lon < minLon {
			minLon = one.Position.Lon
		}
	}

	var startCoordinates = points.Point{Lat: minLat, Lon: minLon}

	pathObjects = append(pathObjects, firstObject)

	var centerX = getMetersDistance(startCoordinates.Lat, startCoordinates.Lon, startCoordinates.Lat, r.Center.Lon)
	var centerY = getMetersDistance(startCoordinates.Lat, startCoordinates.Lon, r.Center.Lat, startCoordinates.Lon)
	var centerCoordinate = Coordinate{centerX, centerY}

	var firstPointX = getMetersDistance(startCoordinates.Lat, startCoordinates.Lon, startCoordinates.Lat, firstObject.Position.Lon)
	var firstPointY = getMetersDistance(startCoordinates.Lat, startCoordinates.Lon, firstObject.Position.Lat, startCoordinates.Lon)
	var firstPointCoordinate = Coordinate{firstPointX, firstPointY}

	var objectData []Data
	for _, one := range allObjects {
		if one == firstObject {
			continue
		}

		var x = getMetersDistance(startCoordinates.Lat, startCoordinates.Lon, startCoordinates.Lat, one.Position.Lon)
		var y = getMetersDistance(startCoordinates.Lat, startCoordinates.Lon, one.Position.Lat, startCoordinates.Lon)

		var angle = getAngle(centerCoordinate, firstPointCoordinate, Coordinate{x, y})
		objectData = append(objectData, Data{angle, Coordinate{x, y}, one})
	}

	sort.Slice(objectData, func(i, j int) bool {
		return objectData[i].angle < objectData[j].angle
	})

	var currentDistanceStartToPoint = getPythagorasDistance(centerCoordinate.x, centerCoordinate.y, firstPointCoordinate.x, firstPointCoordinate.y)
	var distanceNorm = r.getDistanceNorm(currentDistanceStartToPoint)
	var currentAngle = 0.0

	for _, candidate := range objectData {

		if candidate.angle < currentAngle+AgeStep {
			// пропускаем точку, если угол между ней, центром и предыдущей точкой меньше 5 градусов
			// это позволяет делать маршрут равномерным
			continue
		}

		var candidateDistance = getPythagorasDistance(centerCoordinate.x, centerCoordinate.y, candidate.coordinate.x, candidate.coordinate.y)
		if candidateDistance/distanceNorm > 2 || candidateDistance/distanceNorm < 0.5 {
			// расстояние от центра до следующей точки не должно отличаться более чем в два раза от расстояния до предыдущей точки
			continue
		}

		pathObjects = append(pathObjects, candidate.objectEntity)
		currentAngle = candidate.angle
		distanceNorm = r.getDistanceNorm(candidateDistance)
	}

	return pathObjects
}

func (r RoundRoute) getDistanceNorm(distance float64) float64 {
	return math.Max(distance, float64(r.Radius)/3)
}

type Data struct {
	angle        float64
	coordinate   Coordinate
	objectEntity objects.Object
}

type Coordinate struct {
	x float64
	y float64
}

const EarthRadius = 6371000
const EquatorLength = 40075000
const AgeStep = 5

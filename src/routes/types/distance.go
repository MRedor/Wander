package routes

import (
	"math"
	"points"
)

func getAngle(c Coordinate, a Coordinate, b Coordinate) float64 {
	var aFixed = Coordinate{a.x - c.x, a.y - c.y}
	var bFixed = Coordinate{b.x - c.x, b.y - c.y}

	var angle1 = toDegrees(math.Atan2(aFixed.x, aFixed.y))
	var angle2 = toDegrees(math.Atan2(bFixed.x, bFixed.y))
	var result = angle1 - angle2
	if result < 0 {
		result += 360
	}
	return result
}

func getManhattanDistance(a points.Point, b points.Point) float64 {
	return getPythagorasDistance(a.Lat, a.Lon, b.Lat, a.Lon) + getPythagorasDistance(a.Lat, a.Lon, a.Lat, b.Lon)
}

func getPythagorasDistance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	dLat := (lat1 + lat2) / 2
	lon1P := lon1 * math.Cos(toDegrees(dLat)) // поправка из-за разницы широты и долготы
	lon2P := lon2 * math.Cos(toDegrees(dLat))
	return math.Sqrt(math.Pow(lat1-lat2, 2) + math.Pow(lon1P-lon2P, 2))
}

func GetMetersDistanceByPoints(a points.Point, b points.Point) float64 {
	return getMetersDistance(a.Lat, a.Lon, b.Lat, b.Lon)
}

func getMetersDistance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	var dLat = toRadians(lat2 - lat1)
	var dLng = toRadians(lon2 - lon1)
	var q = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	var c = 2 * math.Atan2(math.Sqrt(q), math.Sqrt(1-q))
	return EarthRadius * c
}

func toRadians(angdeg float64) float64 {
	return angdeg / 180.0 * math.Pi
}
func toDegrees(angrad float64) float64 {
	return angrad * 180.0 / math.Pi
}

func MetersToLat(distance float64) float64 {
	return distance / (float64(EquatorLength) / 360.0)
}

func MetersToLon(a points.Point, distance float64) float64 {
	return distance / (EquatorLength / 360 * math.Cos(toRadians(a.Lat)))
}

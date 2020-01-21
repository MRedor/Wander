package routes

import (
	"math"
	"points"
	"testing"
)

func Test_getAngle(t *testing.T) {
	errorValue := 5.0

	type args struct {
		c Coordinate
		a Coordinate
		b Coordinate
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "90 degrees",
			args: args{a: Coordinate{55, 27}, b: Coordinate{26, 31}, c: Coordinate{26, 27}},
			want: 90.0,
		},
		{
			name: "0 degrees",
			args: args{a: Coordinate{0, 1}, b: Coordinate{0, 1}, c: Coordinate{0, 0}},
			want: 0.0,
		},
		{
			name: "180 degrees",
			args: args{a: Coordinate{160, -1}, b: Coordinate{160, 1}, c: Coordinate{160, 0}},
			want: 180.0,
		},
		{
			name: "30 degrees",
			args: args{a: Coordinate{2, 0}, b: Coordinate{2, 1}, c: Coordinate{0, 0}},
			want: 30.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAngle(tt.args.c, tt.args.a, tt.args.b)
			if math.Abs(tt.want-got) > errorValue {
				t.Errorf("getMetersDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMetersDistance(t *testing.T) {
	errorPercents := 0.05

	type args struct {
		lat1 float64
		lon1 float64
		lat2 float64
		lon2 float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "zero value",
			args: args{lat1: 30.52131, lon1: 59.123, lat2: 30.52131, lon2: 59.123},
			want: 0,
		},
		{
			name: "test value",
			args: args{lat1: 59.939069, lon1: 30.315804, lat2: 59.937497, lon2: 30.308669},
			want: 434,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getMetersDistance(tt.args.lat1, tt.args.lon1, tt.args.lat2, tt.args.lon2)
			if math.Abs(tt.want-got)/tt.want > errorPercents {
				t.Errorf("getMetersDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetersToLat(t *testing.T) {
	errorPercents := 0.05

	distance := 1000.0
	latDiff := MetersToLat(distance)
	point1 := points.Point{Lat: 60, Lon: 60}
	point2 := points.Point{Lat: point1.Lat - latDiff, Lon: point1.Lon}
	distanceCheck := GetMetersDistanceByPoints(point1, point2)

	if math.Abs(distance-distanceCheck)/distanceCheck > errorPercents {
		t.Errorf("MetersToLat() = %v, want %v", distanceCheck, distance)
	}
}

func TestMetersToLon(t *testing.T) {
	errorPercents := 0.05

	distance := 1000.0
	point1 := points.Point{Lat: 30, Lon: 10}
	lonDiff := MetersToLon(point1, distance)
	point2 := points.Point{Lat: point1.Lat, Lon: point1.Lon - lonDiff}
	distanceCheck := GetMetersDistanceByPoints(point1, point2)

	if math.Abs(distance-distanceCheck)/distanceCheck > errorPercents {
		t.Errorf("MetersToLat() = %v, want %v", distanceCheck, distance)
	}
}

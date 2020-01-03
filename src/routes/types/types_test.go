package routes

import (
	"objects"
	"points"
	"reflect"
	"testing"
)

func TestABRoute_Build(t *testing.T) {
	type fields struct {
		start  points.Point
		finish points.Point
	}
	type args struct {
		allObjects []objects.Object
	}

	// todo: make real tests
	allObjects := []objects.Object{
		{Id: 1, Position: points.Point{Lat: 1.5, Lon: 1.4}},
		{Id: 2, Position: points.Point{Lat: 1.8, Lon: 1.6}},
		{Id: 3, Position: points.Point{Lat: 1.2, Lon: 1.1}},
		{Id: 4, Position: points.Point{Lat: 1.3, Lon: 1.4}},
		{Id: 5, Position: points.Point{Lat: 1.6, Lon: 1.2}},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []objects.Object
	}{
		{
			name:   "test1",
			fields: fields{start: points.Point{Lat: 1.15, Lon: 1.18}, finish: points.Point{Lat: 1.9, Lon: 1.8}},
			args:   args{allObjects: allObjects},
			want:   []objects.Object{allObjects[3], allObjects[0], allObjects[1]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ABRoute{
				Start:  tt.fields.start,
				Finish: tt.fields.finish,
			}
			if got := r.Build(tt.args.allObjects); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestRoundRoute_Build(t *testing.T) {
	type fields struct {
		center points.Point
		radius int
	}
	type args struct {
		allObjects []objects.Object
	}

	allObjects := []objects.Object{
		{Id: 1, Position: points.Point{Lat: 1, Lon: 1}},
		{Id: 2, Position: points.Point{Lat: 1, Lon: 2.5}},
		{Id: 3, Position: points.Point{Lat: 2, Lon: 3}},
		{Id: 4, Position: points.Point{Lat: 3, Lon: 2.5}},
		{Id: 5, Position: points.Point{Lat: 2.5, Lon: 1}},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []objects.Object
	}{
		{
			name:   "test2",
			fields: fields{center: points.Point{Lat: 2, Lon: 2}, radius: 100},
			args:   args{allObjects: allObjects},
			want:   []objects.Object{allObjects[2], allObjects[3], allObjects[4], allObjects[0], allObjects[1]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := RoundRoute{
				Center: tt.fields.center,
				Radius: tt.fields.radius,
			}
			if got := r.Build(tt.args.allObjects); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

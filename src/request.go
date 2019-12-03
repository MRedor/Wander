package main

type RouteRequest struct {
	Points []Point `json:"points"`
	Radius int     `json:"radius"`
	Type   string  `json:"type"`
}

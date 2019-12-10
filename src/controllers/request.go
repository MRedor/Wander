package controllers

import (
	"points"
)

type RouteRequest struct {
	Points []points.Point `json:"points"`
	Radius int            `json:"radius"`
	Type   string         `json:"type"`
}

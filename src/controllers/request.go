package controllers

import (
	"data"
)

type RouteRequest struct {
	Points []data.Point `json:"points"`
	Radius int          `json:"radius"`
	Type   string       `json:"type"`
}

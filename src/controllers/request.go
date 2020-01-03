package controllers

import (
	"points"
)

type RouteRequest struct {
	Points  []points.Point `json:"points"`
	Radius  int            `json:"radius"`
	Type    string         `json:"type"`
	Filters []string       `json:"filters"`
}

type RemovePointRequest struct {
	ObjectId int
	RouteId  int
}

type GetListsRequest struct {
	Count  int64 `json:"count"`
	Offset int64 `json:"offset"`
}

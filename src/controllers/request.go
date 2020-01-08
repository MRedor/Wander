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

type FeedbackRequest struct {
	Email     string `json:"email"`
	Text      string `json:"text"`
	RelatedTo struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
	} `json:"relatedTo"`
}

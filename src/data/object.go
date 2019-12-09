package data

import "fmt"

type Object struct {
	Id       int    `json:"id"`
	Name     string `json:"title"`
	Position Point  `json:"position"`
	Image    string `json:"image"`
	Type     string `json:"type"`

	Address string `json:"address"`
	Url     string `json:"link"`
	Prices  string `json:"price"`
	//workingTime
	Description string `json:"description"`
}

func (p Object) PositionString() string {
	return fmt.Sprintf("%f,%f", p.Position.Lat, p.Position.Lon)
}

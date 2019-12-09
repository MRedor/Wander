package data

type Route struct {
	Objects []Object `json:"objects"`
	Points  []Point  `json:"points"`
	Id      int      `json:"id"`
	Length  float64  `json:"length"` //meters
	Time    int      `json:"time"`   //seconds
	Name    string   `json:"name"`
	Type    string
	Radius  int
}

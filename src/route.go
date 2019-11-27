package main

type Route struct {
	Objects []Object `json:"objects"`
	Points  []Point  `json:"points"`
	Id      int      `json:"id"`
	Length  int      `json:"length"` //meters
	Time    int      `json:"time"`   //seconds
	Name    string   `json:"name"`
}

func ABRoute(a, b DBObject) Route {
	return Route{}
}

func CircularRoute(start DBObject) Route {
	return Route{}
}

func NameByRoute(path Route) string {
	// хотим генерить что-то умное и триггерное
	return ""
}

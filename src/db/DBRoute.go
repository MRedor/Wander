package db

type DBRoute struct {
	Id         int64     `db:"id"`
	Type       string  `db:"type"`
	Start_lat  float64 `db:"start_lat"`
	Start_lon  float64 `db:"start_lon"`
	Finish_lat float64 `db:"finish_lat"`
	Finish_lon float64 `db:"finish_lon"`
	Radius     int     `db:"radius"`
	Length     float64 `db:"length"`
	Time       int     `db:"time"`
	Objects    string  `db:"objects"`
	Points     string  `db:"points"`
	Name       string  `db:"name"`
	Count      int     `db:"count"`
}
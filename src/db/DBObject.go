package db

import (
	"data"
	"database/sql"
)

type DBObject struct {
	Id                 int            `db:"id"`
	Lat                float64        `db:"lat"`
	Lon                float64        `db:"lon"`
	Type               string         `db:"type"`
	Name               sql.NullString `db:"name"`
	Description        sql.NullString `db:"description"`
	Time               sql.NullString `db:"time"`
	NameEnglish        sql.NullString `db:"name_en"`
	DescriptionEnglish sql.NullString `db:"description_en"`
	Image              sql.NullString `db:"img"`
	IdInGraph          int            `db:"id_point_in_graph"`
	Updated            string         `db:"updated"`
	Active             bool           `db:"active"`
	Address            sql.NullString `db:"address"`
	Prices             sql.NullString `db:"prices"`
	Url                sql.NullString `db:"url"`
	NightDescription   sql.NullString `db:"night_description"`
	NightImage         sql.NullString `db:"night_photo"`
	NightType          sql.NullString `db:"night_type"`
	ActiveOnlyAtNight  sql.NullInt64  `db:"active_only_at_night"`
}

func (o *DBObject) Object() *data.Object {
	return &data.Object{
		Id:          o.Id,
		Name:        o.Name.String,
		Position:    data.Point{Lat: o.Lat, Lon: o.Lon},
		Image:       "https://travelpath.ru" + o.Image.String,
		Type:        o.Type,
		Address:     o.Address.String,
		Url:         o.Url.String,
		Prices:      o.Prices.String,
		Description: o.Description.String,
	}
}

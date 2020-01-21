package db

type DBList struct {
	Id int64 `db:"id"`
	Type string `db:"type"`
	Name string `db:"name"`
	Views string `db:"views"`
}

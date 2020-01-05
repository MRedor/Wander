package lists

import (
	"db"
	"objects"
	"routes"
)

type List struct {
	Id      int64            `json:"id"`
	Name    string           `json:"name"`
	Routes  []routes.Route   `json:"routes"`
	Objects []objects.Object `json:"objects"`
}

func ListByDBList(dblist db.DBList) List {
	return List{
		Id:   dblist.Id,
		Name: dblist.Name,
	}
}

func ListById(id int64) (*List, error) {
	dblist, err := db.DBListById(id)
	if err != nil {
		return nil, err
	}
	list := ListByDBList(*dblist)

	var dbobjects []db.DBObject
	var dbroutes []db.DBRoute
	if dblist.Type == "routes" {
		dbobjects, err = db.ObjectsForList(dblist.Id)
		for _, dbo := range dbobjects {
			list.Objects = append(list.Objects, *objects.ObjectByDBObject(&dbo))
		}
	} else {
		dbroutes, err = db.RoutesForList(dblist.Id)
		for _, dbr := range dbroutes {
			list.Routes = append(list.Routes, *routes.RouteByDBRoute(&dbr))
		}
	}
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func GetSliceOfLists(count, offset int64) ([]List, error) {
	dblists, err := db.GetLists(count, offset)
	if err != nil {
		return nil, err
	}
	var lists []List
	for _, dblist := range dblists {
		lists = append(lists, ListByDBList(dblist))
	}
	return lists, nil
}

func IsLastInSlice(count, offset int64) (bool, error) {
	lastId, err := db.LastListId()
	if err != nil {
		return false, nil
	}

	if int64(lastId) > offset + count {
		return false, nil
	}
	return true, nil
}

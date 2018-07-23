package datastore

type OrderType int

const (
	Asc OrderType = iota + 1
	Desc
)

func (o OrderType) String() string {
	switch o {
	case Asc:
		return "ASC"
	case Desc:
		return "DESC"
	default:
		return ""
	}
}

type Orders map[string]OrderType

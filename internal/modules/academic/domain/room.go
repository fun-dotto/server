package domain

type Floor string

const (
	FloorFloor1 Floor = "Floor1"
	FloorFloor2 Floor = "Floor2"
	FloorFloor3 Floor = "Floor3"
	FloorFloor4 Floor = "Floor4"
	FloorFloor5 Floor = "Floor5"
	FloorFloor6 Floor = "Floor6"
	FloorFloor7 Floor = "Floor7"
)

type Room struct {
	ID    string
	Name  string
	Floor Floor
}

type RoomListFilter struct {
	IDs    []string
	Floors []Floor
}

package utils

type DeleteSet struct {
	Clients map[int][]DeleteItem
}

type DeleteItem struct {
	Clock  int
	Length int
}

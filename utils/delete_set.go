package utils

type DeleteSet struct {
	Clients map[int][]*DeleteItem
}

type DeleteItem struct {
	Clock  int
	Length int
}

func (ds DeleteSet) Add(client, clock, length uint64) {
	t := make([]*DeleteItem, 2)
	if item, ok := ds.Clients[int(client)]; ok {
		t = item
	}
	t = append(t, &DeleteItem{
		Clock:  int(clock),
		Length: int(length),
	})
}

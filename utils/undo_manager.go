package utils

type StackItem struct {
	BeforeState map[int]int
	AfterState  map[int]int
	Meta        map[string]interface{}
	DeleteSet   DeleteSet
}

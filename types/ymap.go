package types

import "YJS-GO/utils"

type YMapEvent struct {
	utils.YEvent
	KeysChanged map[string]struct{}
}

func NewYMapEvent(ymap *YMap, t *utils.Transaction, subs map[string]struct{}) *YMapEvent {
	m := &YMapEvent{
		YEvent: utils.YEvent{
			NewBaseType: utils.NewBaseType{},
			Transaction: t,
		},
		KeysChanged: subs,
	}

	return m
}

type YMap struct {
	AbstractType
}

package types

const (
	YArrayRefId uint = iota
	YMapRefId
	YTextRefId
	YXmlElementRefID
	YXmlFragmentRefID
	YXmlHookRefID
	YXmlTextRefID
)

var TypeMap = map[uint]any{
	YArrayRefId:       ReadArr,
	YMapRefId:         ReadMap,
	YTextRefId:        ReadText,
	YXmlElementRefID:  ReadMap,
	YXmlFragmentRefID: ReadMap,
	YXmlHookRefID:     ReadMap,
	YXmlTextRefID:     ReadMap,
}

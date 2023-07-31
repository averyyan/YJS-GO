package types

import (
	"YJS-GO/structs"
	"YJS-GO/structs/content"
	"YJS-GO/utils"
)

const (
	Insert = iota
	Delete
	Retain
)

type YText struct {
	YArrayBase
}

type ItemTextListPosition struct {
	Left              *structs.Item
	Right             *structs.Item
	Index             uint64
	CurrentAttributes map[string]any
}

func UpdateCurrentAttributes(attributes map[string]any, format content.Format) {
	if format.Value == nil {
		delete(attributes, format.Key)
	} else {
		attributes[format.Key] = format.Value
	}
}

type YTextEvent struct {
	utils.YEvent
	subs             map[string]struct{}
	delta            []*utils.Delta
	ChildListChanged bool
	KeysChanged      map[string]struct{}
}

func NewYTextEvent(arr *YText, transaction *utils.Transaction, subs map[string]struct{},
	target *AbstractType) *YTextEvent {
	textEvent := &YTextEvent{
		YEvent: utils.YEvent{
			NewBaseType: utils.NewBaseType{
				Target:      target,
				Transaction: transaction,
			},
			Transaction: transaction,
		},
		subs: subs,
	}
	for s, _ := range subs {
		if s == "" {
			textEvent.ChildListChanged = true
		} else {
			textEvent.KeysChanged[s] = struct{}{}
		}
	}
	return textEvent
}

func (yte YTextEvent) GetDelta() []*utils.Delta {
	if yte.delta == nil {
		var doc = yte.Target.Doc
		yte.delta = []*utils.Delta{}

		doc.Transact(func(transaction *utils.Transaction) {
			var delta = yte.delta

			// Saves all current attributes for insert.
			var currentAttributes = map[string]any{}
			var oldAttributes = map[string]any{}
			var item = yte.Target.Start
			var action = 0
			var attributes = map[string]any{}

			var (
				insert    any    = ""
				retain    uint64 = 0
				deleteLen uint64 = 0
			)

			addOp := func() {
				if action != 0 {
					var op *utils.Delta

					switch action {
					case Delete:
						op = &utils.Delta{Delete: deleteLen}
						deleteLen = 0
						break
					case Insert:
						op = &utils.Delta{Insert: insert}
						if len(currentAttributes) > 0 {
							op.Attributes = map[string]any{}
							for k, v := range currentAttributes {
								op.Attributes[k] = v
							}
						}
						break
					case Retain:
						op = &utils.Delta{Retain: retain}
						if len(attributes) > 0 {
							op.Attributes = map[string]any{}
							for k, _ := range attributes {
								op.Attributes[k] = attributes[k]
							}
						}
						retain = 0
						break
					default:
						// throw new InvalidOperationException($"Unexpected action: {action}")
						return
					}

					yte.delta = append(yte.delta, op)
					action = 0
				}
			}

			for item != nil {
				switch c := item.Content.(type) {
				case content.Embed:
					if yte.Adds(item) {
						if !yte.Deletes(item) {
							addOp()
							action = Insert
							insert = item.Content.(content.Embed).Embed
							addOp()
						}
					} else if yte.Deletes(item) {
						if action != Delete {
							addOp()
							action = Delete
						}

						deleteLen++
					} else if !item.Deleted {
						if action != Retain {
							addOp()
							action = Retain
						}

						retain++
					}
					break
				case content.String:
					if yte.Adds(item) {
						if !yte.Deletes(item) {
							if action != Insert {
								addOp()
								action = Insert
							}
							insert = insert.(string) + item.Content.(content.String).GetString()
						}
					} else if yte.Deletes(item) {
						if action != Delete {
							addOp()
							action = Delete
						}

						deleteLen += item.Length
					} else if !item.Deleted {
						if action != Retain {
							addOp()
							action = Retain
						}
						retain += item.Length
					}
					break
				case content.Format:
					if yte.Adds(item) {
						if !yte.Deletes(item) {
							var (
								curVal any
								ok     bool
							)
							if curVal, ok = currentAttributes[c.Key]; !ok {
								curVal = nil
							}

							if !utils.EqualAttrs(curVal, c.Value) {
								if action == Retain {
									addOp()
								}
								var (
									oldVal any
									ok     bool
								)
								if oldVal, ok = oldAttributes[c.Key]; !ok {
									oldVal = nil
								}

								if utils.EqualAttrs(c.Value, oldVal) {
									delete(attributes, c.Key)
								} else {
									attributes[c.Key] = c.Value
								}
							} else {
								item.Delete(transaction)
							}
						}
					} else if yte.Deletes(item) {
						oldAttributes[c.Key] = c.Value
						var (
							curVal any
							ok     bool
						)
						if curVal, ok = currentAttributes[c.Key]; !ok {
							curVal = nil
						}
						if !utils.EqualAttrs(curVal, c.Value) {
							if action == Retain {
								addOp()
							}

							attributes[c.Key] = curVal
						}
					} else if !item.Deleted {
						oldAttributes[c.Key] = c.Value
						if attr, ok := attributes[c.Key]; ok {
							if !utils.EqualAttrs(attr, c.Value) {
								if action == Retain {
									addOp()
								}
								if c.Value == nil {
									attributes[c.Key] = nil
								} else {
									delete(attributes, c.Key)
								}
							} else {
								item.Delete(transaction)
							}
						}
					}
					if !item.Deleted {
						if action == Insert {
							addOp()
						}
						UpdateCurrentAttributes(currentAttributes, item.Content.(content.Format))
					}
					break
				}
				item = item.Right.(*structs.Item)
			}
			addOp()
			for len(delta) > 0 {
				var lastOp = delta[len(delta)-1]
				if lastOp.Retain != 0 && lastOp.Attributes != nil {
					// Retain delta's if they don't assign attributes.
					delta = delta[:len(delta)-2]
					// delta.RemoveAt(len(delta) - 1)
				} else {
					break
				}
			}
		}, nil)

	}
	return yte.delta
}

func (p ItemTextListPosition) Forward() {
	if p.Right == nil {
		// throw new Exception("Unexpected");
		return
	}

	switch cf := p.Right.Content.(type) {
	case content.Embed:
	case content.String:
		if !p.Right.Deleted {
			p.Index += p.Right.Length
		}
		break
	case content.Format:
		if !p.Right.Deleted {
			UpdateCurrentAttributes(p.CurrentAttributes, cf)
		}
		break
	}

	p.Left = p.Right
	p.Right = p.Right.Right.(*structs.Item)
}

func ReadText(decoder utils.IUpdateDecoder) *YText {
	return &YText{}
}

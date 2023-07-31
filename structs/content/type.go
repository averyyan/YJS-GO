package content

import (
	"errors"

	"YJS-GO/structs"
	"YJS-GO/types"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Type)(nil)

type Type struct {
	Type *types.AbstractType
}

func NewType(v *types.AbstractType) *Type {
	return &Type{v}
}

func (t Type) SetRef(i int) {
	// Do nothing.
}

func ReadType(decoder utils.IUpdateDecoder) (*Type, error) {
	var typeRef = decoder.ReadTypeRef()
	switch typeRef {
	case types.YArrayRefId:
		var arr = types.ReadArr(decoder)
		return NewType(arr), nil
	case types.YMapRefId:
		var m = types.ReadMap(decoder)
		return NewType(m), nil
	case types.YTextRefId:
		var text = types.ReadText(decoder)
		return NewType(text), nil
	default:
		// throw new NotImplementedException($"Type {typeRef} not implemented")
		return nil, errors.New("Type {typeRef} not implemented")
	}
}

func (t Type) Copy() structs.IContent {
	return NewType(t.Type)
}

func (t Type) Splice(offset uint64) structs.IContent {
	// Do nothing.
	return nil
}

func (t Type) MergeWith(right structs.IContent) bool {
	return false
}

func (t Type) GetContent() any {
	return []any{t.Type}
}

func (t Type) GetLength() int {
	return 1
}

func (t Type) Countable() bool {
	return true
}

func (t Type) Write(encoder utils.IUpdateEncoder, offset int) {
	t.Type.Write(encoder)
}

func (t Type) Gc(store *utils.StructStore) {
	var item = t.Type.Start
	for item != nil {
		item.Gc(store, true)
		item = item.Right.(*structs.Item)
	}

	t.Type.Start = nil
	for _, valueItem := range t.Type.ItemMap {
		for valueItem != nil {
			valueItem.Gc(store, true)
			valueItem = valueItem.Left.(*structs.Item)
		}
	}
	// Clear
	t.Type.ItemMap = map[string]*structs.Item{}
}

func (t Type) Delete(transaction *utils.Transaction) {
	var item = t.Type.Start

	for item != nil {
		if !item.Deleted {
			item.Delete(transaction)
		} else {
			// This will be gc'd later and we want to merge it if possible.
			// We try to merge all deleted items each transaction,
			// but we have no knowledge about that this needs to merged
			// since it is not in transaction. Hence we add it to transaction._mergeStructs.
			transaction.MergeStructs = append(transaction.MergeStructs, item)
		}

		item = item.Right.(*structs.Item)
	}
	for _, valueItem := range t.Type.ItemMap {
		if !valueItem.Deleted {
			valueItem.Delete(transaction)
		} else {
			// Same as above.
			transaction.MergeStructs = append(transaction.MergeStructs, item)
		}
	}
	delete(transaction.Changed, t.Type)
}

func (t Type) Integrate(transaction *utils.Transaction, item *structs.Item) {
	t.Type.Integrate(transaction.Doc, item)
}

func (t Type) GetRef() int {
	return TypeRef
}

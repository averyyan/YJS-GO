package content

import (
	"errors"

	"YJS-GO/structs"
	"YJS-GO/types"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Type)(nil)

type Type struct {
	// Type *types.AbstractType
	Type any
}

func NewType(v any) *Type {
	return &Type{v}
}

func (t Type) SetRef(i int) {
	// Do nothing.
}

func ReadType(decoder utils.IUpdateDecoder) (*Type, error) {
	var typeRef = decoder.ReadTypeRef()
	var tmp any
	switch typeRef {
	case types.YArrayRefId:
		tmp = types.ReadArr(decoder)
	case types.YMapRefId:
		tmp = types.ReadMap(decoder)
	case types.YTextRefId:
		tmp = types.ReadText(decoder)
	case types.YXmlElementRefID:
		var text = types.ReadXmlElement(decoder)
		return NewType(text), nil
	default:
		// throw new NotImplementedException($"Type {typeRef} not implemented")
		return nil, errors.New("Type {typeRef} not implemented")
	}
	return NewType(tmp), nil
}

func (t Type) Copy() structs.IContentExt {
	return NewType(t.GetType())
}

func (t Type) Splice(offset uint64) structs.IContentExt {
	// Do nothing.
	return nil
}

func (t Type) MergeWith(right structs.IContentExt) bool {
	return false
}

func (t Type) GetContent() any {
	return []any{t.GetType()}
}

func (t Type) GetLength() int {
	return 1
}

func (t Type) Countable() bool {
	return true
}

func (t Type) GetType() *types.AbstractType {
	return t.Type.(*types.AbstractType)
}

func (t Type) Write(encoder utils.IUpdateEncoder, offset int) {
	t.GetType().Write(encoder)
}

func (t Type) Gc(store *utils.StructStore) {
	var item = t.GetType().Start
	for item != nil {
		item.Gc(store, true)
		item = item.Right.(*structs.Item)
	}

	t.GetType().Start = nil
	for _, valueItem := range t.GetType().ItemMap {
		for valueItem != nil {
			valueItem.Gc(store, true)
			valueItem = valueItem.Left.(*structs.Item)
		}
	}
	// Clear
	t.GetType().ItemMap = map[string]*structs.Item{}
}

func (t Type) Delete(transaction *utils.Transaction) {
	var item = t.GetType().Start

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
	for _, valueItem := range t.GetType().ItemMap {
		if !valueItem.Deleted {
			valueItem.Delete(transaction)
		} else {
			// Same as above.
			transaction.MergeStructs = append(transaction.MergeStructs, item)
		}
	}
	delete(transaction.Changed, t.GetType())
}

func (t Type) Integrate(transaction *utils.Transaction, item *structs.Item) {
	t.GetType().Integrate(transaction.Doc, item)
}

func (t Type) GetRef() int {
	return TypeRef
}

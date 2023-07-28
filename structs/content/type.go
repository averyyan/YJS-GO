package content

import (
	"YJS-GO/structs"
	"YJS-GO/types"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Type)(nil)

type Type struct {
	Type types.AbstractType
}

func NewType(v any) *Type {
	return &Type{}
}

func (t Type) SetRef(i int) {
	// Do nothing.
}

func ReadType(decoder utils.IUpdateDecoder) (*Type, error) {
	var typeRef = decoder.ReadTypeRef();
	switch (typeRef)
	{
	case YArray.YArrayRefId:
		var arr = YArray.Read(decoder);
		return new ContentType(arr);
	case YMap.YMapRefId:
		var map = YMap.Read(decoder);
		return new ContentType(map);
	case YText.YTextRefId:
		var text = YText.Read(decoder);
		return new ContentType(text);
	default:
		throw new NotImplementedException($"Type {typeRef} not implemented");
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
	var item = Type._start;
	while (item != null)
	{
		item.Gc(store, parentGCd: true);
		item = item.Right as Item;
	}

	Type._start = null;

	foreach (var kvp in Type._map)
	{
	var valueItem = kvp.Value;
	while (valueItem != null)
	{
	valueItem.Gc(store, parentGCd: true);
	valueItem = valueItem.Left as Item;
	}
	}

	Type._map.Clear();
}

func (t Type) Delete(transaction *utils.Transaction) {
	var item = Type._start;

	while (item != null)
	{
		if (!item.Deleted)
		{
			item.Delete(transaction);
		}
		else
		{
			// This will be gc'd later and we want to merge it if possible.
			// We try to merge all deleted items each transaction,
			// but we have no knowledge about that this needs to merged
			// since it is not in transaction. Hence we add it to transaction._mergeStructs.
			transaction._mergeStructs.Add(item);
		}

		item = item.Right as Item;
	}

	foreach (var valueItem in Type._map.Values)
	{
	if (!valueItem.Deleted)
	{
	valueItem.Delete(transaction);
	}
	else
	{
	// Same as above.
	transaction._mergeStructs.Add(valueItem);
	}
	}

	transaction.Changed.Remove(Type);
}

func (t Type) Integrate(transaction *utils.Transaction, item *structs.Item) {
	Type.Integrate(transaction.Doc, item);
}

func (t Type) GetRef() int {
	return TypeRef
}

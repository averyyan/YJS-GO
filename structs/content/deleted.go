package content

import (
	"reflect"

	"YJS-GO/structs"
	"YJS-GO/utils"
)

var _ structs.IContentExt = (*Deleted)(nil)

type Deleted struct {
	Length    uint64
	countable bool
}

func NewDeleted(length int) *Deleted {
	return &Deleted{Length: uint64(length)}
}

func (d Deleted) GetRef() int {
	return DeletedRef
}

func (d Deleted) SetRef(i int) {
	// do nothing
}

func (d Deleted) MergeWith(right structs.IContentExt) bool {
	// Debug.Assert(right is ContentDeleted);
	if reflect.TypeOf(right) != reflect.TypeOf(Deleted{}) {
		return false
	}
	d.Length = uint64(int(d.Length) + right.GetLength())
	return true
}

func ReadDeleted(decoder utils.IUpdateDecoder) (*Deleted, error) {
	var length = decoder.ReadLength()
	return NewDeleted(int(length)), nil
}

func (d Deleted) Copy() structs.IContentExt {
	return NewDeleted(int(d.Length))
}

func (d Deleted) Splice(offset uint64) structs.IContentExt {
	var right = NewDeleted(int(d.Length - offset))
	d.Length = offset
	return right
}

func (d Deleted) GetContent() any {
	// do nothing
	return nil
}

func (d Deleted) GetLength() int {
	return int(d.Length)
}

func (d Deleted) Countable() bool {
	return d.countable
}

func (d Deleted) Write(encoder utils.IUpdateEncoder, offset int) {
	encoder.WriteLength(int(d.Length) - offset)
}

func (d Deleted) Gc(store *utils.StructStore) {
	// Do nothing.
}

func (d Deleted) Delete(transaction *utils.Transaction) {
	// Do nothing.
}

func (d Deleted) Integrate(transaction *utils.Transaction, item *structs.Item) {
	transaction.DeleteSet.Add(item.Id.Client, item.Id.Clock, d.Length)
	item.MarkDeleted()
}

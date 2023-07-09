package utils

import (
	"bytes"
	"io"

	"YJS-GO/structs"
	"github.com/google/uuid"
)

type YDoc struct {
	GC     bool
	Filter GCFilter
	Store  *StructStore
}

type GCFilter struct {
}

var DefaultPredicate *structs.Item = nil
var Store StructStore

type TransactAction func(Transaction)

type YdocOptions struct {
	GC       bool
	GUID     string
	Meta     map[string]string
	AutoLoad bool
}

func (d YdocOptions) Clone() *YdocOptions {
	return &YdocOptions{
		GC:       d.GC,
		GUID:     d.GUID,
		Meta:     d.Meta, // maybe deep clone?
		AutoLoad: d.AutoLoad,
	}
}

func NewDoc() *YDoc {
	return &YDoc{
		GC:     false,
		Filter: GCFilter{},
	}
}

func (d YDoc) EncodeStateVectorV2() []byte {
	var encoder = new(DSEncoderV2)
	writeStateVector(encoder)
	return encoder.ToArray()
}

// EncodeStateAsUpdateV2 Write all the document as a single update message that can be applied on the remote document.
// If you specify the state of the remote client, it will only write the operations that are missing.
// Use 'WriteStateAsUpdate' instead if you are working with Encoder.
func (d YDoc) EncodeStateAsUpdateV2(encodedTargetStateVector []byte) []byte {
	var encoder = NewUpdateEncoderV2()
	var targetStateVector = map[uint64]uint64{}
	if encodedTargetStateVector != nil {
		targetStateVector = DecodeStateVector(bytes.NewReader(encodedTargetStateVector))
	}
	d.WriteStateAsUpdate(encoder, targetStateVector)
	return encoder.ToArray()
}

// WriteStateAsUpdate Write all the document as a single update message.
// If you specify the satte of the remote client, it will only
// write the operations that are missing.
func (d YDoc) WriteStateAsUpdate(encoder IUpdateEncoder, targetStateVector map[uint64]uint64) {
	WriteClientsStructs(encoder, d.Store, targetStateVector)
	NewDeleteSet(d.Store).Write(encoder)
}

func (d YDoc) ApplyUpdateV2(vector []byte, origin interface{}) {
	d.ApplyUpdateV2WithReader(bytes.NewReader(vector), origin)
}

func (d YDoc) ApplyUpdateV2WithReader(reader io.Reader, origin interface{}) {
	var fun = func(tr Transaction) {
		var structDecoder = NewUpdateDecoderV2(reader)
		ReadStructs(structDecoder, tr, Store)
	}
	Transact(fun, origin, false)
}

func Transact(tAction TransactAction, origin interface{}, b bool) {

}

func writeStateVector(encoder IDSEncoder) {
	WriteStateVector(encoder, Store.GetStateVector())

}

func getUUID() string {
	return uuid.NewString()
}

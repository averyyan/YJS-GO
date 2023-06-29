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

func (d YDoc) EncodeStateVectorV2() []byte {
	var encoder = new(DSEncoderV2)
	writeStateVector(encoder)
	return encoder.ToArray()
}

func (d YDoc) EncodeStateAsUpdateV2(b []byte) []byte {
	return b
}

func (d YDoc) ApplyUpdateV2(vector []byte, origin interface{}) {
	d.ApplyUpdateV2WithReader(bytes.NewReader(vector), origin)
}

func (d YDoc) ApplyUpdateV2WithReader(reader io.Reader, origin interface{}) {
	var fun = func(tr Transaction) {
		var structDecoder = UpdateDecoderV2{}
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

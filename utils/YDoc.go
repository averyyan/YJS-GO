package utils

import (
	"bytes"
	"io"
	"math/rand"
	"reflect"

	"YJS-GO/structs"
	"YJS-GO/types"
	"github.com/google/uuid"
)

type YDoc struct {
	GC                      bool
	GCFilter                GCFilter
	Store                   *StructStore
	share                   map[string]*types.AbstractType
	Transaction             *Transaction
	TransactionCleanups     []*Transaction
	BeforeTransaction       func(*Transaction)
	BeforeObserverCalls     func(*Transaction)
	AfterTransaction        func(*Transaction)
	AfterTransactionCleanup func(*Transaction)
	BeforeAllTransactions   func()
	AfterAllTransactions    func([]*Transaction)
	UpdateV2                func(data []byte, origin any, transaction *Transaction)
	Destroyed               func()
	SubdocsChanged          func(Loaded map[*YDoc]struct{}, Added map[*YDoc]struct{}, Removed map[*YDoc]struct{})
	Subdocs                 map[*YDoc]struct{}
	ClientId                uint64
	YDocOptions
	ShouldLoad bool
	// If this document is a subdocument - a document integrated into another document - them _item is defined.
	Item *structs.Item
}

func GenerateNewClientId() uint64 {
	return uint64(rand.Float64())
}

type GCFilter func(*structs.Item) bool

var DefaultPredicate GCFilter = func(item *structs.Item) bool {
	return true
}
var Store StructStore

type TransactAction func(*Transaction)

type YDocOptions struct {
	GC       bool
	GUID     string
	Meta     map[string]string
	AutoLoad bool
}

func ReadYDocOptions(decoder IUpdateDecoder) *YDocOptions {
	var dict = (decoder.ReadAny()).(map[string]any)

	var result = &YDocOptions{}
	var ok bool
	result.GC, ok = dict["gc"].(bool)
	if ok {
		result.GC = true
	}
	result.GUID, ok = dict["guid"].(string)
	if ok {
		result.GUID = getUUID()
	}
	result.Meta, ok = dict["meta"].(map[string]string)
	if ok {
		result.Meta = nil
	}
	result.AutoLoad, ok = dict["autoLoad"].(bool)
	if ok {
		result.AutoLoad = false
	}
	return result
}

func (d *YDocOptions) Clone() *YDocOptions {
	return &YDocOptions{
		GC:       d.GC,
		GUID:     d.GUID,
		Meta:     d.Meta, // maybe deep clone?
		AutoLoad: d.AutoLoad,
	}
}

func NewDoc(options *YDocOptions) *YDoc {
	doc := &YDoc{}
	if options == nil {
		options = &YDocOptions{}
	}
	doc.TransactionCleanups = []*Transaction{}
	doc.ClientId = GenerateNewClientId()
	doc.share = map[string]*types.AbstractType{}
	doc.Store = &StructStore{}
	doc.Subdocs = map[*YDoc]struct{}{}
	doc.ShouldLoad = options.AutoLoad
	return doc
}

func NewDocWithOption(options *YDocOptions) *YDoc {
	return &YDoc{
		GC:       options.GC,
		GCFilter: DefaultPredicate,
	}
}

func (d YDocOptions) Write(encoder IUpdateEncoder, offset int) {
	var dict = map[string]any{}
	dict["gc"] = d.GC
	dict["guid"] = d.GUID
	dict["autoLoad"] = d.AutoLoad
	if d.Meta != nil {
		dict["meta"] = d.Meta
	}
	encoder.WriteAny(dict)
}

func (d *YDoc) EncodeStateVectorV2() []byte {
	var encoder = new(DSEncoderV2)
	writeStateVector(encoder)
	return encoder.ToArray()
}

// EncodeStateAsUpdateV2 Write all the document as a single update message that can be applied on the remote document.
// If you specify the state of the remote client, it will only write the operations that are missing.
// Use 'WriteStateAsUpdate' instead if you are working with Encoder.
func (d *YDoc) EncodeStateAsUpdateV2(encodedTargetStateVector []byte) []byte {
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
func (d *YDoc) WriteStateAsUpdate(encoder IUpdateEncoder, targetStateVector map[uint64]uint64) {
	WriteClientsStructs(encoder, d.Store, targetStateVector)
	NewDeleteSet(d.Store).Write(encoder)
}

func (d *YDoc) ApplyUpdateV2(vector []byte, origin interface{}, local ...bool) {
	d.ApplyUpdateV2WithReader(bytes.NewReader(vector), origin, local...)
}

func (d *YDoc) ApplyUpdateV2WithReader(reader io.Reader, origin interface{}, local ...bool) {
	var fun = func(tr *Transaction) {
		var structDecoder = NewUpdateDecoderV2(reader)
		ReadStructs(structDecoder, tr, d.Store)
		Store.ReadAndApplyDeleteSet(structDecoder, tr)
	}
	d.Transact(fun, origin, local...)
}

func (d *YDoc) Get(key string) *types.AbstractType {
	var (
		Type *types.AbstractType
		ok   bool
	)
	Type, ok = d.share[key]
	if !ok {
		Type = &types.AbstractType{}
	}
	// if (typeof(T) != typeof(AbstractType) && !typeof(T).IsAssignableFrom(type.GetType()))
	if reflect.TypeOf(Type) != reflect.TypeOf(&types.AbstractType{}) {
		if reflect.TypeOf(Type) != reflect.TypeOf(&types.AbstractType{}) {
			t := &types.AbstractType{}
			t.ItemMap = Type.ItemMap
			for _, v := range Type.ItemMap {
				for ; v != nil; v = v.Left.(*structs.Item) {
					v.Parent = t
				}
			}
			t.Start = Type.Start
			for n := t.Start; n != nil; n = n.Right.(*structs.Item) {
				n.Parent = t
			}
			t.Length = Type.Length
			d.share[key] = t
			t.Integrate(d, nil)
		} else {
			//   throw new Exception($"Type with the name {name} has already been defined with a different constructor");
			return nil
		}
	}
	return Type

}

func (d *YDoc) Transact(tAction TransactAction, origin interface{}, local ...bool) {
	var initialCall = false
	if d.Transaction == nil {
		initialCall = true
		d.Transaction = NewTransaction(d, origin, local...)
		d.TransactionCleanups = append(d.TransactionCleanups, d.Transaction)
		if len(d.TransactionCleanups) == 1 {
			if d.BeforeAllTransactions != nil {
				d.BeforeAllTransactions()
			}
		}
		if d.BeforeTransaction != nil {
			d.BeforeTransaction(d.Transaction)
		}
	}

	tAction(d.Transaction)
	if initialCall && d.TransactionCleanups[0] == d.Transaction {
		// The first transaction ended, now process observer calls.
		// Observer call may create new transacations for which we need to call the observers and do cleanup.
		// We don't want to nest these calls, so we execute these calls one after another.
		// Also we need to ensure that all cleanups are called, even if the observers throw errors.
		CleanupTransactions(d.TransactionCleanups, 0)
	}
}

func (d *YDoc) CloneOptionsWithNewGuid() *YDocOptions {
	newOpts := d.YDocOptions.Clone()
	newOpts.GUID = getUUID()
	return newOpts
}

func writeStateVector(encoder IDSEncoder) {
	WriteStateVector(encoder, Store.GetStateVector())

}

func getUUID() string {
	return uuid.NewString()
}

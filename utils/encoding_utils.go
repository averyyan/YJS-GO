package utils

import (
	"encoding/binary"
	"errors"
	"io"
	"sort"

	"YJS-GO/lib0"
	"YJS-GO/structs"
	"YJS-GO/structs/content"
)

// ReadItemcontent. We use first five bits in the info flag for determining the type of the struct.
// 0: GC
// 1: Deleted content.
// 2: JSON content.
// 3: Binary content.
// 4: String content.
// 5: Embed content. (for richtext content.)
// 6: Format content. (a formatting marker for richtext content.)
// 7: Type content.
// 8: Any content.
// 9: Doc content.
func ReadItemContent(decoder IUpdateDecoder, info byte) (structs.IContent, error) {
	switch uint(info) & lib0.Bits5 {
	case 0: // GC
		return nil, errors.New("GC is not Itemcontent.")
	case 1: // Deleted
		return content.ReadDeleted(decoder)
	case 2: // JSON
		return content.ReadJson(decoder)
	case 3: // Binary
		return content.ReadBinary(decoder)
	case 4: // String
		return content.ReadString(decoder)
	case 5: // Embed
		return content.ReadEmbed(decoder)
	case 6: // Format
		return content.ReadFormat(decoder)
	case 7: // Type
		return content.ReadType(decoder)
	case 8: // Any
		return content.ReadAny(decoder)
	case 9: // Doc
		return content.ReadDoc(decoder)
	}
	return nil, errors.New("dont implement type")
}

func WriteStateVector(encoder IDSEncoder, sv map[uint64]uint64) {
	binary.Write(encoder.RestWriter(), binary.LittleEndian, len(sv))
	for k, v := range sv {
		binary.Write(encoder.RestWriter(), binary.LittleEndian, k)
		binary.Write(encoder.RestWriter(), binary.LittleEndian, v)
	}
}

// ReadStructs Read the next ItemId in a Decoder and fill this ItemId with the read data.
// <br/>
// This is called when data is received from a remote peer.
// public static void ReadStructs(IUpdateDecoder decoder, transaction *Transaction, StructStore store)
// {
// var clientStructRefs = ReadClientStructRefs(decoder, transaction.Doc);
// store.MergeReadStructsIntoPendingReads(clientStructRefs);
// store.ResumeStructIntegration(transaction);
// store.CleanupPendingStructs();
// store.TryResumePendingDeleteReaders(transaction);
// }
func ReadStructs(decoder IUpdateDecoder, transaction *Transaction, store *StructStore) {
	var clientStructRefs = ReadClientStructRefs(decoder, transaction.Doc)
	store.MergeReadStructsIntoPendingReads(clientStructRefs)
	store.ResumeStructIntegration(transaction)
	store.CleanupPendingStructs()
	store.TryResumePendingDeleteReaders(transaction)
}

func ReadClientStructRefs(decoder IUpdateDecoder, doc *YDoc) map[uint64][]structs.IAbstractStruct {
	var clientRefs = map[uint64][]structs.IAbstractStruct{}
	numOfStateUpdates, err := binary.ReadUvarint(decoder.Reader())
	if err != nil {
		return nil
	}
	for i := 0; i < int(numOfStateUpdates); i++ {
		numberOfStructs, err := binary.ReadUvarint(decoder.Reader())
		if err != nil {
			return nil
		}
		var refs = make([]structs.IAbstractStruct, numberOfStructs)
		var client = decoder.ReadClient()
		clock, err := binary.ReadUvarint(decoder.Reader())
		if err != nil {
			return nil
		}
		clientRefs[client] = refs
		for i := 0; i < int(numberOfStructs); i++ {
			var info = decoder.ReadInfo()
			if lib0.Bits5&uint(info) != 0 {
				var (
					leftOrigin    *ID
					rightOrigin   *ID
					hasParentYKey bool
				)
				// The item that was originally to the left of this item.
				if (uint(info) & lib0.Bit8) == lib0.Bit8 {
					leftOrigin = decoder.ReadLeftId()
				}
				// The item that was originally to the right of this item.
				if (uint(info) & lib0.Bit7) == lib0.Bit7 {
					leftOrigin = decoder.ReadRightId()
				}
				var cantCopyParentInfo = (uint(info) & (lib0.Bit7 | lib0.Bit8)) == 0
				if cantCopyParentInfo {
					hasParentYKey = decoder.ReadParentInfo()
				}
				// If parent == null and neither left nor right are defined, then we know that 'parent' is child of 'y'
				// and we read the next string as parentYKey.
				// It indicates how we store/retrieve parent from 'y.share'.
				var parentYKey string
				if cantCopyParentInfo && hasParentYKey {
					parentYKey = decoder.ReadString()
				}
				var parent interface{} = doc.Get(parentYKey)
				if cantCopyParentInfo && !hasParentYKey {
					parent = decoder.ReadLeftId()
				}
				content, _ := ReadItemContent(decoder, info)

				var parentSub string
				if cantCopyParentInfo && uint(info)&lib0.Bit6 == lib0.Bit6 {
					parentSub = decoder.ReadString()
				}
				var str = &structs.Item{
					Id: &ID{
						Client: client,
						Clock:  clock,
					},
					LeftOrigin:  leftOrigin,
					RightOrigin: rightOrigin,
					Parent:      parent,
					ParentSub:   parentSub,
					Content:     content,
				}
				refs = append(refs, str)
				clock += str.Length

			} else {
				length := decoder.ReadLength()
				gc := &structs.GC{
					Id:     &ID{client, clock},
					Length: length,
				}
				refs = append(refs, gc)
				clock += length
			}
		}

	}
	return clientRefs
}

func ReadStateVector(decoder IDSDecoder) map[uint64]uint64 {
	ssLength, err := binary.ReadUvarint(decoder.Reader())
	if err != nil {
		return nil
	}
	var ss = make(map[uint64]uint64, ssLength)
	for i := 0; i < int(ssLength); i++ {
		client, err := binary.ReadUvarint(decoder.Reader())
		if err != nil {
			continue
		}
		clock, err := binary.ReadUvarint(decoder.Reader())
		if err != nil {
			continue
		}
		ss[client] = clock
	}

	return ss
}

func DecodeStateVector(input io.Reader) map[uint64]uint64 {
	return ReadStateVector(NewDsDecoderV2(input))
}

func WriteClientsStructs(encoder IUpdateEncoder, store *StructStore, _sm map[uint64]uint64) {
	// We filter all valid _sm entries into sm.
	var sm map[uint64]uint64
	for client, clock := range _sm {
		// Only write if new structs are available.
		if store.GetState(client) > clock {
			sm[client] = clock
		}
	}

	for client, _ := range store.GetStateVector() {
		if _, ok := sm[client]; !ok {
			sm[client] = 0
		}
	}
	// Write # states that were updated.
	binary.Write(encoder.RestWriter(), binary.LittleEndian, uint(len(sm)))
	// Write items with higher client ids first.
	// This heavily improves the conflict resolution algorithm.
	sortedClients := sortClients(sm)
	for _, client := range sortedClients {
		WriteStructs(encoder, store.Clients[client], client, sm[client])
	}

	encoder.RestWriter()
}

func WriteStructs(encoder IUpdateEncoder, structs []structs.IAbstractStruct, client, clock uint64) {
	// Write first id.
	startNewStructs := FindIndexSS(structs, clock)

	// Write # encoded structs.
	binary.Write(encoder.RestWriter(), binary.LittleEndian, uint(len(structs))-startNewStructs)
	encoder.WriteClient(client)
	binary.Write(encoder.RestWriter(), binary.LittleEndian, uint(clock))

	// Write first struct with offset.
	var firstStruct = structs[startNewStructs]
	firstStruct.Write(encoder, (int)(clock-firstStruct.ID().Clock))

	for i := startNewStructs + 1; i < uint(len(structs)); i++ {
		structs[i].Write(encoder, 0)
	}
}

func sortClients(sm map[uint64]uint64) []uint64 {
	var a []uint64
	for _, u := range sm {
		a = append(a, u)
	}
	sort.SliceStable(a, func(i, j int) bool {
		return a[i] < a[j]
	})

	return a
}

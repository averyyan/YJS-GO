package utils

import (
	"container/list"
	"encoding/binary"
	"errors"

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

//    public static void WriteStateVector(IDSEncoder encoder, IDictionary<long, long> sv)
//        {
//            encoder.RestWriter.WriteVarUint((uint)sv.Count);
//
//            foreach (var kvp in sv)
//            {
//                var Client = kvp.Key;
//                var Clock = kvp.Value;
//
//                encoder.RestWriter.WriteVarUint((uint)Client);
//                encoder.RestWriter.WriteVarUint((uint)Clock);
//            }
//        }

func WriteStateVector(encoder IDSEncoder, sv map[uint64]uint64) {
	binary.Write(encoder.RestWriter(), binary.LittleEndian, len(sv))
	for k, v := range sv {
		binary.Write(encoder.RestWriter(), binary.LittleEndian, k)
		binary.Write(encoder.RestWriter(), binary.LittleEndian, v)
	}
}

// / <summary>
// / Read the next Item in a Decoder and fill this Item with the read data.
// / <br/>
// / This is called when data is received from a remote peer.
// / </summary>
// public static void ReadStructs(IUpdateDecoder decoder, Transaction transaction, StructStore store)
// {
// var clientStructRefs = ReadClientStructRefs(decoder, transaction.Doc);
// store.MergeReadStructsIntoPendingReads(clientStructRefs);
// store.ResumeStructIntegration(transaction);
// store.CleanupPendingStructs();
// store.TryResumePendingDeleteReaders(transaction);
// }
func ReadStructs(decoder IUpdateDecoder, transaction Transaction, store StructStore) {
	var clientStructRefs = ReadClientStructRefs(decoder, transaction.Doc)
	store.MergeReadStructsIntoPendingReads(clientStructRefs)
	store.ResumeStructIntegration(transaction)
	store.CleanupPendingStructs()
	store.TryResumePendingDeleteReaders(transaction)
}

func ReadClientStructRefs(decoder IUpdateDecoder, doc YDoc) map[int]list.List {
	return nil
}

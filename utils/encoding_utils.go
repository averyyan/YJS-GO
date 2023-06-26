package utils

import (
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

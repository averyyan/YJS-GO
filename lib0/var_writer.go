package lib0

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

// WriteUint16  Writes two bytes as an unsigned unteger.
func WriteUint16(writer *bufio.Writer, num uint) {
	binary.Write(writer, binary.LittleEndian, num&Bits8)
	writer.WriteByte((byte)(num & Bits8))
	writer.WriteByte((byte)((num >> 8) & Bits8))
}

// WriteUint32  Writes four bytes as an unsigned integer.
func WriteUint32(writer *bufio.Writer, num uint) {
	for i := 0; i < 4; i++ {
		writer.WriteByte((byte)(num & Bits8))
		num >>= 8
	}
}

// WriteVarUint Writes a variable length unsigned integer.
// Encodes integers in the range <c>[0, 4294967295] / [0, 0xFFFFFFFF]</c>.
func WriteVarUint(writer *bufio.Writer, num uint) {
	for num > Bits7 {
		writer.WriteByte((byte)(Bit8 | (Bits7 & num)))
		num >>= 7
	}

	writer.WriteByte((byte)(Bits7 & num))
}

// WriteVarInt Writes a variable length integer.
// <br/>
// Encodes integers in the range <c>[-2147483648, -2147483647]</c>.
// <br/>
// We don't use zig-zag encoding because we want to keep the option open
// to use the same function for <c>BigInt</c> and 53-bit integers (doubles).
// <br/>
//
//	We use the 7-th bit instead for signalling that this is a negative number.
func WriteVarInt(writer *bufio.Writer, num int, treatZeroAsNegative bool) {
	var isNegative = treatZeroAsNegative
	if num < 0 {
		isNegative = true
	}
	if isNegative {
		num = -num
	}

	// |   whether to continue reading   |         is negative         | value.
	var tmpA uint
	var tmpB uint
	if uint(num) > Bits6 {
		tmpA = Bit8
	}
	if isNegative {
		tmpB = Bit7
	}
	// |   whether to continue reading   |         is negative         | value.
	writer.WriteByte((byte)(tmpA | tmpB | Bit6&uint(num)))
	num >>= 6

	// We don't need to consider the case of num == 0 so we can use a different pattern here than above.
	var tmpC uint
	if uint(num) > Bits7 {
		tmpC = Bit8
	}
	for num > 0 {
		writer.WriteByte((byte)(tmpC | (Bits7 & uint(num))))
		num >>= 7
	}
}

// WriteVarString  Writes a variable length string.
func WriteVarString(writer *bufio.Writer, str string) {
	var data = []byte(str)
	WriteVarUint8Array(writer, data)
}

// WriteVarUint8Array  Appends a byte array to the writer.
func WriteVarUint8Array(writer *bufio.Writer, array []byte) {
	WriteVarUint(writer, uint(len(array)))
	writer.Write(array)
}

// WriteAny Encodes data with efficient binary format.
// <br/>
// Differences to JSON:
//   - Transforms data to a binary format (not to a string).
//   - Encodes undefined, NaN, and ArrayBuffer (these can't be represented in JSON).
//   - Numbers are efficiently encoded either as a variable length integer, as a 32-bit
//     float, or as a 64-bit float.
//
// <br/>
// Encoding table:
// | Data Type                      | Prefix | Encoding method   | Comment                                                        |
// | ------------------------------ | ------ | ----------------- | -------------------------------------------------------------- |
// | undefined                      | 127    |                   | Functions, symbol, and everything that cannot be identified    |
// |                                |        |                   | is encdoded as undefined.                                      |
// | null                           | 126    |                   |                                                                |
// | integer                        | 125    | WriteVarInt       | Only encodes 32-bit signed integers.                           |
// | float                          | 124    | SingleToInt32Bits |                                                                |
// | double                         | 123    | DoubleToInt64Bits |                                                                |
// | boolean (false)                | 121    |                   | True and false are different data types so we save the         |
// | boolean (true)        |        | 120    |                   | following byte (0b_01111000) so the last bit determines value. |
// | string                         | 119    | WriteVarString    |                                                                |
// | IDictionary&lt;string, any&gt; | 118    | custom            | Writes length, then key-value pairs.                           |
// | ICollection&lt;any&gt;         | 117    | custom            | Writes length, then values.                                    |
// | byte[]                         | 116    |                   | We use byte[] for any kind of binary data.                     |
// <br/>
// Reasons for the decreasing prefix:
// We need the first bit for extendability (later we may want to encode the prefix with <see cref="WriteVarUint(BinaryWriter, uint)"/>).
// The remaining 7 bits are divided as follows:
// [0-30]   The beginning of the data range is used for custom purposes
//
//	      (defined by the function that uses this library).
//	[31-127] The end of the data range is used for data encoding.
func WriteAny(writer *bufio.Writer, o any) {
	switch t := o.(type) {
	case string: // TYPE 119: STRING
		writer.WriteByte(119)
		WriteVarString(writer, t)
	case bool: // TYPE 120/121: boolean (true/false)
		var b uint = 121
		if t {
			b = 120
		}
		writer.WriteByte((byte)(b))
	case float64: // TYPE 123: FLOAT64
		bits := math.Float64bits(t)
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, bits)
		if !IsLittleEndian() {
			bytes = Reverse(bytes)
		}
		writer.WriteByte(123)
		writer.Write(bytes)
	case float32: // TYPE 124: FLOAT32
		bits := math.Float32bits(t)
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, bits)
		if !IsLittleEndian() {
			bytes = Reverse(bytes)
		}
		writer.WriteByte(124)
		writer.Write(bytes)
	case int: // TYPE 125: INTEGER
		writer.WriteByte(125)
		WriteVarInt(writer, t, false)
	case int64: // Special case: treat LONG as INTEGER.
		writer.WriteByte(125)
		WriteVarInt(writer, int(t), false)
	case nil: // TYPE 126: null
		// TYPE 127: undefined
		writer.WriteByte(126)
	case []byte: // TYPE 116: ArrayBuffer
		writer.WriteByte(116)
		WriteVarUint8Array(writer, t)
	case map[string]any: // TYPE 118: object (Dictionary<string, object>)
		writer.WriteByte(118)
		WriteVarUint(writer, uint(len(t)))
		for key, value := range t {
			WriteVarString(writer, key)
			WriteAny(writer, value)
		}
	case []any: // TYPE 117: Array
		writer.WriteByte(117)
		WriteVarUint(writer, uint(len(t)))
		for _, item := range t {
			WriteAny(writer, item)
		}

	default:
		// throw new NotSupportedException($"Unsupported object type: {o?.GetType()}")
		fmt.Printf("Unsupported object type:%v", reflect.TypeOf(o))
	}
}

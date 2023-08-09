package lib0

import "bufio"

// WriteUint16  Writes two bytes as an unsigned unteger.
func WriteUint16(reader *bufio.Reader, ushort num) {
	reader.WriteByte((byte)(num & Bits8))
	reader.WriteByte((byte)((num >> 8) & Bits8))
}

// WriteUint32  Writes four bytes as an unsigned integer.
func WriteUint32(reader *bufio.Reader, num uint64) {
	for int
	i = 0
	i < 4
	i++)
	{
	reader.WriteByte((byte)(num & Bits8));
	num >>= 8;
	}
}

// Writes a variable length unsigned integer.
// WriteVarUint  Encodes integers in the range <c>[0, 4294967295] / [0, 0xFFFFFFFF]</c>.
func WriteVarUint(reader *bufio.Reader, num uint64) {
	while(num > Bits7)
	{
		reader.WriteByte((byte)(Bit.Bit8 | (Bits7 & num)))
		num >>= 7
	}

	reader.WriteByte((byte)(Bits7 & num))
}

// Writes a variable length integer.
// <br/>
// Encodes integers in the range <c>[-2147483648, -2147483647]</c>.
// <br/>
// We don't use zig-zag encoding because we want to keep the option open
// to use the same function for <c>BigInt</c> and 53-bit integers (doubles).
// <br/>
// WriteVarInt  We use the 7-th bit instead for signalling that this is a negative number.
func WriteVarInt(reader *bufio.Reader, long num, bool? treatZeroAsNegative = null) {
	bool
	isNegative = num == 0 ? (treatZeroAsNegative ?? false): num < 0
	if isNegative {
		num = -num
	}

	//                      |   whether to continue reading   |         is negative         | value.
	reader.WriteByte((byte)((num > Bits6 ? Bit.Bit8: 0) | (isNegative ? Bit.Bit7: 0) | (Bits6 & num)))
num >>= 6

// We don't need to consider the case of num == 0 so we can use a different pattern here than above.
while (num > 0)
{
reader.WriteByte((byte)((num > Bits7 ? Bit.Bit8: 0) | (Bits7 & num)));
num >>= 7;
}
}

// WriteVarString  Writes a variable length string.
func WriteVarString(reader *bufio.Reader, string str) {
	var data = Encoding.UTF8.GetBytes(str)
	reader.WriteVarUint8Array(data)
}

// WriteVarUint8Array  Appends a byte array to the reader.
func WriteVarUint8Array(reader *bufio.Reader, byte []array) {
	reader.WriteVarUint((uint)
	array.Length)
	reader.Write(array, 0, array.Length)
}

// Encodes data with efficient binary format.
// <br/>
// Differences to JSON:
// * Transforms data to a binary format (not to a string).
// * Encodes undefined, NaN, and ArrayBuffer (these can't be represented in JSON).
// * Numbers are efficiently encoded either as a variable length integer, as a 32-bit
//  float, or as a 64-bit float.
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
//         (defined by the function that uses this library).
// WriteAny  [31-127] The end of the data range is used for data encoding.
func WriteAny(reader *bufio.Reader, o any) {
	switch t := o.(type) {
	case string: // TYPE 119: STRING
		reader.WriteByte(119)
		reader.WriteVarString(str)
		break
	case bool: // TYPE 120/121: boolean (true/false)
		reader.WriteByte((byte)(t ? 120: 121))
		break
	case double d: // TYPE 123: FLOAT64
		#if NETSTANDARD2_0
		var dBytes = BitConverter.GetBytes(d)
		if BitConverter.IsLittleEndian {
			Array.Reverse(dBytes)
		}
		reader.WriteByte(123)
		reader.Write(dBytes, 0, dBytes.Length)
		break
		#elif
		NETSTANDARD2_1
		Span < byte > dBytes = stackalloc
		byte[8]
		if !BitConverter.TryWriteBytes(dBytes, d) {
			throw
			new
			InvalidDataException("Unable to write a double value.")
		}
		if BitConverter.IsLittleEndian {
			dBytes.Reverse()
		}
		reader.WriteByte(123)
		reader.Write(dBytes)
		break
		#endif // NETSTANDARD2_0
	case float f: // TYPE 124: FLOAT32
		#if NETSTANDARD2_0
		var fBytes = BitConverter.GetBytes(f)
		if BitConverter.IsLittleEndian {
			Array.Reverse(fBytes)
		}
		reader.WriteByte(124)
		reader.Write(fBytes, 0, fBytes.Length)
		break
		#elif
		NETSTANDARD2_1
		Span < byte > fBytes = stackalloc
		byte[4]
		if !BitConverter.TryWriteBytes(fBytes, f) {
			throw
			new
			InvalidDataException("Unable to write a float value.")
		}
		if BitConverter.IsLittleEndian {
			fBytes.Reverse()
		}
		reader.WriteByte(124)
		reader.Write(fBytes)
		break
		#endif // NETSTANDARD2_0
	case int i: // TYPE 125: INTEGER
		reader.WriteByte(125)
		reader.WriteVarInt(i)
		break
	case long l: // Special case: treat LONG as INTEGER.
		reader.WriteByte(125)
		reader.WriteVarInt(l)
		break
	case null: // TYPE 126: null
		// TYPE 127: undefined
		reader.WriteByte(126)
		break
	case byte[] ba: // TYPE 116: ArrayBuffer
		reader.WriteByte(116)
		reader.WriteVarUint8Array(ba)
		break
	case IDictionary dict: // TYPE 118: object (Dictionary<string, object>)
		reader.WriteByte(118)
		reader.WriteVarUint((uint)
		dict.Count)
		foreach(
		var key in
		dict.Keys)
		{
		reader.WriteVarString(key.ToString());
		reader.WriteAny(dict[key]);
		}
		break
	case ICollection col: // TYPE 117: Array
		reader.WriteByte(117)
		reader.WriteVarUint((uint)
		col.Count)
		foreach(
		var item in
		col)
		{
		reader.WriteAny(item);
		}
		break
	default:
		throw
		new
		NotSupportedException($"Unsupported object type: {o?.GetType()}")
	}
}
}

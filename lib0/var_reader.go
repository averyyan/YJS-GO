package lib0

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
)

type VarReader struct {
	*bufio.Reader
}

func New(reader io.Reader) *VarReader {
	return &VarReader{Reader: bufio.NewReader(reader)}
}

func ReadVarString(reader *bufio.Reader) string {
	remainingLen, err := binary.ReadUvarint(reader)
	if err != nil || remainingLen == 0 {
		return ""
	}
	readBytes := make([]byte, remainingLen)
	n, err := reader.Read(readBytes)
	if err != nil {
		return ""
	}
	if uint64(n) != remainingLen {
		return ""
	}

	// reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
	// strBytes, _ = ioutil.ReadAll(reader)
	// return string(strBytes)
	return string(bytes.Runes(readBytes))
}

func ReadVarUint8ArrayAsStream(reader *bufio.Reader) *bufio.Reader {
	var data = ReadVarUint8Array(reader)
	return bufio.NewReader(bytes.NewReader(data))
}

func ReadVarUint8Array(reader *bufio.Reader) []byte {
	// uint len = stream.ReadVarUint();
	length, err := binary.ReadUvarint(reader)
	if err != nil {
		return nil
	}
	var ret = make([]byte, length)
	n, err := reader.Read(ret)
	if err != nil || uint64(n) != length {
		return nil
	}
	return ret
}

// Contains <see cref="Stream"/> extensions compatible with the <c>lib0</c>:
// <see href="https://github.com/dmonad/lib0"/>.

// Reads two bytes as an unsigned integer.

func ReadUint16(reader *bufio.Reader) ushort {
	return (ushort)(stream._ReadByte() + (stream._ReadByte() << 8))
}

// Reads four bytes as an unsigned integer.

func ReadUint32(reader *bufio.Reader) uint {
	return (uint)((stream._ReadByte() + (stream._ReadByte() << 8) + (stream._ReadByte() << 16) + (stream._ReadByte() << 24)) >> 0)
}

// Reads unsigned integer (32-bit) with variable length.
// 1/8th of the storage is used as encoding overhead.
// * Values &lt; 2^7 are stored in one byte.
// * Values &lt; 2^14 are stored in two bytes.

// <exception cref="InvalidDataException">Invalid binary format.</exception>

func ReadVarUint(reader *bufio.Reader) uint {
	uint
	num = 0
	int
	len = 0

	while(true)
	{
		byte
		r = stream._ReadByte()
		num |= (r & Bits.Bits7) << len
		len += 7

		if r < Bit.Bit8 {
			return num
		}

		if len > 35 {
			throw
			new
			InvalidDataException("Integer out of range.")
		}
	}
}

// Reads a 32-bit variable length signed integer.
// 1/8th of storage is used as encoding overhead.
// * Values &lt; 2^7 are stored in one byte.
// * Values &lt; 2^14 are stored in two bytes.

// <exception cref="InvalidDataException">Invalid binary format.</exception>

func ReadVarInt(reader *bufio.Reader) (long Value, int Sign) {
	byte
	r = stream._ReadByte()
	uint
	num = r & Bits.Bits6
	int
	len = 6
	int
	sign = (r & Bit.Bit7) > 0 ? -1: 1

	if (r & Bit.Bit8) == 0 {
		// Don't continue reading.
		return sign * num, sign)
	}

	while(true)
	{
		r = stream._ReadByte()
		num |= (r & Bits.Bits7) << len
		len += 7

		if r < Bit.Bit8 {
			return sign * num, sign)
		}

		if len > 41 {
			throw
			new
			InvalidDataException("Integer out of range")
		}
	}
}

// Reads a variable length string.

// <remarks>
// <see cref="StreamEncodingExtensions.WriteVarUint(Stream, uint)"/> is used to store the length of the string.
// </remarks>

func ReadVarString(reader *bufio.Reader) string {
	uint
	remainingLen = stream.ReadVarUint()
	if remainingLen == 0 {
		return string.Empty
	}

	var data = stream._ReadBytes(int
	remainingLen)

	var str = Encoding.UTF8.GetString(data)
	return str
}

// Reads a variable length byte array.

func ReadVarUint8Array(reader *bufio.Reader) byte[] {
	uint
	len = stream.ReadVarUint()
	return stream._ReadBytes(int
	len)
}

// Reads variable length byte array as a readable <see cref="MemoryStream"/>.

func ReadVarUint8ArrayAsStream(reader *bufio.Reader) MemoryStream {
	var data = stream.ReadVarUint8Array()
	return new
	MemoryStream(data, writable: false)
}

// Decodes data from the stream.

func ReadAny(reader *bufio.Reader) any {
	byte
	type = stream._ReadByte()
	switch

	type
)
	{
	case 119: // String
	return stream.ReadVarString();
	case 120: // boolean true
	return true;
	case 121: // boolean false
	return false;
	case 123: // Float64

	var dBytes = new byte[8];
	stream._ReadBytes(dBytes);

	if (BitConverter.IsLittleEndian)
	{
	Array.Reverse(dBytes);
	}

	return BitConverter.ToDouble(dBytes, 0);

	Span<byte> dBytes = stackalloc byte[8];
	stream._ReadBytes(dBytes);

	if (BitConverter.IsLittleEndian)
	{
	dBytes.Reverse();
	}

	return BitConverter.ToDouble(dBytes);

	case 124: // Float32

	var fBytes = new byte[4];
	stream._ReadBytes(fBytes);

	if (BitConverter.IsLittleEndian)
	{
	Array.Reverse(fBytes);
	}

	return BitConverter.ToSingle(fBytes, 0);

	Span<byte> fBytes = stackalloc byte[4];
	stream._ReadBytes(fBytes);

	if (BitConverter.IsLittleEndian)
	{
	fBytes.Reverse();
	}

	return BitConverter.ToSingle(fBytes);

	case 125: // integer
	return (int)stream.ReadVarInt().Value;
	case 126: // null
	case 127: // undefined
	return null;
	case 116: // ArrayBuffer
	return stream.ReadVarUint8Array();
	case 117: // Array<any>
	{
	var len = (int)stream.ReadVarUint();
	var arr = new List<any>(len);

	for (int i = 0; i < len; i++)
	{
	arr.Add(stream.ReadAny());
	}

	return arr;
	}
	case 118: // any (Dictionary<string, any>)
	{
	var len = (int)stream.ReadVarUint();
	var obj = new Dictionary<string, any>(len);

	for (int i = 0; i < len; i++)
	{
	var key = stream.ReadVarString();
	obj[key] = stream.ReadAny();
	}

	return obj;
	}
	default:
	throw new InvalidDataException($"Unknown any type: {type}");
	}
}

// Reads a byte from the stream and advances the position within the stream by one byte.

// <exception cref="EndOfStreamException">End of stream reached.</exception>

func _ReadByte(reader *bufio.Reader) byte {
	int
	v = stream.ReadByte()
	if v < 0 {
		throw
		new
		EndOfStreamException()
	}

	return Convert.ToByte(v)
}

// Reads a sequence of bytes from the current stream and advances the position
// within the stream by the number of bytes read.

// <exception cref="EndOfStreamException">End of stream reached.</exception>

func _ReadBytes(reader *bufio.Reader, int count) byte[] {
	Debug.Assert(count >= 0)

	var result = new
	byte[count]
	for int
	i = 0
	i < count
	i++)
	{
	int v = stream.ReadByte();
	if (v < 0)
	{
	throw new EndOfStreamException();
	}

	result[i] = Convert.ToByte(v);
	}

	return result
}

// Reads a sequence of bytes from the current stream and advances the position
// within the stream by the number of bytes read.

// <exception cref="EndOfStreamException">End of stream reached.</exception>

func _ReadBytes(reader *bufio.Reader, byte []buffer) void {
	if buffer.Length == 0 {
		return
	}

	if buffer.Length != stream.Read(buffer, 0, buffer.Length) {
		throw
		new
		EndOfStreamException()
	}
}

func _ReadBytes(reader *bufio.Reader, Span<byte> buffer) void {
	if buffer.Length == 0 {
		return
	}

	if buffer.Length != stream.Read(buffer) {
		throw
		new
		EndOfStreamException()
	}
}

}

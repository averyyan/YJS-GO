package lib0

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"unsafe"
)

type VarReader struct {
	*bufio.Reader
}

func New(reader io.Reader) *VarReader {
	return &VarReader{Reader: bufio.NewReader(reader)}
}

// ReadVarString Reads a variable length string.
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

// ReadVarUint8ArrayAsReader Reads variable length byte array as a readable <see cref="Memoryreader"/>.
func ReadVarUint8ArrayAsReader(reader *bufio.Reader) *bufio.Reader {
	var data = ReadVarUint8Array(reader)
	return bufio.NewReader(bytes.NewReader(data))
}

// ReadVarUint8Array Reads a variable length byte array.
func ReadVarUint8Array(reader *bufio.Reader) []byte {
	// uint len = reader.ReadVarUint();
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

// ReadUint16 Contains <see cref="reader"/> extensions compatible with the <c>lib0</c>:
// <see href="https://github.com/dmonad/lib0"/>.
// Reads two bytes as an unsigned integer.
func ReadUint16(reader *bufio.Reader) uint16 {
	a, errA := reader.ReadByte()
	b, errB := reader.ReadByte()
	if errA != nil || errB != nil {
		return 0
	}
	return (uint16)(a + b<<8)
}

// ReadUint32 Reads four bytes as an unsigned integer.
func ReadUint32(reader *bufio.Reader) uint32 {
	a, errA := reader.ReadByte()
	b, errB := reader.ReadByte()
	c, errC := reader.ReadByte()
	d, errD := reader.ReadByte()
	if errA != nil || errB != nil || errC != nil || errD != nil {
		return 0
	}
	return (uint32)((a + (b << 8) + (c << 16) + (d << 24)) >> 0)
}

// ReadVarUint Reads unsigned integer (32-bit) with variable length.
// 1/8th of the storage is used as encoding overhead.
// * Values &lt; 2^7 are stored in one byte.
// * Values &lt; 2^14 are stored in two bytes.
// <exception cref="InvalidDataException">Invalid binary format.</exception>
func ReadVarUint(reader *bufio.Reader) uint {
	var num uint = 0
	var length = 0

	for true {
		r, err := reader.ReadByte()
		if err != nil {
			return 0
		}
		num |= (uint(r) & Bits7) << length
		length += 7

		if uint(r) < Bit8 {
			return num
		}

		if length > 35 {
			// throw new InvalidDataException("Integer out of range.")
		}
	}
	return 0
}

// ReadVarInt Reads a 32-bit variable length signed integer.
// 1/8th of storage is used as encoding overhead.
// * Values &lt; 2^7 are stored in one byte.
// * Values &lt; 2^14 are stored in two bytes.
// <exception cref="InvalidDataException">Invalid binary format.</exception>
func ReadVarInt(reader *bufio.Reader) (uint, uint, error) {
	byteReader := bufio.NewReader(reader)
	r, err := byteReader.ReadByte()
	if err != nil {
		return 0, 0, err
	}
	var num = uint(r) & Bits6
	var length = 6
	var sign uint = 1
	if uint(r)&Bit7 > 0 {
		sign = -1
	}

	if (uint(r) & Bit8) == 0 {
		// Don't continue reading.
		return sign * num, sign, nil
	}
	for true {
		r, err := byteReader.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		num |= (uint(r) & Bits7) << length
		length += 7

		if uint(r) < Bit8 {
			return sign * num, sign, nil
		}

		if length > 41 {
			// throw new InvalidDataException("Integer out of range")
			return 0, 0, errors.New("Integer out of range")
		}
	}
	return 0, 0, err
}

// ReadAny Decodes data from the reader.
func ReadAny(reader *bufio.Reader) any {
	Type, err := reader.ReadByte()
	if err != nil {
		return nil
	}
	switch uint(Type) {
	case 119: // String
		return ReadVarString(reader)
	case 120: // boolean true
		return true
	case 121: // boolean false
		return false
	case 123: // Float64
		var dBytes = make([]byte, 8)
		l, err := reader.Read(dBytes)
		if err != nil || l != 8 {
			return nil
		}
		if !IsLittleEndian() {
			dBytes = Reverse(dBytes)
		}
		bits := binary.LittleEndian.Uint64(dBytes)
		return math.Float64frombits(bits)
	case 124: // Float32
		var dBytes = make([]byte, 4)
		l, err := reader.Read(dBytes)
		if err != nil || l != 4 {
			return nil
		}
		if !IsLittleEndian() {
			dBytes = Reverse(dBytes)
		}
		bits := binary.LittleEndian.Uint32(dBytes)
		return math.Float32frombits(bits)
	case 125: // integer
		_, b, _ := ReadVarInt(reader)
		return b
	case 126: // null
	case 127: // undefined
		return nil
	case 116: // ArrayBuffer
		return ReadVarUint8Array(reader)
	case 117: // Array<any>
		var length = ReadVarUint(reader)
		var arr = make([]any, length)
		for i := 0; uint(i) < length; i++ {
			arr = append(arr, ReadAny(reader))
		}
		return arr
	case 118: // any (Dictionary<string, any>)
		var length = ReadVarUint(reader)
		var obj = make(map[string]any, length)
		for i := 0; i < int(length); i++ {
			var key = ReadVarString(reader)
			obj[key] = ReadAny(reader)
		}
		return obj
	default:
		// throw new InvalidDataException($"Unknown any type: {type}")
	}
	return nil
}

// Reads a byte from the reader and advances the position within the reader by one byte.
// <exception cref="EndOfreaderException">End of reader reached.</exception>
func ReadByte(reader *bufio.Reader) byte {
	v, err := reader.ReadByte()
	if err != nil {
		return v
	}
	return v
}

// ReadBytes Reads a sequence of bytes from the current reader and advances the position
// within the reader by the number of bytes read.
// <exception cref="EndOfreaderException">End of reader reached.</exception>
func ReadBytes(reader *bufio.Reader, buffer []byte) error {
	if len(buffer) == 0 {
		return errors.New("ReadBytes error")
	}
	readLength, err := reader.Read(buffer)
	if len(buffer) != readLength || err != nil {
		// throw new EndOfreaderException()
		return errors.New("ReadBytes error")
	}
	return nil
}

func IsLittleEndian() bool {
	var value int32 = 1 // 占4byte 转换成16进制 0x00 00 00 01
	// 大端(16进制)：00 00 00 01
	// 小端(16进制)：01 00 00 00
	pointer := unsafe.Pointer(&value)
	pb := (*byte)(pointer)
	if *pb != 1 {
		return false
	}
	return true
}

func Reverse(bytes []byte) []byte {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

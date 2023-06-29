package protocols

import (
	"bufio"
	"encoding/binary"
	"io"

	"YJS-GO/utils"
)

type SyncStep uint

const (
	YjsSyncStep1 SyncStep = 0
	YjsSyncStep2 SyncStep = 1
	YjsUpdate    SyncStep = 2
)

type SyncProtocol struct {
}

func (*SyncProtocol) WriteSyncStep1(stream io.Writer, doc utils.YDoc) {
	binary.Write(stream, binary.LittleEndian, uint64(YjsSyncStep1))

	var sv = doc.EncodeStateVectorV2()
	binary.Write(stream, binary.LittleEndian, sv)
}

func (s *SyncProtocol) ReadSyncStep1(reader io.Reader, writer io.Writer, doc utils.YDoc) {
	bufioReader := bufio.NewReader(reader)
	length, err := binary.ReadUvarint(bufioReader)
	if err != nil {
		return
	}
	encodedStateVector := make([]byte, length)
	binary.Read(bufioReader, binary.LittleEndian, encodedStateVector)
	s.WriteSyncStep2(writer, doc, encodedStateVector)
}

func (*SyncProtocol) WriteSyncStep2(stream io.Writer, doc utils.YDoc, encodedStateVector []byte) {

	binary.Write(stream, binary.LittleEndian, uint64(YjsSyncStep2))

	var update = doc.EncodeStateAsUpdateV2(encodedStateVector)
	binary.Write(stream, binary.LittleEndian, update)

}

func (*SyncProtocol) ReadSyncStep2(reader io.Reader, doc utils.YDoc, transactionOrigin interface{}) {
	bufioReader := bufio.NewReader(reader)
	length, err := binary.ReadUvarint(bufioReader)
	if err != nil {
		return
	}
	encodedStateVector := make([]byte, length)
	binary.Read(bufioReader, binary.LittleEndian, encodedStateVector)

	doc.ApplyUpdateV2(encodedStateVector, transactionOrigin)
}

func (*SyncProtocol) WriteUpdate(stream io.Writer, update []byte) {
	binary.Write(stream, binary.LittleEndian, uint64(YjsUpdate))
	binary.Write(stream, binary.LittleEndian, update)
}

func (s *SyncProtocol) ReadUpdate(stream io.Reader, doc utils.YDoc, transactionOrigin interface{}) {
	s.ReadSyncStep2(stream, doc, transactionOrigin)
}

func (s *SyncProtocol) ReadSyncMessage(reader io.Reader, writer io.Writer, doc utils.YDoc,
	transactionOrigin interface{}) uint {
	messageType, err := binary.ReadUvarint(bufio.NewReader(reader))
	if err != nil {
		return 0
	}
	switch SyncStep(messageType) {
	case YjsSyncStep1:
		s.ReadSyncStep1(reader, writer, doc)
		break
	case YjsSyncStep2:
		s.ReadSyncStep2(reader, doc, transactionOrigin)
		break
	case YjsUpdate:
		s.ReadUpdate(reader, doc, transactionOrigin)
		break
	default:
		return 0
		// Exception($"Unknown message type: {messageType}")
	}

	return uint(messageType)
}

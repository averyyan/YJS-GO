package protocols

import (
	"encoding/binary"
	"io"
)

type SyncStep uint

const (
	YjsSyncStep1 SyncStep = 0
	YjsSyncStep2 SyncStep = 1
	YjsUpdate    SyncStep = 2
)

type SyncProtocol struct {
}

func (*SyncProtocol) WriteSyncStep1(stream io.ByteReader, YDoc doc) {
	stream.WriteVarUint(YjsSyncStep1)

	var sv = doc.EncodeStateVectorV2()
	stream.WriteVarUint8Array(sv)
}

func (*SyncProtocol) WriteSyncStep2(stream io.ByteReader, YDoc doc, byte []encodedStateVector) {
	stream.WriteVarUint(YjsSyncStep2)

	var update = doc.EncodeStateAsUpdateV2(encodedStateVector)
	stream.WriteVarUint8Array(update)
}

func (*SyncProtocol) ReadSyncStep1(reader io.ByteReader, writer io.ByteWriter, doc YDoc) {
	var encodedStateVector = reader.ReadVarUint8Array()
	WriteSyncStep2(writer, doc, encodedStateVector)
}

func (*SyncProtocol) ReadSyncStep2(stream io.ByteReader, YDoc doc, object transactionOrigin) {
	var update = stream.ReadVarUint8Array()
	doc.ApplyUpdateV2(update, transactionOrigin)
}

func (*SyncProtocol) WriteUpdate(stream io.ByteWriter, byte []update) {
	stream.WriteVarUint(YjsUpdate)
	stream.WriteVarUint8Array(update)
}

func (*SyncProtocol) ReadUpdate(stream io.ByteReader, YDoc doc, object transactionOrigin) {
	ReadSyncStep2(stream, doc, transactionOrigin)
}

func (s *SyncProtocol) ReadSyncMessage(reader io.ByteReader, Stream writer, YDoc doc, object transactionOrigin) uint {
	messageType, err := binary.ReadUvarint(reader)
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

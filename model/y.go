package model

type AbstractConnector interface {
}

type AbstractPersistence interface {
}

type Y[K string, V any] struct {
	um          *UndoManager
	connected   bool
	connector   AbstractConnector
	ds          string
	gcEnabled   bool
	os          string
	persistence AbstractPersistence
	room        string
	share       map[K]V
	ss          string
	userID      string
}

type YEvent struct {
}

func createRelativePositionFromTypeIndex(ytext string, ty int) {

}

func createAbsolutePositionFromRelativePosition() {

}

func (*Y[K, V]) writeDeleteSet(e encoder) {

}
func (*Y[K, V]) toBinary() encoder {
	return nil
}

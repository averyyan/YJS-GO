package model

type connector struct {
	authInfo string

	broadcastBuffer string

	broadcastBufferSize int

	broadcastBufferSizePos string

	checkAuth string

	connections string

	currentSyncTarget string

	debug string

	isSynced bool

	log string

	logMessage string

	maxBufferLength string

	opts string

	protocolVersion int

	role string

	userEventListeners []string

	whenSyncedListeners []string

	y string
}

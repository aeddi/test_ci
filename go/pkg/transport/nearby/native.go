package nearby

type NativeDriver interface {
	Start(localPeerID string) error
	Stop() error
	Dial(remotePeerID string) bool
}

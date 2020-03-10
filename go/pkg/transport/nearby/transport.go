package nearby

import (
	"context"
	"fmt"

	nearby_ma "github.com/ipfs-shipyard/gomobile-ipfs/go/pkg/transport/nearby/multiaddr"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	ma "github.com/multiformats/go-multiaddr"

	logging "github.com/ipfs/go-log"
	"github.com/pkg/errors"
)

const DefaultBind = "/nearby/Qmeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"

var log = logging.Logger("nearby_transport")

// Transport is a Nearby tpt.transport.
var _ tpt.Transport = &Transport{}

// Transport represents any device by which you can connect to and accept
// connections from other peers.
type Transport struct {
	host     host.Host
	upgrader *tptu.Upgrader
}

// NewTransport creates a Nearby transport object that tracks dialers and listener.
// It also starts the discovery service.
func NewTransport(h host.Host, u *tptu.Upgrader) (*Transport, error) {
	return &Transport{
		host:     h,
		upgrader: u,
	}, nil
}

// Dial dials the peer at the remote address.
// With Nearby you can only dial a device that is already connected with the native driver.
func (t *Transport) Dial(ctx context.Context, remoteMa ma.Multiaddr, remotePID peer.ID) (tpt.CapableConn, error) {
	// Nearby transport needs to have a running listener in order to dial other peer
	// because native driver is initialized during listener creation.
	if gListener == nil {
		return nil, errors.New("transport dialing peer failed: no active listener")
	}

	// remoteAddr is supposed to be equal to remotePID since with Nearby transport:
	// multiaddr == /nearby/<peerID>
	remoteAddr, err := remoteMa.ValueForProtocol(nearby_ma.P_NEARBY)
	if err != nil || remoteAddr != remotePID.Pretty() {
		return nil, errors.Wrap(err, "transport dialing peer failed: wrong multiaddr")
	}

	// Check if native driver is already connected to peer's device.
	// With Nearby you can't really dial, only auto-connect with peer nearby.
	// if bledrv.DialPeer(remoteAddr) == false {
	// 	return nil, errors.New("transport dialing peer failed: peer not connected through Nearby")
	// }

	// Can't have two connections on the same multiaddr
	if _, ok := connMap.Load(remoteAddr); ok {
		return nil, errors.New("transport dialing peer failed: already connected to this address")
	}

	// Returns an outbound conn.
	return newConn(ctx, t, remoteMa, remotePID, false)
}

// CanDial returns true if this transport believes it can dial the given
// multiaddr.
func (t *Transport) CanDial(remoteMa ma.Multiaddr) bool {
	return nearby_ma.NEARBY.Matches(remoteMa)
}

// Listen listens on the given multiaddr.
// Nearby can't listen on more than one listener.
func (t *Transport) Listen(localMa ma.Multiaddr) (tpt.Listener, error) {
	// localAddr is supposed to be equal to localPID or to DefaultBind since with
	// Nearby transport: multiaddr == /nearby/<peerID>
	localPID := t.host.ID().Pretty()
	localAddr, err := localMa.ValueForProtocol(nearby_ma.P_NEARBY)
	if err != nil || (localMa.String() != DefaultBind && localAddr != localPID) {
		return nil, errors.Wrap(err, "transport listen failed: wrong multiaddr")
	}

	// Replaces default bind by local host peerID.
	if localMa.String() == DefaultBind {
		localMa, err = ma.NewMultiaddr(fmt.Sprintf("/nearby/%s", localPID))
		if err != nil { // Should never occur.
			panic(err)
		}
	}

	// If a global listener already exists, returns an error.
	if gListener != nil {
		// TODO: restore this when published as generic lib / fixed in Berty network
		// config update
		// return nil, errors.New("transport listen failed: one listener maximum")
		gListener.Close()
	}

	return newListener(localMa, t)
}

// Proxy returns true if this transport proxies.
func (t *Transport) Proxy() bool {
	return false
}

// Protocols returns the set of protocols handled by this transport.
func (t *Transport) Protocols() []int {
	return []int{nearby_ma.P_NEARBY}
}

func (t *Transport) String() string {
	return "Nearby"
}

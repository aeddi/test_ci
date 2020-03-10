package multiaddr

import (
	ma "github.com/multiformats/go-multiaddr"
)

// Add Nearby to the list of libp2p's multiaddr protocols.
func init() {
	err := ma.AddProtocol(protoNearby)
	if err != nil {
		panic(err) // Should never occur.
	}
}

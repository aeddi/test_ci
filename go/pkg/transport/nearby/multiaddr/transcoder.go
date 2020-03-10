package multiaddr

import (
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

// Nearby multiaddr transcoder
// See https://github.com/multiformats/go-multiaddr/blob/master/transcoders.go
var TranscoderNearby = ma.NewTranscoderFromFunctions(nearbyStB, nearbyBtS, nearbyVal)

func nearbyStB(s string) ([]byte, error) {
	_, err := peer.IDB58Decode(s)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}

func nearbyBtS(b []byte) (string, error) {
	_, err := peer.IDB58Decode(string(b))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func nearbyVal(b []byte) error {
	_, err := peer.IDB58Decode(string(b))
	return err
}

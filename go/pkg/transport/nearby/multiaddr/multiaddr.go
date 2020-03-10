package multiaddr

import (
	ma "github.com/multiformats/go-multiaddr"
)

// Nearby multiaddr protocol definition
// See https://github.com/multiformats/go-multiaddr/blob/master/protocols.go
// See https://github.com/multiformats/multiaddr/blob/master/protocols.csv
const P_NEARBY = 0x0042

var protoNearby = ma.Protocol{
	Name:       "nearby",
	Code:       P_NEARBY,
	VCode:      ma.CodeToVarint(P_NEARBY),
	Size:       -1,
	Path:       false,
	Transcoder: TranscoderNearby,
}

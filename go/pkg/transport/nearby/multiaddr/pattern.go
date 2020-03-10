package multiaddr

import (
	mafmt "github.com/multiformats/go-multiaddr-fmt"
)

// Nearby multiaddr validation checker
// See https://github.com/multiformats/go-multiaddr-fmt
var NEARBY = mafmt.Base(P_NEARBY)

package state

import "github.com/coniks-sys/coniks-go/crypto/vrf"

type User struct {
	sk vrf.PrivateKey
	pk vrf.PublicKey

	weight uint64

	pp ProbParams

	sortitionIntervals []ProbInterval
}
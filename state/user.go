package state

import (
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"log"
)

// TODO: factor sortition-specific attributes into sub-struct?
type User struct {
	sk vrf.PrivateKey
	pk vrf.PublicKey

	weight uint64

	pp ProbParams

	sortitionIntervals []ProbInterval
}

func MakeTestUser(weight uint64, privKeyBytes []byte) *User {
	var sk vrf.PrivateKey
	if privKeyBytes == nil {
		sk, _ = vrf.GenerateKey(nil) // TODO: rando?
	} else {
		sk = vrf.PrivateKey(privKeyBytes)
	}
	return &User{sk, emptyKey, weight, ppInit, make([]ProbInterval, 0)}
}


func (u *User) isPubKeyDefined() bool {
	return len([]byte(u.pk)) != 0
}

func (u *User) initPublicFromPrivateKey() {
	pk, ok := u.sk.Public()
	if !ok {
		log.Panic("Public key init failed")
	}
	u.pk = pk
}

func (u *User) GetPublicKeyBytes() []byte {
	return u.pk
}

func (u *User) Sign(msgs ...[]byte) []byte {
	return make([]byte, 0) // TODO !!!
}

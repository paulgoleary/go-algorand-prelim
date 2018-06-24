package crypto

import (
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"golang.org/x/crypto/ed25519"
	"crypto"
	"log"
)

func Sign(sk vrf.PrivateKey, message []byte) []byte {
	sig, err := ed25519.PrivateKey(sk).Sign(nil, message, crypto.Hash(0))
	if err != nil {
		log.Panic("Unable to sign message")
	}
	return sig
}

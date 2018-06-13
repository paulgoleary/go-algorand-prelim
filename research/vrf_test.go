package research

import (
	"bytes"
	"testing"

	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"log"
)

// cribbed from coniks-go vrf_test
func TestHonestComplete(t *testing.T) {
	sk, err := vrf.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	pk, _ := sk.Public()
	alice := []byte("alice")
	aliceVRF := sk.Compute(alice)
	log.Printf("vrf: %v", aliceVRF)
	aliceVRFFromProof, aliceProof := sk.Prove(alice)

	// fmt.Printf("pk:           %X\n", pk)
	// fmt.Printf("sk:           %X\n", *sk)
	// fmt.Printf("alice(bytes): %X\n", alice)
	// fmt.Printf("aliceVRF:     %X\n", aliceVRF)
	// fmt.Printf("aliceProof:   %X\n", aliceProof)

	if !pk.Verify(alice, aliceVRF, aliceProof) {
		t.Error("Gen -> Compute -> Prove -> Verify -> FALSE")
	}
	if !bytes.Equal(aliceVRF, aliceVRFFromProof) {
		t.Error("Compute != Prove")
	}
}
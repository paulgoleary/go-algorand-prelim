package crypto

import (
	"testing"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"golang.org/x/crypto/ed25519"
	"reflect"
	goCrypto "crypto"
)

func TestCryptoCompat(t *testing.T) {

	privKeyBytes := []byte {
		114, 48, 107, 12, 26, 4, 242, 59, 80, 179, 244, 231, 237, 138, 203, 76,
		231, 118, 0, 87, 31, 67, 89, 47, 122, 37, 216, 236, 48, 137, 81, 35,
		217, 39, 229, 179, 68, 179, 86, 114, 231, 184, 80, 232, 95, 78, 4, 73,
		195, 194, 47, 181, 191, 128, 41, 43, 76, 136, 36, 170, 123, 115, 59, 99 }

	skVRF := vrf.PrivateKey(privKeyBytes)
	pkVRF, _ := skVRF.Public()

	skGoX := ed25519.PrivateKey(privKeyBytes)
	pkGoX := skGoX.Public()

	bytesVRF := []byte(pkVRF)
	bytesGoX := []byte(pkGoX.(ed25519.PublicKey)) // YOINKS!
	if !reflect.DeepEqual(bytesVRF, bytesGoX) {
		t.Errorf("Expected public keys to be binary compatible")
	}

	someBytes := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	sigGoX, _ := skGoX.Sign(nil, someBytes, goCrypto.Hash(0))

	sigThunk := Sign(skVRF, someBytes)

	if !reflect.DeepEqual(sigGoX, sigThunk) {
		t.Errorf("Expected signatures to be binary compatible")
	}

}

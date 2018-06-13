package crypto

import (
	"testing"
)

func TestSortition(t *testing.T) {

	tau := uint64(1000)
	totalWeights := uint64(1000 * 1000)

	privKeyBytes := []byte {
		114, 48, 107, 12, 26, 4, 242, 59, 80, 179, 244, 231, 237, 138, 203, 76,
		231, 118, 0, 87, 31, 67, 89, 47, 122, 37, 216, 236, 48, 137, 81, 35,
		217, 39, 229, 179, 68, 179, 86, 114, 231, 184, 80, 232, 95, 78, 4, 73,
		195, 194, 47, 181, 191, 128, 41, 43, 76, 136, 36, 170, 123, 115, 59, 99 }

	user := MakeTestUser(1000, privKeyBytes)

	role := "person"
	seed := []byte {0xDE, 0xAD, 0xBE, 0xEF}

	if user.checkSortitionPrecalc(tau, totalWeights) {
		t.Error("User should be init'ed with no interval pre-calcs")
	}

	j := user.Sortition(role, seed, tau, totalWeights)
	if j != 2 {
		t.Error("Given fixed user (sk), role, seed, weights and tau chosen j should be stable.")
	}

	if !user.checkSortitionPrecalc(tau, totalWeights) {
		t.Error("For fixed weights and tau User should pre-calc and save intervals.")
	}

	tau = uint64(10 * 1000)

	if user.checkSortitionPrecalc(tau, totalWeights) {
		t.Error("Changing params to Sortition should invalidate pre-calc.")
	}

	j = user.Sortition(role, seed, tau, totalWeights)
	if j != 13 {
		t.Error("Given fixed user (sk), role, seed, weights and tau chosen j should be stable.")
	}

	user.SetWeight(100)

	if user.checkSortitionPrecalc(tau, totalWeights) {
		t.Error("Changing user weight should invalidate pre-calc.")
	}

}
package state

import (
	"testing"
	"math/big"
	"fmt"
	"math/rand"
	"log"
)

func kindaCompare( li ProbInterval, ri ProbInterval ) bool {
	lstr := fmt.Sprintf("%.4f %.4f", li.start, li.end)
	rstr := fmt.Sprintf("%.4f %.4f", ri.start, ri.end)
	return lstr == rstr
}

func kindaCompIntervals(leftIntervals []ProbInterval, rightIntervals []ProbInterval) bool {
	if len(leftIntervals) != len(rightIntervals) {
		return false
	}

	for i, li := range leftIntervals {
		ri := rightIntervals[i]
		if !kindaCompare(li, ri) {
			return false
		}
	}

	return true
}

func TestBinoms(t *testing.T) {

	// binomials from 10 coin tosses
	expectBinoms := []float64 {0.000977, 0.009766, 0.043945, 0.117188, 0.205078, 0.246094, 0.205078, 0.117188, 0.043945, 0.009766, 0.000977}
	expectIntervals := make([]ProbInterval, len(expectBinoms))

	for i := 0; i < len(expectIntervals); i++ {
		if i == 0 {
			expectIntervals[i].start = big.NewFloat(0.0)
			expectIntervals[i].end = big.NewFloat(expectBinoms[i])
		} else {
			expectIntervals[i].start = expectIntervals[i - 1].end
			expectIntervals[i].end = big.NewFloat(expectBinoms[i])
			expectIntervals[i].end.Add(expectIntervals[i].end, expectIntervals[i].start)
		}
	}

	// coin toss experiment
	tau := uint64(50)
	totalWeights := uint64(100)

	role := "person"
	seed := []byte {0xDE, 0xAD, 0xBE, 0xEF}

	user := MakeTestUser(10, nil)

	user.Sortition(role, seed, tau, totalWeights)

	if !kindaCompIntervals(expectIntervals[0:11], user.sortitionIntervals) {
		t.Error("Expected (kinda) equivalent results")
	}
}

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

	_, _, j := user.Sortition(role, seed, tau, totalWeights)
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

	_, _, j = user.Sortition(role, seed, tau, totalWeights)
	if j != 13 {
		t.Error("Given fixed user (sk), role, seed, weights and tau chosen j should be stable.")
	}

	user.SetWeight(100)

	if user.checkSortitionPrecalc(tau, totalWeights) {
		t.Error("Changing user weight should invalidate pre-calc.")
	}
}

// maybe want more powerful test framework?
// see discussion here: https://stackoverflow.com/questions/31595791/how-to-test-panics
func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestSortAndVerify(t *testing.T) {

	tau := uint64(50)
	totalWeights := uint64(100)

	role := "person"
	seed := []byte {0xDE, 0xAD, 0xBE, 0xEF}

	user := MakeTestUser(10, nil)

	hashBytes, proofBytes, j := user.Sortition(role, seed, tau, totalWeights)

	// should panic because test user does not yet have a public key
	assertPanic(t, func() {user.VerifySort(role, seed, tau, totalWeights, hashBytes, proofBytes)})

	user.initPublicFromPrivateKey()

	testJ, _ := user.VerifySort(role, seed, tau, totalWeights, hashBytes, proofBytes)
	if j != testJ {
		t.Errorf("Inconsistent value for j: expect %v, got %v", j, testJ)
	}

	_, err := user.VerifySort(role, seed, tau, totalWeights, append(hashBytes, 0), proofBytes)
	if err == nil {
		t.Error("Verify should fail if hash is different")
	}

	_, err = user.VerifySort(role, seed, tau, totalWeights, hashBytes, append(proofBytes, 0))
	if err == nil {
		t.Error("Verify should fail if proof is different")
	}

	_, err = user.VerifySort(role, append(seed, 0), tau, totalWeights, hashBytes, proofBytes)
	if err == nil {
		t.Error("Verify should fail if seed is different")
	}
}

func randoSortAndVerify(t *testing.T, weightHisto map[int]int) {

	tau := uint64(50)
	totalWeights := uint64(1000 * 1000 * 1000) // 1B token supply

	role := "person"
	seed := make([]byte, seedSize)
	rand.Read(seed) // rando seed

	// user has 1m 'staked' tokens
	user := MakeTestUser(1000 * 1000, nil) // rando user - i.e. private key

	hashBytes, proofBytes, j := user.Sortition(role, seed, tau, totalWeights)

	user.initPublicFromPrivateKey()

	testJ, _ := user.VerifySort(role, seed, tau, totalWeights, hashBytes, proofBytes)

	if j != testJ {
		t.Errorf("Inconsistent value for j: expect %v, got %v", j, testJ)
	}
	weightHisto[j] += 1
}

func TestRandoSortAndVerify(t *testing.T) {

	weightHistoMap := make(map[int]int)

	for i := 0; i < 1000; i++ {
		randoSortAndVerify(t, weightHistoMap)
	}

	maxJ := -1
	for j, _ := range weightHistoMap {
		if j > maxJ {
			maxJ = j
		}
	}

	weightHisto := make([]int, maxJ + 1)

	for j, cnt := range weightHistoMap {
		weightHisto[j] = cnt
	}

	log.Printf("test j (user count) histo: %v", weightHisto)
}
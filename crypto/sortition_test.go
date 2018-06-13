package crypto

import (
	"testing"
	"log"
	"math/rand"
	"time"
)

func TestSortition(t *testing.T) {

	tau := int64(1000)
	totalWeight := int64(1000 * 1000)
	user := MakeTestUser(1000)

	role := "person"
	seed := user.sk.Compute([]byte(role))

	j := user.Sortition(role, seed, tau, totalWeight)
	log.Printf("got j: %v", j)
}

func randoSelect(tau int64, totalWeight int64, userWeight int64) int {
	cntSelect := 0
	for i := int64(0); i < userWeight; i++ {
		if rand.Int63n(totalWeight) < tau {
			cntSelect += 1
		}
	}
	return cntSelect
}

func histoRandoSelect(numTrials int, tau int64, totalWeight int64, userWeight int64) []int {
	rand.Seed(time.Now().UnixNano())
	selectCnts := make([]int, userWeight)
	for i := 0; i < numTrials; i++ {
		s := randoSelect(tau, totalWeight, userWeight)
		selectCnts[s] += 1
	}
	return selectCnts
}

func TestMath(t *testing.T) {

	test10 := f64Facto(10)
	if test10 != 3628800.0 {
		t.Errorf( "Wrong value for 10 factorial: %v", test10)
	}
	test10Again := f64Facto(10)
	if test10 != test10Again {
		t.Errorf("Value for 10 should not have changed")
	}

	test := f64BinomDist(10, int64(100), 0.5)
	// TODO: verify answer?
	if test != 1.3655426387463099e-17 {
		t.Errorf("Wrong value for test: %v", test)
	}

	w := int64(100)
	totalWeight := int64(1000 * 1000)
	tau := int64(30)

	selectHisto := histoRandoSelect(100 * 1000, tau, totalWeight, w)
	for i, s := range selectHisto {
		if s != 0 {
			log.Printf("histo bucket, cnt: %v, %v", i, s)
		}
	}

	log.Printf("")

	prob := float64(tau) / float64(totalWeight)
	log.Printf("prob: %v", prob)

	sumBinomialDist := func(j int64) float64 {
		res := 0.0
		for k := int64(0); k < j; k++ {
			res += f64BinomDist(k, w, prob)
		}
		return res
	}

	for j := int64(0); j < w; j++ {
		xStart := sumBinomialDist(j)
		xEnd := sumBinomialDist(j + 1)
		if xStart != xEnd {
			log.Printf( "interval %v: %v to %v, %v", j, xStart, xEnd, xEnd - xStart)
		}
	}
}

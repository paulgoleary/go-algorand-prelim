package crypto

import (
	"testing"
	"log"
	"math/rand"
	"time"
	"fmt"
)

func TestSortition(t *testing.T) {
	tau := int64(1000)
	totalWeight := int64(1000 * 1000)
	user := MakeTestUser(1000)

	for i := 0; i < 10; i++ {
		role := fmt.Sprintf("person-%v", i)
		seed := user.sk.Compute([]byte(role))
		j := user.Sortition(role, seed, tau, totalWeight)
		log.Printf("got j: %v", j)
	}
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
	w := int64(100)
	totalWeight := int64(1000 * 1000)
	tau := int64(30)

	selectHisto := histoRandoSelect(100 * 1000, tau, totalWeight, w)
	for i, s := range selectHisto {
		if s != 0 {
			log.Printf("histo bucket, cnt: %v, %v", i, s)
		}
	}
}

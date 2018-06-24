package research

import (
	"testing"
	"math/big"
	"log"
	"github.com/paulgoleary/go-algorand/state"
)

/*
re-write p^k * (1-p)^(w-k) where p = (t/Wt) as t^k * (Wt - t)^(w-k) / Wt^w
 */
func bigScaledBinomDist(k , w, t, Wt uint64, scaleFactor *big.Float) *big.Int {

	binomCoeff := big.NewInt(0).Binomial(int64(w), int64(k))

	// t^k
	nl := big.NewInt(int64(t))
	nl = nl.Exp(nl, big.NewInt(int64(k)), nil)

	// OOOOOPH - these exponents are slow bigly with large integers!
	// (Wt - t)^(w-k)
	nr := big.NewInt(int64(Wt - t))
	nr = nr.Exp(nr, big.NewInt(int64(w - k)), nil)

	// Wt^w
	d := big.NewInt(int64(Wt))
	d = d.Exp(d, big.NewInt(int64(w)), nil)

	n := nl.Mul(nl, nr)
	n.Mul(n, binomCoeff)

	scaled := &big.Float{}
	scaled.SetInt(n)
	scaled.Mul(scaled, scaleFactor)
	scaled.Quo(scaled, big.NewFloat(0.0).SetInt(d))

	scaledInt, _ := scaled.Int(nil)

	return scaledInt
}

func TestBigBinomials(t *testing.T) {

	tau := uint64(1000)
	Wt := uint64(1000 * 1000)

	userWeight := uint64(1000)

	user := state.MakeTestUser(userWeight, nil)
	user.Sortition("", make([]byte,0), tau, Wt)
	log.Printf("test user: %v", user)

	// TODO: same as crypto
	hashDenomInt := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil)
	hashDenom := big.NewFloat(0.0).SetInt(hashDenomInt)

	scale := hashDenom.Quo(hashDenom, big.NewFloat(float64(Wt)))
	log.Printf("%v", scale)

	bigZero := big.NewInt(0)
	checkTotal := big.NewInt(0)

	checkProbLoss := big.NewInt(0)

	bigBinoms := make([]*big.Int, userWeight)
	for i := uint64(0); i < userWeight; i++ {
		bigBinoms[i] = bigScaledBinomDist(i, userWeight, tau, Wt, scale)
		if bigBinoms[i].Cmp(bigZero) == 0 {
			log.Printf("Got binom of zero - quiting")
			break
		}
		log.Printf("binom %v: %v", i, bigBinoms[i])
		checkTotal.Add(checkTotal, bigBinoms[i])
		if i >= 18 {
			checkProbLoss.Add(checkProbLoss, bigBinoms[i])
		}
	}
	checkPercent := big.NewFloat(0.0).SetInt(checkTotal)
	checkPercent.Quo(checkPercent, hashDenom)
	checkPercentLoss := big.NewFloat(0.0).SetInt(checkProbLoss)
	checkPercentLoss.Quo(checkPercentLoss, hashDenom)
	log.Printf("check total percent: %v", checkPercent)
	log.Printf("check prob loss: %v", checkPercentLoss)
}
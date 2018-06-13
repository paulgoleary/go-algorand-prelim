package crypto

import (
	"math/big"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	// "log"
	"math"
	"sort"
)

// algorand paper says they use curve 25519 and sha-256 for hash function

type ProbInterval struct {
	piStart *big.Float
	piEnd *big.Float
}

type ProbParams struct {
	tau int64
	totalWeights int64
}

type User struct {
	sk vrf.PrivateKey

	weight int64

	pp ProbParams

	sortitionIntervals []ProbInterval
}

func MakeTestUser(weight int64) *User {
	sk, _ := vrf.GenerateKey(nil) // TODO
	return &User{sk, weight, ProbParams{-1, -1}, make([]ProbInterval, 0)}
}

// in algorand the seed for a round is based on a vrf from the previous round
var seedSize = vrf.Size
var bigHashLen = big.NewInt(vrf.Size * 8)

// pre-compute 2^hashLen
var hashDenomInt = big.NewInt(0).Exp(big.NewInt(2), bigHashLen, nil)
var hashDenom = big.NewFloat(0.0).SetInt(hashDenomInt)

func binomDist(k int64, w int64, prob float64 ) *big.Float {

	binomCoeff := big.NewFloat(0.0).SetInt(big.NewInt(0).Binomial(w, k))

	// go currently has no big.Float exponentiation :/
	pExp := big.NewFloat(0.0).SetFloat64(math.Pow(prob, float64(k)))
	npExp := big.NewFloat(0.0).SetFloat64(math.Pow(1.0 - prob, float64(w - k)))

	return binomCoeff.Mul(binomCoeff, pExp.Mul(pExp, npExp))
}

func (u *User) getSortitionIntervals(tau, totalWeights int64) []ProbInterval {
	if tau == u.pp.tau && totalWeights == u.pp.totalWeights {
		return u.sortitionIntervals
	}
	// ok that this is float64 since golang does not support exponentiation with big.Float
	prob := float64(tau) / float64(totalWeights)

	sumBinomialDist := func(j int64) *big.Float {
		res := big.NewFloat(0.0)
		for k := int64(0); k < j; k++ {
			res = res.Add(res, binomDist(k, u.weight, prob))
		}
		return res
	}

	intervals := make([]ProbInterval, 0)
	for j := int64(0); j < u.weight; j++ {
		xStart := sumBinomialDist(j)
		xEnd := sumBinomialDist(j + 1)
		if xStart.Cmp(xEnd) == 0 {
			break
		}
		intervals = append(intervals, ProbInterval{xStart, xEnd})
	}
	u.pp.tau = tau
	u.pp.totalWeights = totalWeights
	u.sortitionIntervals = intervals
	return intervals
}

func (u *User) Sortition(role string, seed []byte, tau, totalWeights int64) int {
	intervals := u.getSortitionIntervals(tau, totalWeights)

	// <hash, pi> <- VRFsk(seed||role)
	roleBytes := []byte(role)
	hash := u.sk.Compute(append(seed, roleBytes...))
	hashInt := big.NewInt(0).SetBytes(hash)

	randProb := new(big.Float).SetInt(hashInt).Quo(big.NewFloat(0.0).SetInt(hashInt), hashDenom)

	cmpInterval := func(i int) bool {
		xx := intervals[i]
		cmp := xx.piEnd.Cmp(randProb)
		return cmp >= 0
	}
	j := sort.Search(len(intervals), cmpInterval)
	return j
}
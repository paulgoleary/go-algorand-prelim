package crypto

import (
	"math/big"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	// "log"
	"math"
	"sort"
)

var facts map[int64]*big.Int

func bigFactoMemo(n int64) *big.Int {
	if facts[n] != nil {
		res := facts[n]
		return res
	}

	if n > 0 {
		x := bigFactoMemo(n - 1)
		res := big.NewInt(n)
		res = res.Mul(res, x)
		facts[n] = res
		return res
	}

	return big.NewInt(1)
}

var f64Facts = map[int64]float64{}

func f64FactoMemo(n int64) float64 {
	if f64Facts[n] != 0.0 {
		res := f64Facts[n]
		return res
	}

	if n > 0 {
		x := f64FactoMemo(n - 1)
		res := x * float64(n)
		f64Facts[n] = res
		return res
	}

	return 1.0
}

func f64Facto(n int64) float64 {
	return f64FactoMemo(n)
}

// algorand paper says they use curve 25519 and sha-256 for hash function

type ProbInterval struct {
	piStart float64
	piEnd float64
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
var hashDenom = big.NewInt(0).Exp(big.NewInt(2), bigHashLen, nil)
var bfHashDemom = big.NewFloat(0.0).SetInt(hashDenom)

// TODO: return big.Float?
func f64BinomDist(k int64, w int64, prob float64 ) *big.Float {
	/*
	wFacto := f64Facto(w)
	kFacto := f64Facto(k)
	wkFacto := f64Facto(w - k)
	binomCoeff := wFacto / (kFacto * wkFacto)
	*/

	bfBinomCoeff := big.NewFloat(0.0).SetInt(big.NewInt(0).Binomial(w, k))

	// go currently has no big.Float exponentiation :/
	pExp := big.NewFloat(0.0).SetFloat64(math.Pow(prob, float64(k)))
	npExp := big.NewFloat(0.0).SetFloat64(math.Pow(1.0 - prob, float64(w - k)))

	res, _ := bfBinomCoeff.Mul(bfBinomCoeff, pExp.Mul(pExp, npExp)).Float64()
	return res
}

func (u *User) getSortitionIntervals(tau, totalWeights int64) []ProbInterval {
	if tau == u.pp.tau && totalWeights == u.pp.totalWeights {
		return u.sortitionIntervals
	}
	prob := float64(tau) / float64(totalWeights)

	sumBinomialDist := func(j int64) float64 {
		res := 0.0
		for k := int64(0); k < j; k++ {
			res += f64BinomDist(k, u.weight, prob)
		}
		return res
	}

	intervals := make([]ProbInterval, 0)
	for j := int64(0); j < u.weight; j++ {
		xStart := sumBinomialDist(j)
		xEnd := sumBinomialDist(j + 1)
		if xStart == xEnd {
			break
		}
		intervals = append(intervals, ProbInterval{xStart, xEnd})
	}
	return intervals
}

func (u *User) Sortition(role string, seed []byte, tau, totalWeights int64) int {

	intervals := u.getSortitionIntervals(tau, totalWeights)

	// <hash, pi> <- VRFsk(seed||role)
	roleBytes := []byte(role)
	hash := u.sk.Compute(append(seed, roleBytes...))
	hashInt := big.NewInt(0).SetBytes(hash)

	bfRandProb := new(big.Float).SetInt(hashInt).Quo(big.NewFloat(0.0).SetInt(hashInt), bfHashDemom)
	randProb, _ := bfRandProb.Float64()

	cmpInterval := func(i int) bool {
		xx := intervals[i]
		return xx.piEnd >= randProb
	}
	j := sort.Search(len(intervals), cmpInterval)
	return j
}
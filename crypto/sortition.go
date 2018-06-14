package crypto

import (
	"math/big"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"math"
	"sort"
	"fmt"
)

// algorand paper says they use curve 25519 and sha-256 for hash function

type ProbInterval struct {
	start *big.Float
	end *big.Float
}

func (pi ProbInterval) String() string {
	return fmt.Sprintf("[%v, %v)", pi.start, pi.end)
}

type ProbParams struct {
	tau uint64
	totalWeights uint64
}

type User struct {
	sk vrf.PrivateKey

	weight uint64

	pp ProbParams

	sortitionIntervals []ProbInterval
}

var ppInit = ProbParams{^uint64(0), ^uint64(0)}

func MakeTestUser(weight uint64, privKeyBytes []byte) *User {
	var sk vrf.PrivateKey
	if privKeyBytes == nil {
		sk, _ = vrf.GenerateKey(nil) // TODO: rando?
	} else {
		sk = vrf.PrivateKey(privKeyBytes)
	}
	return &User{sk, weight, ppInit, make([]ProbInterval, 0)}
}

// in algorand the seed for a round is based on a vrf from the previous round
var seedSize = vrf.Size
var bigHashLen = big.NewInt(vrf.Size * 8)

// pre-compute 2^hashLen
var hashDenomInt = big.NewInt(0).Exp(big.NewInt(2), bigHashLen, nil)
var hashDenom = big.NewFloat(0.0).SetInt(hashDenomInt)

func binomDist(k , w uint64, prob float64 ) *big.Float {

	binomCoeff := big.NewFloat(0.0).SetInt(big.NewInt(0).Binomial(int64(w), int64(k)))

	// go currently has no big.Float exponentiation :/
	pExp := big.NewFloat(0.0).SetFloat64(math.Pow(prob, float64(k)))
	npExp := big.NewFloat(0.0).SetFloat64(math.Pow(1.0 - prob, float64(w - k)))

	return binomCoeff.Mul(binomCoeff, pExp.Mul(pExp, npExp))
}

func (u *User) checkSortitionPrecalc(tau, totalWeights uint64) bool {
	if tau == u.pp.tau && totalWeights == u.pp.totalWeights {
		return true
	} else {
		return false
	}
}

func (u *User) getSortitionIntervals(tau, totalWeights uint64) []ProbInterval {
	if u.checkSortitionPrecalc(tau, totalWeights) {
		return u.sortitionIntervals
	}
	// ok that this is float64 since golang does not support exponentiation with big.Float
	prob := float64(tau) / float64(totalWeights)

	sumBinomialDist := func(j uint64) *big.Float {
		res := big.NewFloat(0.0)
		for k := uint64(0); k < j; k++ {
			res = res.Add(res, binomDist(k, u.weight, prob))
		}
		return res
	}

	intervals := make([]ProbInterval, 0)
	for j := uint64(0); j < u.weight + 1; j++ {
		start := sumBinomialDist(j)
		end := sumBinomialDist(j + 1)
		if start.Cmp(end) == 0 {
			break
		}
		intervals = append(intervals, ProbInterval{start, end})
	}
	u.pp.tau = tau
	u.pp.totalWeights = totalWeights
	u.sortitionIntervals = intervals
	return intervals
}

func (u *User) getHashInt(role string, seed []byte) *big.Int {
	// <hash, pi> <- VRFsk(seed||role)
	roleBytes := []byte(role)
	hash := u.sk.Compute(append(seed, roleBytes...))
	return big.NewInt(0).SetBytes(hash)
}

func (u *User) Sortition(role string, seed []byte, tau, totalWeights uint64) int {
	intervals := u.getSortitionIntervals(tau, totalWeights)

	hashInt := u.getHashInt(role, seed)

	randProb := new(big.Float).SetInt(hashInt).Quo(big.NewFloat(0.0).SetInt(hashInt), hashDenom)

	cmpInterval := func(i int) bool {
		xx := intervals[i]
		cmp := xx.end.Cmp(randProb)
		return cmp >= 0
	}
	j := sort.Search(len(intervals), cmpInterval)
	return j
}

func (u *User) SetWeight(newWeight uint64) {
	u.weight = newWeight
	u.pp = ppInit
}

func (u *User) String() string {
	return fmt.Sprintf("%v", u.sortitionIntervals) // TODO: more !!!
}
package state

import (
	"math/big"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"math"
	"sort"
	"fmt"
	"log"
	"github.com/pkg/errors"
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

var ppInit = ProbParams{^uint64(0), ^uint64(0)}
var emptyKey = make([]byte, 0)

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

	// TODO: maybe could be more optimal. could calc all the binomials first and then calc the intervals
	// since we calc and re-use this is likely not a significant issue but i suppose that depends on how
	//  effective the re-use is; e.g. if total weights change frequently between calculations
	// on the other hand, this method quits fairly early when the weight is high because the probabilities converge
	//  quickly and we seem to run out of precision ...
	intervals := make([]ProbInterval, 0)
	for j := uint64(0); j < u.weight + 1; j++ {
		var start *big.Float
		if j == 0 {
			start = big.NewFloat(0.0)
		} else {
			start = intervals[j - 1].end
		}
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

func makeMessage(role string, seed []byte) []byte {
	roleBytes := []byte(role)
	return append(seed, roleBytes...)
}

func (u *User) calcSubUsers(tau, totalWeights uint64, hashBytes []byte) int {
	intervals := u.getSortitionIntervals(tau, totalWeights)

	hashInt := big.NewInt(0).SetBytes(hashBytes)

	randProb := new(big.Float).SetInt(hashInt).Quo(big.NewFloat(0.0).SetInt(hashInt), hashDenom)

	cmpInterval := func(i int) bool {
		xx := intervals[i]
		cmp := xx.end.Cmp(randProb)
		return cmp >= 0
	}
	j := sort.Search(len(intervals), cmpInterval)

	return j
}

func (u *User) Sortition(role string, seed []byte, tau, totalWeights uint64) ([]byte, []byte, int) {

	hashBytes, hashProof := u.sk.Prove(makeMessage(role, seed))

	j := u.calcSubUsers(tau, totalWeights, hashBytes)

	return hashBytes, hashProof, j
}

func (u *User) VerifySort(role string, seed []byte, tau, totalWeights uint64, hashBytes, proofBytes []byte) (int, error) {

	if !u.isPubKeyDefined() {
		log.Panic("Public key must be defined on User for this operation.")
	}

	if !u.pk.Verify(makeMessage(role, seed), hashBytes, proofBytes) {
		return -1, errors.New("Failed VRF verification")
	}

	j := u.calcSubUsers(tau, totalWeights, hashBytes)

	return j, nil
}

func (u *User) SetWeight(newWeight uint64) {
	u.weight = newWeight
	u.pp = ppInit
}

func (u *User) String() string {
	return fmt.Sprintf("%v", u.sortitionIntervals) // TODO: more !!!
}
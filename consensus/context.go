package consensus

import (
	"fmt"
	"github.com/paulgoleary/go-algorand/state"
	"math/big"
)

type Context struct {
	theUser *state.User
	currentSeed []byte
}

func (c *Context) getTotalWeights() uint64 {
	return 0 // TODO !!!
}

func (c *Context) getLastBlockHash() []byte {
	return make([]byte, 0) // TODO!
}

func makeRoleString(roleName string, round, step int) string {
	return fmt.Sprintf("%v|%v|%v", roleName, round, step)
}

func (c *Context) CommitteeVote(round, step int, tau uint64, value []byte) {
	role := makeRoleString("committee", round, step)

	hashBytes, hashProof, j := c.theUser.Sortition(role, c.currentSeed, tau, c.getTotalWeights())

	if j > 0 {
		// i won !!!
		signBytes := [6][]byte {
			big.NewInt(int64(round)).Bytes(),
			big.NewInt(int64(step)).Bytes(),
			hashBytes,
			hashProof,
			c.getLastBlockHash(),
			value,
		}
		c.Gossip(c.theUser.GetPublicKeyBytes(), c.theUser.Sign(signBytes[:]...))
	}
}

func (c *Context) Gossip(userPubKeyBytes []byte, committeeSig []byte) {
	return // TODO
}
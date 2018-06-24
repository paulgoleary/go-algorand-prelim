package consensus

import (
	"fmt"
	"github.com/paulgoleary/go-algorand/state"
)

type Context struct {
	theUser *state.User
	currentSeed []byte
}

func (c *Context) getTotalWeights() uint64 {
	return 0 // TODO !!!
}

func makeRoleString(roleName string, round, step int) string {
	return fmt.Sprintf("%v|%v|%v", roleName, round, step)
}

func (c *Context) CommitteeVote(round, step int, tau uint64, value []byte) {
	role := makeRoleString("committee", round, step)

	hashBytes, hashProof, j := c.theUser.Sortition(role, c.currentSeed, tau, c.getTotalWeights())

	if j > 0 {
		// i won !!!
		c.Gossip(c.theUser.GetPublicKeyBytes(), )
	}
}
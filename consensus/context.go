package consensus

import "fmt"

type Context struct {

}

func makeRoleString(roleName string, round, step int) string {
	return fmt.Sprintf("%v|%v|%v", roleName, round, step)
}

func (c *Context) CommitteeVote(round, step int, tau uint64, value []byte) {
	role := makeRoleString("committee", round, step)


}
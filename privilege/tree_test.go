package privilege

import (
	"fmt"
	"testing"
)

func TestInsertChild(t *testing.T) {
	n := &node{}
	n.addPrivilege("/admin/:name", true)
	n.addPrivilege("/admin/:name/1", true)

	// 打印所有节点信息
gg:
	for {
		fmt.Println(n)
		for _, c := range n.children {
			n = c
			continue gg
		}
		break
	}
}

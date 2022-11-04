package privilege

import (
	"fmt"
	"testing"
)

func TestInsertChild(t *testing.T) {
	n := &node{}
	n.addPrivilege("/admin/add/1")
	n.addPrivilege("/admin/1")

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

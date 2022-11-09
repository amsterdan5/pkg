package privilege

import (
	"sync"
)

var (
	// 用户权限列表
	userPrivilegeList = make(map[uint]*userPrivilege)
	// 锁
	syncP sync.Mutex
)

type userPrivilege struct {
	privilege methodTrees
}

// 新建权限
func NewPrivilege(uid uint, authList map[string][]string) *userPrivilege {
	if p, ok := userPrivilegeList[uid]; ok {
		return p
	}

	return refreshUserPrivilege(uid, authList)
}

// 刷新权限树
func refreshUserPrivilege(uid uint, authList map[string][]string) *userPrivilege {
	syncP.Lock()

	p := &userPrivilege{
		privilege: make(methodTrees, 0, 9),
	}

	for m, uris := range authList {
		if !checkMethod(m) {
			panic("存在无效方法")
		}

		for _, path := range uris {
			p.AddPrivilege(m, path)
		}
	}

	userPrivilegeList[uid] = p
	syncP.Unlock()
	return p
}

// 添加权限
func (p *userPrivilege) AddPrivilege(method, path string) {
	root := p.privilege.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		p.privilege = append(p.privilege, methodTree{method: method, tree: root})
	}

	root.addPrivilege(path, true)
}

// 返回节点
func (p *userPrivilege) GetNode(method string) *node {
	return p.privilege.get(method)
}

// 鉴权方法
func (p *userPrivilege) CheckPrivilege(method, path string) bool {
	if pn := p.privilege.get(method); pn != nil {
		return pn.checkPrivilege(path)
	}
	return false
}

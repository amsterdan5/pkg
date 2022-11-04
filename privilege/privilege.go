package privilege

type userPrivilege struct {
	Uid       int
	privilege methodTrees
}

// 新建权限
func NewPrivilege(uid int) *userPrivilege {
	return &userPrivilege{
		Uid:       uid,
		privilege: make(methodTrees, 0, 9),
	}
}

// 检查权限
func (p *userPrivilege) AddPrivilege(method, path string) {
	root := p.privilege.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		p.privilege = append(p.privilege, methodTree{method: method, tree: root})
	}

	root.addPrivilege(path)
}

package privilege

import (
	"strings"
)

type methodTree struct {
	method string
	tree   *node
}

type methodTrees []methodTree

// 返回节点
func (m methodTrees) get(method string) *node {
	method = strings.ToLower(method)

	for _, t := range m {
		if t.method == method {
			return t.tree
		}
	}
	return nil
}

// 节点类型
type nodeType uint8

const (
	static nodeType = iota
	root
	param
	catchAll
)

// 节点信息
type node struct {
	path      string   // 节点对应字符串路径
	nType     nodeType // 节点类型
	wildChild bool     // 是否通配符
	maxParams uint8    // 最大字符数
	priority  uint32   // 节点层级
	fullPath  string   // 原始路径
	indices   string   // 节点与子节点分裂的第一个字符
	children  []*node  // 子节点
	valid     bool     // 是否有效路由
}

// 添加权限
// 参数：
//   path 路由
func (n *node) addPrivilege(path string, valid bool) {
	fullPath := path
	n.priority++
	numParams := countParams(path)

	// 创建根节点
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(numParams, path, fullPath, valid)
		n.nType = root
		return
	}

	parentFullPathIndex := 0

insert:
	for {
		// 更新节点深度
		if numParams > n.maxParams {
			n.maxParams = numParams
		}

		// 查找公共前缀
		i := longesCommonPrefix(path, n.path)

		// 处理子节点内容
		if i < len(n.path) {
			child := node{
				path:      n.path[i:],
				wildChild: n.wildChild,
				indices:   n.indices,
				children:  n.children,
				priority:  n.priority - 1,
				fullPath:  n.fullPath,
			}

			// 更新所有子节点规则长度
			for _, v := range child.children {
				if v.maxParams > child.maxParams {
					child.maxParams = v.maxParams
				}
			}

			n.children = []*node{&child}
			n.indices = string([]byte{n.path[i]})
			n.path = path[:i]
			n.wildChild = false
			n.fullPath = fullPath[:parentFullPathIndex+i]
		}

		if i < len(path) {
			path = path[i:]

			if n.wildChild {
				parentFullPathIndex += len(n.path)
				n = n.children[0]
				n.priority++

				// 更新节点长度
				if numParams > n.maxParams {
					n.maxParams = numParams
				}
				numParams--

				// 查找通配符
				if len(path) >= len(n.path) && n.path == path[:len(n.path)] {
					// 查找长通配符，例如 :name、:names
					if len(n.path) >= len(path) || path[len(n.path)] == '/' {
						continue insert
					}
				}

				pathSeg := path
				// 不是全匹配
				if n.nType != catchAll {
					pathSeg = strings.SplitN(path, "/", 2)[0]
				}

				prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path
				panic("'" + fullPath + "'与已知路由'" + prefix + "'重叠")
			}

			c := path[0]

			// 缩减参数
			if n.nType == param && c == '/' && len(n.children) == 1 {
				parentFullPathIndex += len(n.path)
				n = n.children[0]
				n.priority++
				continue insert
			}

			// 检查节点时候还有其他参数
			for i, max := 0, len(n.indices); i < max; i++ {
				if c == n.indices[i] {
					parentFullPathIndex += len(n.path)
					i = n.incrementChildPrio(i)
					n = n.children[i]
					continue insert
				}
			}

			// 普通路由
			if c != ':' && c != '*' {
				n.indices += string([]byte{c})
				child := &node{
					maxParams: numParams,
					fullPath:  fullPath,
				}
				n.children = append(n.children, child)
				n.incrementChildPrio(len(n.indices) - 1)
				n = child
			}
			n.insertChild(numParams, path, fullPath, valid)
			return
		}
		n.valid = valid
		return
	}

}

// 插入节点
// 参数：
//   numParms 路由长度
//   path     当前路由
//   fullPath 原始路由
func (n *node) insertChild(numParams uint8, path, fullPath string, valid bool) {
	for numParams > 0 {
		wildcard, position, valid := findWildcard(path)
		if position < 0 {
			break
		}

		if !valid {
			panic("无效的路由规则")
		}

		if wildcard != "*" && len(wildcard) < 2 {
			panic("通配符必须2个字符以上")
		}

		if len(n.children) > 0 {
			panic("该节点已经有子节点")
		}

		// 通配符节点
		if wildcard[0] == ':' {
			if position > 0 {
				n.path = path[:position]
				path = path[position:]
			}

			n.wildChild = true
			child := &node{
				nType:     param,
				path:      wildcard,
				maxParams: numParams,
				fullPath:  fullPath,
			}

			n.children = []*node{child}
			n = child
			n.priority++
			numParams--

			// 通配符后还有路由规则
			if len(wildcard) < len(path) {
				path = path[len(wildcard):]

				child := &node{
					maxParams: numParams,
					priority:  1,
					fullPath:  fullPath,
				}
				n.children = []*node{child}
				n = child
				continue
			}
			n.valid = valid
			return
		}

		// 全匹配
		if position+len(wildcard) != len(path) || numParams > 1 {
			panic("* 只允许在路由最后")
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == '/' {
			panic("* 不能包含 /")
		}

		// 检测”/”符
		position--
		if path[position] != '/' {
			panic("* 前必须为 /")
		}

		n.path = path[:position]
		child := &node{
			wildChild: true,
			nType:     catchAll,
			maxParams: 1,
			fullPath:  fullPath,
		}

		// 重置父节点长度
		if n.maxParams < 1 {
			n.maxParams = 1
		}
		n.children = []*node{child}
		n.indices = string('/')
		n = child
		n.priority++

		// 包含参数的子节点
		child = &node{
			path:      path[position:],
			nType:     catchAll,
			maxParams: 1,
			valid:     valid,
			priority:  1,
			fullPath:  fullPath,
		}
		n.children = []*node{child}

		return
	}

	// 没有通配符，直接插入节点
	n.path = path
	n.fullPath = fullPath
	n.valid = valid
}

// 提升子节点层级
func (n *node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// 调整位置
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		// 交换节点位置
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	// 位置需调整，新建字符串
	if newPos != pos {
		n.indices = n.indices[:newPos] + // 不变的前缀
			n.indices[pos:pos+1] + // 移除首字符
			n.indices[newPos:pos] + n.indices[pos+1:] // 在pos位置开始调整
	}

	return newPos
}

// 计算路由长度
// 参数：
// 	 path 要解析的链接
func countParams(path string) uint8 {
	n := uint(0)
	for i := 0; i < len(path); i++ {
		if path[i] == ':' || path[i] == '*' {
			n++
		}
	}
	if n > 255 {
		n = 255
	}
	return uint8(n)
}

// 查找通配符位置
// 参数：
// 	 path 要解析的链接
// 返回：
//   通配符
//   通配符位置
//   是否通配符
func findWildcard(path string) (string, int, bool) {
	for start, name := range []byte(path) {
		// ':'为参数, '*'为匹配所有
		if name != ':' && name != '*' {
			continue
		}

		// 规则是否有效
		valid := true
		// 通配符截止为止
		for end, w := range []byte(path[start+1:]) {
			switch w {
			case '/':
				return path[start : start+1+end], start, valid
			case ':':
			case '*':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}

// 匹配路由规则
// 参数：
// 	 path 要解析的链接
func (n *node) checkPrivilege(path string) bool {
walk:
	for {
		prefix := n.path
		if path == prefix {
			// 当前路由有效
			if n.valid {
				return true
			}

			if path == "/" && n.wildChild && n.nType != root {
				return true
			}

			indices := n.indices
			indLen := len(indices)
			// 是否有公共前缀
			for i := 0; i < indLen; i++ {
				if indices[i] == '/' {
					n = n.children[i]
					if n.nType == catchAll {
						return true
					}
				}
			}

			return n.valid
		}

		// 循环查找
		if len(path) > len(prefix) && path[:len(prefix)] == prefix {
			path = path[len(prefix):]
			// 非通配符情况循环匹配
			if !n.wildChild {
				c := path[0]
				indices := n.indices
				for i, max := 0, len(indices); i < max; i++ {
					if c == indices[i] {
						n = n.children[i]
						continue walk
					}
				}

				// 查看是否 / 结尾
				return path == "/" && n.valid
			}

			// 查找节点
			n = n.children[0]
			switch n.nType {
			case param:
				// 找到下一个 / 以获取通配符长度
				end := 0
				for end < len(path) && path[end] != '/' {
					end++
				}

				// 循环查找
				if end < len(path) {
					if len(n.children) > 0 {
						path = path[end:]
						n = n.children[0]
						continue walk
					}

					return len(path) == end+1
				}

				if n.valid {
					return true
				}

				if len(n.children) == 1 {
					// 查看是否 / 结尾
					n = n.children[0]
					return n.path == "/" && n.valid
				}
				return false

			case catchAll:
				return true
			default:
				panic("无效节点类型")
			}
		}

		// 根目录 或 请求地址最后多个 /f
		return (path == "/") ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && n.valid)
	}
}

package privilege

import "strings"

// 返回最小值
func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// 查找公共前缀
func longesCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}

	return i
}

// 检测请求方法是否正确
func checkMethod(method string) bool {
	// 请求方式
	methods := []string{
		"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH",
	}

	method = strings.ToUpper(method)
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}

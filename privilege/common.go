package privilege

import "strings"

// 返回最小值
func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// 查找公共前缀
func LongesCommonPrefix(a, b string) int {
	i := 0
	max := Min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}

	return i
}

// 检测请求方法是否正确
func CheckMethod(method string) bool {
	// 请求方式
	methods := []string{
		"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH",
	}

	for _, m := range methods {
		if m == strings.ToLower(method) {
			return true
		}
	}
	return false
}

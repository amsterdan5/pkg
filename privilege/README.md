## restful 权限校验
> 基于 radix tree 基数树实现的路由鉴权

### 使用方式
```
p := NewPrivilege(1)

routes := []string{
	"/admin/:name",
	"/admin/:name/1",
}

apis := map[string][]string{
	"get": routes,
}

for method, api := range apis {
	if CheckMethod(method) {
		for _, a := range api {
			p.AddPrivilege(method, a)
		}
	}
}

fmt.Println(p.privilege.get("get").checkPrivilege("/admin"))
```
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

### 增加发布订阅方式
> 发布订阅用于异步通知更新

#### 发布端
```
c := &ManageConfig{
	Type: "redis",
	Addr: "127.0.0.1",
	Port: "6379",
}
ok, err := Publish(NewManage(c), key, "123")
if err != nil {
	fmt.Println(err)
}
fmt.Println(ok)
```

#### 订阅端
```
c := &ManageConfig{
	Type: "redis",
	Addr: "127.0.0.1",
	Port: "6379",
}
m := NewManage(c)
msg := Subscribe(m, key)
fmt.Println(msg, 1)

f := func() map[string][]string {
	data := map[string][]string{
		"get": []string{
			"/index",
		},
	}
	return data
}
fmt.Println(m.SetFunc(f).Do(1))
```
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

>> 使用redis
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
c := privilege.ManageConfig{
	Type: "redis",
	Addr: "127.0.0.1",
	Port: "6379",
}

msg := make(chan string)
go privilege.Subscribe(privilege.NewManage(c), key, msg)

for m := range msg {
	fmt.Println(m)
}

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
>> 使用rabbitmq
#### 发布端
```
c := privilege.ManageConfig{
	Type:     "rabbitmq",
	Addr:     "www.local.com",
	Port:     "5672",
	UserName: "myuser",
	Password: "mypass",
}
_, err := privilege.Publish(privilege.NewManage(c), "queue", "test mq")
if err != nil {
	fmt.Println(err)
}
```
#### 订阅端
```
c := privilege.ManageConfig{
	Type:     "rabbitmq",
	Addr:     "www.local.com",
	Port:     "5672",
	UserName: "myuser",
	Password: "mypass",
}

msg := make(chan string)
go privilege.Subscribe(privilege.NewManage(c), "queue", msg)
for {
	select {
	case mm := <-msg:
		fmt.Println(mm)
	}
}
```
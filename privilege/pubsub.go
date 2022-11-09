package privilege

import (
	"context"
	"errors"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	once sync.Once
	m    manage
)

// 发布订阅接口
type manage interface {
	// 发布
	publish(key, data string) (bool, error)
	// 订阅
	sublish(key string) chan string
	// 设置方法
	SetFunc(f ManageFunc) manage
	// 执行方法
	Do(uid uint) bool
}

// 发布消息
// 	p       manage  管理类型
// 	key     string  订阅字段
// 	data    string  发布内容
func Publish(p manage, key, data string) (bool, error) {
	return p.publish(key, data)
}

// 订阅消息
// 	p       manage  管理类型
// 	key     string  订阅字段
func Subscribe(p manage, key string) chan string {
	return p.sublish(key)
}

// 配置管理
type ManageConfig struct {
	Type     string
	Addr     string
	Port     string
	UserName string
	Password string
}

// 自定义方法
type ManageFunc func() map[string][]string

// redis订阅类
func (c *ManageConfig) newRedisManage() manage {
	once.Do(func() {
		m = redisManage{
			ctx: context.Background(),
			rdb: redis.NewClient(&redis.Options{
				Addr: c.Addr + ":" + c.Port,
			}),
		}
	})

	return m
}

// 订阅类
func NewManage(c *ManageConfig) (m manage) {
	switch c.Type {
	case "redis":
		m = c.newRedisManage()
	default:
		panic("未找到方法")
	}

	return
}

type redisManage struct {
	ctx context.Context
	rdb *redis.Client
	f   ManageFunc
}

// 发布消息
func (r redisManage) publish(key, data string) (bool, error) {
	res := r.rdb.Publish(r.ctx, key, data)
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

// 订阅消息
func (r redisManage) sublish(key string) chan string {
	message := <-r.rdb.Subscribe(r.ctx, key).Channel()

	m := make(chan string)
	m <- message.Payload
	return m
}

// 设置数据来源方法
func (r redisManage) SetFunc(f ManageFunc) manage {
	r.f = f
	return r
}

// 刷新权限
func (r redisManage) Do(uid uint) bool {
	if r.f == nil {
		panic(errors.New("缺少数据来源方法"))
	}

	refreshUserPrivilege(uid, r.f())
	return true
}

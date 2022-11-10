package privilege

import (
	"context"
	"errors"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	once sync.Once
	m    manageInterface
)

// 发布订阅接口
type manageInterface interface {
	// 发布
	publish(key, data string) (bool, error)
	// 订阅
	sublish(key string, s chan string)
	// 设置上下文
	SetCtx(c context.Context)
	// 设置方法
	SetFunc(f ManageFunc)
	// 执行方法
	Do(uid uint) bool
}

// 发布消息
// 	p       manage  管理类型
// 	key     string  订阅字段
// 	data    string  发布内容
func Publish(p manageInterface, key, data string) (bool, error) {
	return p.publish(key, data)
}

// 订阅消息
// 	p       manage  管理类型
// 	key     string  订阅字段
func Subscribe(p manageInterface, key string, recevier chan string) {
	p.sublish(key, recevier)
}

// 公共参数
type manage struct {
	f   ManageFunc
	ctx context.Context
}

// 设置上下文
func (m *manage) SetCtx(c context.Context) {
	m.ctx = c
}

// 设置方法
func (m *manage) SetFunc(f ManageFunc) {
	m.f = f
}

// 执行方法
func (m *manage) Do(uid uint) bool {
	if m.f == nil {
		panic(errors.New("缺少数据来源方法"))
	}

	refreshUserPrivilege(uid, m.f())
	return true
}

// 配置管理
type ManageConfig struct {
	Type     string
	Addr     string
	Port     string
	UserName string
	Password string
}

// 订阅类
func NewManage(c ManageConfig) (m manageInterface) {
	switch c.Type {
	case "redis":
		m = c.newRedisManage()
	case "rabbitmq":
		m = c.newRabbitmqManage()
	default:
		panic("未找到方法")
	}

	return
}

// 自定义方法
type ManageFunc func() map[string][]string

// redis订阅类
func (c ManageConfig) newRedisManage() manageInterface {
	once.Do(func() {
		m = &redisManage{
			rdb: redis.NewClient(&redis.Options{
				Addr: c.Addr + ":" + c.Port,
			}),
		}

		m.SetCtx(context.Background())
	})

	return m
}

// mq订阅类
func (c ManageConfig) newRabbitmqManage() manageInterface {
	once.Do(func() {
		m = newRabbitmqManage(c)
	})

	return m
}

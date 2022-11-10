package privilege

import (
	"github.com/go-redis/redis/v8"
)

type redisManage struct {
	manage
	rdb *redis.Client
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
func (r redisManage) sublish(key string, recevier chan string) {
	b := make(chan bool)
	go func() {
		for message := range r.rdb.Subscribe(r.ctx, key).Channel() {
			recevier <- message.Payload
		}
	}()
	<-b
}

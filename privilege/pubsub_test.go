package privilege

import (
	"fmt"
	"testing"
	"time"
)

var (
	key = "test_pub_sub"
)

func TestStart(t *testing.T) {
	go testSubscribe()

	time.Sleep(time.Duration(3))
	testPublish()

	time.Sleep(time.Duration(5))
}

func testPublish() {
	c := ManageConfig{
		Type: "redis",
		Addr: "127.0.0.1",
		Port: "6379",
	}
	_, err := Publish(NewManage(c), key, "123")
	if err != nil {
		fmt.Println(err)
	}
}

func testSubscribe() {
	c := ManageConfig{
		Type: "redis",
		Addr: "127.0.0.1",
		Port: "6379",
	}

	msg := make(chan string)
	Subscribe(NewManage(c), key, msg)
	fmt.Println(<-msg)
}

func TestFunc(t *testing.T) {

	c := ManageConfig{
		Type: "redis",
		Addr: "127.0.0.1",
		Port: "6379",
	}
	m := NewManage(c)

	f := func() map[string][]string {
		data := map[string][]string{
			"get": []string{
				"/index",
			},
		}
		return data
	}
	m.SetFunc(f)
	fmt.Println(m.Do(1))
}

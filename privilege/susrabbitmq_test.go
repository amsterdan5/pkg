package privilege

import (
	"fmt"
	"testing"
	"time"
)

func TestMq(t *testing.T) {

	go testMqSublish()

	time.Sleep(time.Duration(3))
	testMqPublish()

	time.Sleep(time.Duration(5))
}

func testMqPublish() {
	c := ManageConfig{
		Type:     "rabbitmq",
		Addr:     "www.local.com",
		Port:     "5672",
		UserName: "myuser",
		Password: "mypass",
	}
	_, err := Publish(newRabbitmqManage(c), "queue", "test mq")
	if err != nil {
		fmt.Println(err)
	}
}

func testMqSublish() {
	c := ManageConfig{
		Type:     "rabbitmq",
		Addr:     "www.local.com",
		Port:     "5672",
		UserName: "myuser",
		Password: "mypass",
	}
	mm := make(chan string)

	Subscribe(newRabbitmqManage(c), "queue", mm)
	fmt.Println(<-mm)
}

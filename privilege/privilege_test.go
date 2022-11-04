package privilege

import (
	"fmt"
	"testing"
)

func TestNewPrivilege(t *testing.T) {
	p := NewPrivilege(1)

	routes := []string{
		"/admin/:name/add",
		"/admin/:name",
	}

	apis := map[string][]string{
		"get": routes,
	}

	for method, api := range apis {
		if checkMethod(method) {
			for _, a := range api {
				p.AddPrivilege(method, a)
			}
		}
	}

	fmt.Println(p.CheckPrivilege("get", "/admin/1"))

}

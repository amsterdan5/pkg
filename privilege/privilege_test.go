package privilege

import (
	"fmt"
	"testing"
)

func TestNewPrivilege(t *testing.T) {
	p := NewPrivilege(1)

	routes := []string{
		"/admin/:name",
		"/admin/:name/add",
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

	fmt.Println(p.GetNode("get").checkPrivilege("/admin"))

}

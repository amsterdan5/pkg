package privilege

import (
	"fmt"
	"testing"
)

func TestNewPrivilege(t *testing.T) {
	routes := []string{
		"/admin/:name/add",
		"/admin/:name",
	}

	apis := map[string][]string{
		"get": routes,
	}

	p := NewPrivilege(1, apis)

	fmt.Println(p.CheckPrivilege("get", "/admin/1"))

}

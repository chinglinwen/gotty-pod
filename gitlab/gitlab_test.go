package gitlab

import (
	"fmt"
	"testing"
)

func TestGitlabVerify(t *testing.T) {
	p, err := GetProjects("wenzhenglin", "360860aA1")
	if err != nil {
		t.Error("verify err ", err)
		return
	}
	fmt.Println(p)
}

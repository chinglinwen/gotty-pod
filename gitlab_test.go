package main

import (
	"fmt"
	"testing"
)

func init() {
	// fmt.Println("set test init")
	// *GitlabAccessToken = "cAiBkXYjcPckzfJxcPnK" //robot
}

func TestGetProject(t *testing.T) {
	p, err := GetGitProject("yunwei/worktile")
	// p, err := GetGitProject("flow_center/agent_system")
	if err != nil {
		t.Error("get project err", err)
		return
	}
	fmt.Println("got project", p)
}

func TestCheckPerm(t *testing.T) {
	envs, err := CheckPerm("flow_center/audit-api", 75)
	if err != nil {
		t.Error("check perm err", err)
	}
	fmt.Println("got envs", envs)
}

func TestGetBinds(t *testing.T) {
	src := "/data/fluentd"
	envs := []string{"pre-online", "test", "online"}
	path := "/data/fluentd/flow-center/df-openapi"
	fmt.Println(getBinds(src, path, envs))
}

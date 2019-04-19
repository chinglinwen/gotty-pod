package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func init() {
	// fmt.Println("set test init")
	// *GitlabAccessToken = "cAiBkXYjcPckzfJxcPnK" //robot
}

func TestGetProject(t *testing.T) {
	// p, err := GetGitProject("xindaiquan/main")
	p, err := GetGitProject("flow_center/agent-system") //this need not agent_system
	// p, err := GetGitProject("flow_center/grade-users")

	if err != nil {
		t.Error("get project err", err)
		return
	}
	fmt.Println("got project", p)
}

func TestGetAccessLevel(t *testing.T) {
	// p, err := GetGitProject("xindaiquan/main")
	p, err := GetGitProject("flow_center/agent-system")
	// p, err := GetGitProject("flow_center/grade-users")
	if err != nil {
		t.Error("get project err", err)
		return
	}

	u, _ := GetUser(UserToken)
	fmt.Printf("got project: %v\nuser id: %v\n", p, u.ID)
	a, err := getAccessLevel(p, u.ID)
	if err != nil {
		t.Error("getAccessLevel err", err)
		return
	}
	fmt.Println("getAccessLevel", a)

}

func TestGetUser(t *testing.T) {
	u, err := GetUser(UserToken)
	if err != nil {
		t.Error("verify err ", err)
		return
	}
	// spew.Dump("user", u)
	b, err := json.Marshal(u)
	fmt.Println(string(b))
}

// need upgrade gitlab
// func TestUserInGroup(t *testing.T) {
// 	_, gs, _ := GetGroups(UserToken)
// 	u, _ := GetUser(UserToken)

// 	for _, v := range gs {
// 		fmt.Println("for", v.Path, u.Name, userIsInGroup(v, u.Name))
// 	}
// }

func TestGetProjects(t *testing.T) {
	_, pss, err := GetProjects(UserToken)
	if err != nil {
		t.Error("get projects err ", err)
		return
	}
	fmt.Println("got", len(pss))

	for _, v := range pss {
		if strings.Contains(v.WebURL, "trx") {
			fmt.Println("got", v.WebURL)
		}
	}
}

// func TestGetProjectsOld(t *testing.T) {
// 	_, pss, err := GetProjectsOld(UserToken)
// 	if err != nil {
// 		t.Error("get projects err ", err)
// 		return
// 	}
// 	fmt.Println("got", len(pss))
// }

func TestGetProjectLists(t *testing.T) {
	ps, err := GetProjectLists(UserToken, "/data/fluentd")
	fmt.Printf("got %v projects, err: %v\n", len(ps), err)
	for _, v := range ps {
		fmt.Println("got", v)
	}

	// data := fmt.Sprintf("%v", ps)

	// a := compress(data)
	// fmt.Println("a:", len(a))

	// b := uncompress(a)
	// fmt.Println("b:", len(b), b)
}

func TestGetGroups(t *testing.T) {
	// _, p, err := GetGroups(UserToken)
	_, p, err := GetGroups(UserToken)
	if err != nil {
		t.Error("get groups err ", err)
		return
	}
	for _, v := range p {
		fmt.Println(v.WebURL)
	}
	// spew.Dump(p)
}

func TestCheckPerm(t *testing.T) {
	// envs, err := CheckPerm("xindaiquan/base-service", UserToken)
	envs, err := CheckPerm("yunwei/worktile", UserToken)
	if err != nil {
		t.Error("check perm err", err)
		return
	}
	fmt.Println("got envs", envs)
}

func TestGetBinds(t *testing.T) {
	src := "/data/fluentd"
	envs := []string{"pre", "test", "online"}
	path := "/data/fluentd/flow-center/df-openapi"
	fmt.Println(getBinds(src, path, envs))
}

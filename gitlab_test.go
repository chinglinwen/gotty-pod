package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func init() {
	// fmt.Println("set test init")
	// *GitlabAccessToken = "cAiBkXYjcPckzfJxcPnK" //robot
}

func TestGetProject(t *testing.T) {
	p, err := GetGitProject("xindaiquan/main")
	// p, err := GetGitProject("flow_center/agent_system")
	if err != nil {
		t.Error("get project err", err)
		return
	}
	fmt.Println("got project", p)
}

func TestGetUser(t *testing.T) {
	u, err := GetUser("sHJm7wrnsZbnVtxNFsye")
	if err != nil {
		t.Error("verify err ", err)
		return
	}
	// spew.Dump("user", u)
	b, err := json.Marshal(u)
	fmt.Println(string(b))
}

func TestGetProjects(t *testing.T) {
	_, pss, err := GetProjects("sHJm7wrnsZbnVtxNFsye")
	if err != nil {
		t.Error("get projects err ", err)
		return
	}
	fmt.Println("got", len(pss))
}

func TestGetProjectsOld(t *testing.T) {
	_, pss, err := GetProjectsOld("sHJm7wrnsZbnVtxNFsye")
	if err != nil {
		t.Error("get projects err ", err)
		return
	}
	fmt.Println("got", len(pss))
}

func TestGetProjectLists(t *testing.T) {
	ps, err := GetProjectLists("sHJm7wrnsZbnVtxNFsye", "/data/fluentd")
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
	// _, p, err := GetGroups("sHJm7wrnsZbnVtxNFsye")
	_, p, err := GetGroups("iVLDbsPst5VCZjFdTFQo")
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
	envs, err := CheckPerm("xindaiquan/base-service", "sHJm7wrnsZbnVtxNFsye")
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

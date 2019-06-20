package main

import (
	"fmt"
	"testing"
)

func TestVerifyPermission(t *testing.T) {
	err := VerifyPermission()
	if err != nil {
		t.Error(err)
		return
	}
}
func TestValidateJWT(t *testing.T) {
	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoid2VuIiwidG9rZW4iOiJoZWxsbyIsImV4cCI6MTU2MTA5MTQwOX0.wWNoipQxwSUIf3HBgxaFVmFwJU2UuASjD_iQk59JPYM"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoid2VuIiwidG9rZW4iOiJoZWxsbyIsImV4cCI6MTU2MTA5NTQzOX0.CSQLZh2WB1Rn--1jAKGKIv5HbsuRhEImfL9-tCUOTZA"
	err := validateJWT(token)
	if err != nil {
		t.Error(err)
		return
	}
}

// func TestFilter(t *testing.T) {

// 	_, gitlist, err := GetProjectLists(UserToken)
// 	fmt.Printf("got %v projects, err: %v\n", len(gitlist), err)

// 	list, e := Walk("/data/fluentd")
// 	if e != nil {
// 		err = fmt.Errorf("walk error %v", e)
// 		return
// 	}

// 	sort.Slice(list, func(i, j int) bool { return list[i] < list[j] })
// 	sort.Slice(gitlist, func(i, j int) bool { return gitlist[i] < gitlist[j] })

// 	for _, v := range list {
// 		fmt.Printf("list %v\n", v)
// 	}
// 	for _, v := range gitlist {
// 		fmt.Printf("gitlist %v\n", v)
// 	}

// 	loglist := Filter(list, gitlist)
// 	fmt.Println("loglist", loglist)
// }

// go test -v -run TestGetProjectFromInput
func TestGetProjectFromInput(t *testing.T) {
	admin, ps, err := GetProjectLists(UserToken)
	fmt.Printf("got %v projects, err: %v\n", len(ps), err)

	//git := []string{"flow-center/df-openapi", "yunwei/trx", "yunwei/trx1"}
	// srcDir := "/data/fluentd"
	p, err := GetProjectFromInput(ps, admin)
	fmt.Println(p, err)
}

func TestSearchInput(t *testing.T) {
	inputlist := []string{"0 quit", "1 flow-center/df-openapi", "2 yunwei/trx", "3 yunwei/trx1"}

	i := searchInput("trx", inputlist)
	if i != 2 {
		t.Error("search input trx err, expect 2 got", i)
	}

	i = searchInput("quit", inputlist)
	if i != 0 {
		t.Error("search input quit err, expect 0 got", i)
	}
}

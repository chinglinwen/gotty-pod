package main

import (
	"fmt"
	"testing"
)

func TestGetProject(t *testing.T) {
	git := "flow-center/df-openapi yunwei/trx yunwei/trx1"
	srcDir := "/data/fluentd"
	p, err := GetProject(git, srcDir)
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

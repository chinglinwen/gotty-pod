package main

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	//k8s.Pod{Name:"tangguo-online", PodName:"tangguo-online-597fc44cb4-q8lsc", Env:"online", GitName:"flow-center/tangguo", Namespace:"flow-center"}
	out, err := run("flow-center", "tangguo-online-597fc44cb4-q8lsc")
	if err != nil {
		t.Error("run err", err)
		return
	}
	fmt.Println(out)
}

func TestRun2(t *testing.T) {
	//k8s.Pod{Name:"tangguo-online", PodName:"tangguo-online-597fc44cb4-q8lsc", Env:"online", GitName:"flow-center/tangguo", Namespace:"flow-center"}
	out, err := run2("flow-center", "tangguo-online-597fc44cb4-q8lsc")
	if err != nil {
		t.Error("run err", err)
		return
	}
	fmt.Println(out)
}

func TestRunPTY(t *testing.T) {
	//k8s.Pod{Name:"tangguo-online", PodName:"tangguo-online-597fc44cb4-q8lsc", Env:"online", GitName:"flow-center/tangguo", Namespace:"flow-center"}
	out, err := runpty("flow-center", "tangguo-online-597fc44cb4-q8lsc")
	if err != nil {
		t.Error("run err", err)
		return
	}
	fmt.Println(out)
}

// func TestRunStdin(t *testing.T) {
// 	//k8s.Pod{Name:"tangguo-online", PodName:"tangguo-online-597fc44cb4-q8lsc", Env:"online", GitName:"flow-center/tangguo", Namespace:"flow-center"}
// 	out, err := runstdin("flow-center", "tangguo-online-597fc44cb4-q8lsc")
// 	if err != nil {
// 		t.Error("run err", err)
// 		return
// 	}
// 	fmt.Println(out)
// }

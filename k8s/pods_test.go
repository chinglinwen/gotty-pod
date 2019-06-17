package k8s

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestPodList(t *testing.T) {
	ss, err := PodItems()
	if err != nil {
		t.Error("secretlist err", err)
		return
	}
	b, _ := json.MarshalIndent(ss, "", "  ")
	fmt.Println(string(b))
}

func TestgetNameEnv(t *testing.T) {
	cases := []struct {
		podname string
		want    Pod
	}{
		{"codis-ha-flow-center-8-yun-d659f8bd6-66mhs", Pod{Name: "codis-ha-flow-center-8-yun"}},
		{"adm-old-online-59fc977b48-mmv8k", Pod{Name: "adm-old", Env: "online"}},
		{"xpartner-pre-7c745b86dc-cthbn", Pod{Name: "xpartner", Env: "pre"}},
	}
	for _, v := range cases {
		name, env := getNameEnv(v.podname)
		if name != v.want.Name {
			t.Errorf("err got %v, want %v", name, v.want.Name)
			return
		}
		if env != v.want.Env {
			t.Errorf("err got %v, want %v", env, v.want.Env)
			return
		}
	}
}

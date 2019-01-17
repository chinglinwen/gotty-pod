package container

import (
	"fmt"
	"testing"
)

// normal user can't run this
func TestRun(t *testing.T) {
	c := &Container{
		Arg:        []string{"bash"},
		Src:        "/home/wen/log",
		Rootfs:     "/home/wen/alpine",
		Dst:        "/home/wen/targetrootfs",
		CGroupName: "wen",
		Hostname:   "container",
	}
	t.Error(c.Run())
}
func TestMount(t *testing.T) {
	// logdir should be one more folder, it flat with the / path
	out, err := Mount("/home/wen/log", "/home/wen/alpine", "/home/wen/targetrootfs")
	if err != nil {
		t.Errorf("mount err: %v", err)
	}
	fmt.Println("mount ok: ", string(out))
}

func TestUnMount(t *testing.T) {
	out, err := UnMount("/home/wen/targetrootfs")
	if err != nil {
		t.Errorf("umount err: %v", err)
	}
	fmt.Println("umount ok: ", string(out))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"./container"
)

var (
	srcDir = flag.String("src", "/data/fluentd", "log base directory on the host")
	dstDir = flag.String("dstDir", "/tmp", "destination directory on the host")
	rootFS = flag.String("rootfs", "/data/alpine", "rootfs on the host( eg. alpine rootfs")

	runChild = flag.Bool("child", false, "run container")
)

// the container can't be run directly, there need an parent
func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"-child"}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	log.Fatal(cmd.Run())
}

// link doesn't mount files, can't set workDir
// workDir can set in the log file structure only
func child(org, repo string) {
	// l := filepath.Join(*srcDir, org, repo)
	// t := filepath.Join("/tmp/", org, repo)
	// src := filepath.Join("/tmp/", org)
	// CreateLink(l, t)
	c := &container.Container{
		Arg:        []string{"sh"},
		Src:        filepath.Join(*srcDir, org, repo),
		Rootfs:     *rootFS,
		Dst:        filepath.Join(*dstDir, org, repo),
		CGroupName: repo,
		Hostname:   repo,
		WorkDir:    "/" + repo,
	}
	log.Fatal(c.Run())
	// RemoveLink(t)
}

func main() {
	flag.Parse()

	if len(os.Args) <= 2 {
		fmt.Println("error: not enough arguments")
		os.Exit(1)
	}

	var user, pass string
	git := "no-git-provide"
	for _, v := range os.Args[1:] {
		arg := strings.Split(v, "=")
		if len(arg) != 2 {
			continue
		}
		if arg[0] == "user" {
			user = arg[1]
		}
		if arg[0] == "pass" {
			pass = arg[1]
		}
		if arg[0] == "git" {
			git = arg[1]
		}
		continue
	}
	var org, repo string
	if git != "" {
		giturl := strings.Split(git, "/")
		if len(giturl) == 2 {
			org = giturl[0]
			repo = giturl[1]
		}
	}

	_, _ = user, pass

	if *runChild {
		fmt.Println("===welcome===")
		fmt.Println("logbase: ", git)
		child(org, repo)
	}
	//proceed if auth ok
	//auth check here

	// do user init
	// create link? or change user?

	run()

	// check command done

}

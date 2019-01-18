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

	listenPort = flag.Int("p", 0, "local port number")
	dialAddr   = flag.String("dial", "", "dailing address for test listening (eg. localhost:8081)")
)

// the container can't be run directly, there need an parent
func run(org, repo, env string) {

	cmd := exec.Command("/proc/self/exe", append([]string{"-child"}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	cmd.Run()

	//dst := filepath.Join(*dstDir, org, repo)

	// _, err := container.UnMount(dst)
	// if err != nil {
	// 	fmt.Println("umount err", err)
	// }
	//log.Println("exit")
}

// link doesn't mount files, can't set workDir
// workDir can set in the log file structure only
func child(org, repo, env string) {
	//Listen()

	// l := filepath.Join(*srcDir, org, repo)
	// t := filepath.Join("/tmp/", org, repo)
	// src := filepath.Join("/tmp/", org)
	// CreateLink(l, t)
	c := &container.Container{
		Arg:        []string{"sh"},
		Src:        filepath.Join(*srcDir, env, org, repo),
		Rootfs:     *rootFS,
		Dst:        filepath.Join(*dstDir, org, repo),
		BindDst:    filepath.Join(*dstDir, org, repo, "logs"),
		CGroupName: repo,
		Hostname:   repo,
		WorkDir:    filepath.Join("/logs"),
	}
	err := c.Run()
	if err != nil {
		log.Println("run err", err)
	}
	log.Println("exit")
	// RemoveLink(t)
}

func main() {
	flag.Parse()

	if *dialAddr != "" {
		Dial(*dialAddr)
		os.Exit(0)
	}

	if len(os.Args) <= 1 {
		fmt.Println("error: no git provided")
		os.Exit(1)
	}

	var (
		user string
		git  string
		env  string
	)
	for _, v := range os.Args[1:] {
		arg := strings.Split(v, "=")
		if len(arg) != 2 {
			continue
		}
		if arg[0] == "user" {
			user = arg[1]
		}
		if arg[0] == "git" {
			git = arg[1]
		}
		if arg[0] == "env" {
			env = arg[1]
		}
		continue
	}
	if user == "" {
		fmt.Println("user arg not provided")
		return
	}
	if git == "" {
		fmt.Println("git arg not provided")
		return
	}
	if env == "" {
		fmt.Println("env arg not provided")
		return
	}

	var org, repo string
	giturl := strings.Split(git, "/")
	if len(giturl) == 2 {
		org = giturl[0]
		repo = giturl[1]
	}

	if *runChild {
		fmt.Println("===welcome===")
		fmt.Printf("logbase: %v, env: %v\n", git, env)

		err := UserValidate(user, git)
		if err != nil {
			log.Println("user validate error: ", err)
			//return
		}

		child(org, repo, env)
		return
	}
	//proceed if auth ok
	//auth check here

	// do user init
	// create link? or change user?

	run(org, repo, env)

	// check command done

}

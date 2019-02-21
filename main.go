package main

import (
	"flag"
	"fmt"
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
	rootFS = flag.String("rootfs", "/data/alpine", "rootfs on the host( eg. alpine rootfs )")

	runChild    = flag.Bool("child", false, "run container")
	user        = flag.String("user", "", "gitlab user info")
	childgit    = flag.String("childgit", "", "gitlab git(org/repo) for container")
	childuserid = flag.Int("childuserid", 0, "gitlab userid for container")

	listenPort = flag.Int("p", 0, "local port number")
	// dialAddr = flag.String("dial", "", "dailing address for test listening (eg. localhost:8081)")
)

// the container can't be run directly, there need an parent
func run() error {
	// args := append([]string{"-child"}, "-childorg="+org, "-childrepo="+repo, "-childenv="+env)
	args := []string{"-child"}
	args = append(args, os.Args[1:]...)
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
		// Credential:   &syscall.Credential{Uid: 65534, Gid: 65534}, //set at child level
	}
	// fmt.Println("start the child")
	return cmd.Run()

	//dst := filepath.Join(*dstDir, org, repo)

	// _, err := container.UnMount(dst)
	// if err != nil {
	// 	fmt.Println("umount err", err)
	// }
	//log.Println("exit")
}

// link doesn't mount files, can't set workDir
// workDir can set in the log file structure only
func child(org, repo string, envs []string) error {
	//Listen()

	// l := filepath.Join(*srcDir, org, repo)
	// t := filepath.Join("/tmp/", org, repo)
	// src := filepath.Join("/tmp/", org)
	// CreateLink(l, t)
	dst := filepath.Join(*dstDir, org, repo)
	src := filepath.Join(*srcDir, org, repo)

	// err := os.MkdirAll(filepath.Join(dst, "logs/online"), 0755)
	// if err != nil {
	// 	fmt.Println("make err", err)
	// }
	// _, err = os.Stat(filepath.Join(dst, "logs/online"))
	// if err != nil {
	// 	return fmt.Errorf("dst %v does not exist", dst)
	// }
	// fmt.Println("create ok", filepath.Join(dst, "logs/online"))

	okenvs, binds, err := getBinds(src, dst, envs)
	if err != nil {
		return err
	}
	fmt.Println("exist envs: ", okenvs)
	for _, v := range okenvs {
		os.MkdirAll(filepath.Join(dst, "logs", v), 0755)
	}
	c := &container.Container{
		Arg: []string{"sh"},
		// Src:        filepath.Join(dst, env),
		Rootfs: *rootFS,
		Dst:    dst,
		// BindDst:    filepath.Join(dst, "logs"),
		CGroupName: repo,
		Hostname:   repo,
		WorkDir:    filepath.Join("/logs"),
		Binds:      binds,
	}

	return c.Run()
	// if err != nil {
	// 	log.Println("run err", err)
	// 	return err
	// }
	// return nil
	//log.Println("exit")
	// RemoveLink(t)
}

func main() {
	flag.Parse()
	if *GitlabAccessToken == "" {
		fmt.Println("token not set, exit")
		return
	}

	// if *dialAddr != "" {
	// 	Dial(*dialAddr)
	// 	os.Exit(0)
	// }

	// if len(os.Args) <= 1 {
	// 	fmt.Println("error: no git provided")
	// 	os.Exit(1)
	// }

	if *runChild {
		// var (
		// 	user    string
		// 	gitlist string

		// 	env string
		// )
		// for _, v := range os.Args[1:] {
		// 	arg := strings.Split(v, "=")
		// 	if len(arg) != 2 {
		// 		continue
		// 	}
		// 	if arg[0] == "user" {
		// 		user = arg[1]
		// 	}
		// 	if arg[0] == "git" {
		// 		gitlist = UnCompress(arg[1])
		// 	}
		// 	if arg[0] == "env" {
		// 		env = arg[1]
		// 	}
		// 	continue
		// }
		// if user == "" {
		// 	fmt.Println("user arg not provided")
		// 	return
		// }
		// // if git == "" {
		// // 	fmt.Println("git arg not provided")
		// // 	return
		// // }
		// if env == "" {
		// 	//fmt.Println("env arg not provided, using default env: online")
		// 	env = "online"
		// }
		// fmt.Printf("Hi %v\n", *user)

		// fmt.Println("going exist")
		// os.Exit(0)

		// fmt.Println("You can change env by visit: http://logs.devops.haodai.net:8001/?env=pre-online")
		// fmt.Println("You can change repo by visit: http://logs.devops.haodai.net:8001/?git=yunwei/worktile")

		// git := GetProject(gitlist, *srcDir)

		// var org, repo string
		// giturl := strings.Split(git, "/")
		// if len(giturl) == 2 {
		// 	org = giturl[0]
		// 	repo = giturl[1]
		// }

		envs, err := CheckPerm(*childgit, *childuserid)
		if err != nil {
			fmt.Printf("check permission err: %v\n", err)
			return
		}

		k8sgit := strings.Replace(*childgit, "_", "-", -1)

		var org, repo string
		giturl := strings.Split(k8sgit, "/")
		if len(giturl) != 2 {
			fmt.Printf("git %v format err, expect lens 2, got  %v\n", *childgit, len(giturl))
			return
		}
		org = giturl[0]
		repo = giturl[1]

		fmt.Printf("\n=== Welcome ===\n")
		fmt.Printf("logbase: %v, permit envs: %v\n", k8sgit, strings.Join(envs, ","))

		child(org, repo, envs)
		// child(org, repo, env)
		return
	}
	//proceed if auth ok
	//auth check here

	// do user init
	// create link? or change user?

	// runparent(gitlist, *srcDir, env)
	run()
	// check command done

	fmt.Println("exited")
}

// func runparent() {
// 	var i int // to prevent dead loop
// 	for {
// 		i++
// 		//ar x int
// 		// fmt.Println("enter interger input:")
// 		// _, err := fmt.Scanf("%d", &x)
// 		// fmt.Println("got input x:", x, err)
// 		// git := GetProject(gitlist, srcDir)
// 		// git := "yunwei/trx"
// 		// var org, repo string
// 		// giturl := strings.Split(git, "/")
// 		// if len(giturl) == 2 {
// 		// 	org = giturl[0]
// 		// 	repo = giturl[1]
// 		// }

// 		// _, _ = org, repo
// 		fmt.Println("starting child...")
// 		err := run()
// 		if err != nil {
// 			fmt.Println("run err: ", err)
// 			break
// 		}
// 		if i >= 10 {
// 			fmt.Println("it's runned too many times, try refresh the page to enter again")
// 			break
// 		}
// 	}
// }

// func runparent(gitlist, srcDir, env string) {
// 	var i int // to prevent dead loop
// 	for {
// 		i++
// 		//ar x int
// 		// fmt.Println("enter interger input:")
// 		// _, err := fmt.Scanf("%d", &x)
// 		// fmt.Println("got input x:", x, err)
// 		// git := GetProject(gitlist, srcDir)
// 		git := "yunwei/trx"
// 		var org, repo string
// 		giturl := strings.Split(git, "/")
// 		if len(giturl) == 2 {
// 			org = giturl[0]
// 			repo = giturl[1]
// 		}

// 		// _, _ = org, repo
// 		fmt.Println("starting child...")
// 		err := run(org, repo, env)
// 		if err != nil {
// 			fmt.Println("run err: ", err)
// 			break
// 		}
// 		if i >= 10 {
// 			fmt.Println("it's runned too many times, try refresh the page to enter again")
// 			break
// 		}
// 	}
// }

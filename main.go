package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"./k8s"
)

var (
	// srcDir = flag.String("src", "/data/fluentd", "log base directory on the host")
	// dstDir = flag.String("dstDir", "/tmp", "destination directory on the host")
	// rootFS = flag.String("rootfs", "/data/alpine", "rootfs on the host( eg. alpine rootfs )")

	GitlabAccessToken = flag.String("gitlabtoken", "", "gitlab access token")
	GitlabEndpoint    = flag.String("gitlaburl", "http://g.haodai.net", "gitlab base url")

	// runChild = flag.Bool("child", false, "run container")
	token = flag.String("token", "", "gitlab user token")
	// user        = flag.String("user", "", "gitlab user info")
	// childgit    = flag.String("childgit", "", "gitlab git(org/repo) for container")
	// childuserid = flag.Int("childuserid", 0, "gitlab userid for container")

	listenPort = flag.Int("p", 0, "local port number")
	// dialAddr = flag.String("dial", "", "dailing address for test listening (eg. localhost:8081)")
)

var Usage = func() {
	fmt.Printf("Usage of %v:\n", os.Args[0])
	flag.PrintDefaults()

	fmt.Println(`Environments:
	GOTTY_USERTOKEN	-- user token for gitlab access
	`)
}

func parseArgs(args []string) (git, user, token string) {
	token = os.Getenv("GOTTY_USERTOKEN")
	if token == "" {
		fmt.Printf("error, no token found")
		os.Exit(1)
	}
	for _, v := range args[1:] {
		arg := strings.Split(v, "=")
		if len(arg) < 2 {
			continue
		}
		// if arg[0] == "token" {
		// 	// fmt.Println("got user ", arg[1])
		// 	s := strings.Split(v, "token=")
		// 	if len(s) != 2 {
		// 		log.Printf("parse token err expect prefix: token= ")
		// 		continue
		// 	}
		// 	token = s[1]
		// }
		if arg[0] == "user" {
			// fmt.Println("got user ", arg[1])
			s := strings.Split(v, "user=")
			if len(s) != 2 {
				log.Printf("parse user err expect prefix: user= ")
				continue
			}
			user = s[1]
		}
		if arg[0] == "git" {
			// fmt.Println("got user ", arg[1])
			s := strings.Split(v, "git=")
			if len(s) != 2 {
				log.Printf("parse user err expect prefix: git= ")
				continue
			}
			git = s[1]
		}
		continue
	}
	return
}

// // the container can't be run directly, there need an parent
// func run() error {
// 	// args := append([]string{"-child"}, "-childorg="+org, "-childrepo="+repo, "-childenv="+env)
// 	args := []string{"-child"}
// 	args = append(args, os.Args[1:]...)
// 	cmd := exec.Command("/proc/self/exe", args...)
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	cmd.SysProcAttr = &syscall.SysProcAttr{
// 		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
// 		Unshareflags: syscall.CLONE_NEWNS,
// 		// Credential:   &syscall.Credential{Uid: 65534, Gid: 65534}, //set at child level
// 	}
// 	// fmt.Println("start the child")
// 	return cmd.Run()

// 	//dst := filepath.Join(*dstDir, org, repo)

// 	// _, err := container.UnMount(dst)
// 	// if err != nil {
// 	// 	fmt.Println("umount err", err)
// 	// }
// 	//log.Println("exit")
// }

// // link doesn't mount files, can't set workDir
// // workDir can set in the log file structure only
// func child(org, repo string, envs []string) error {
// 	//Listen()

// 	// l := filepath.Join(*srcDir, org, repo)
// 	// t := filepath.Join("/tmp/", org, repo)
// 	// src := filepath.Join("/tmp/", org)
// 	// CreateLink(l, t)
// 	dst := filepath.Join(*dstDir, org, repo)
// 	src := filepath.Join(*srcDir, org, repo)

// 	// err := os.MkdirAll(filepath.Join(dst, "logs/online"), 0755)
// 	// if err != nil {
// 	// 	fmt.Println("make err", err)
// 	// }
// 	// _, err = os.Stat(filepath.Join(dst, "logs/online"))
// 	// if err != nil {
// 	// 	return fmt.Errorf("dst %v does not exist", dst)
// 	// }
// 	// fmt.Println("create ok", filepath.Join(dst, "logs/online"))

// 	okenvs, binds, err := getBinds(src, dst, envs)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("exist envs: ", okenvs)
// 	for _, v := range okenvs {
// 		os.MkdirAll(filepath.Join(dst, "logs", v), 0755)
// 	}
// 	c := &container.Container{
// 		Arg: []string{"sh"},
// 		// Src:        filepath.Join(dst, env),
// 		Rootfs: *rootFS,
// 		Dst:    dst,
// 		// BindDst:    filepath.Join(dst, "logs"),
// 		CGroupName: repo,
// 		Hostname:   repo,
// 		WorkDir:    filepath.Join("/logs"),
// 		Binds:      binds,
// 	}

// 	return c.Run()
// 	// if err != nil {
// 	// 	log.Println("run err", err)
// 	// 	return err
// 	// }
// 	// return nil
// 	//log.Println("exit")
// 	// RemoveLink(t)
// }

func main() {
	flag.Usage = Usage
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

	// if *runChild {
	git, user, token := parseArgs(os.Args)
	if token == "" {
		fmt.Println("token arg not provided")
		return
	}
	fmt.Printf("Hi %v\n", strings.TrimSpace(user))

	_ = git

	var pod k8s.Pod

	// if git == "" {
	//fmt.Printf("\ntry append gitlab info for quicker access:\n")
	// fmt.Printf("\n想快点？ 添加Gitlab项目信息直接进入:\n")
	// fmt.Printf("示例:    http://logs.devops.haodai.net:8001/?git=flow_center/df-openapi\n\n")
	cancelprint := printprogress()
	admin, gitlist, err := GetProjectLists(token)
	cancelprint()
	if err != nil {
		fmt.Printf("get project lists err: %v\n", err)
		return
	}

	pod, err = GetProjectFromInput(gitlist, admin)
	if err != nil {
		fmt.Println("get project err: ", err)
		os.Exit(1)
	}
	// git = pod.Namespace + "/" + pod.Name
	// }

	if !admin {
		// check user's permission, need to ignore no-exist error
		envs, err := CheckPerm(pod.GitName, token)
		if err != nil {
			// if !strings.Contains(err.Error(), "Project Not Found") {
			fmt.Printf("check permission err: %v, for git: %v\n", err, pod.GitName)
			return
			// }
		}
		if !envok(pod.Env, envs) {
			fmt.Printf("env: %v permission not allowed, allowed env: %v\n", pod.Env, envs)
			return
		}
	}
	// k8sgit := strings.Replace(git, "_", "-", -1)

	// var org, repo string
	// giturl := strings.Split(k8sgit, "/")
	// if len(giturl) != 2 {
	// 	fmt.Printf("git %v format err, expect lens 2, got  %v\n", git, len(giturl))
	// 	return
	// }
	// org = giturl[0]
	// repo = giturl[1]

	fmt.Printf("\n=== Welcome ===\n")
	// fmt.Printf("logbase: %v, permit envs: %v\n", k8sgit, strings.Join(envs, ","))

	run(pod.Namespace, pod.PodName)
	return
	// }

	// err := run()
	// if err != nil {
	// 	fmt.Printf("run err: %v\n", err)
	// 	return
	// }

	fmt.Println("exited")
	fmt.Printf("\nTry refresh the page to enter again.\n")
}

func run(ns, pod string) (out string, err error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("kubectl exec -it -n %v %v", ns, pod))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("build execute build err: %v\noutput: %v\n", err, string(output))
		return
	}
	out = string(output)
	return
}

func envok(env string, envs []string) bool {
	for _, v := range envs {
		if env == v {
			return true
		}
	}
	return false
}

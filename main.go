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
	rootFS = flag.String("rootfs", "/data/alpine", "rootfs on the host( eg. alpine rootfs )")

	GitlabAccessToken = flag.String("gitlabtoken", "", "gitlab access token")
	GitlabEndpoint    = flag.String("gitlaburl", "http://g.haodai.net", "gitlab base url")

	runChild = flag.Bool("child", false, "run container")
	token    = flag.String("token", "", "gitlab user token")
	// user        = flag.String("user", "", "gitlab user info")
	// childgit    = flag.String("childgit", "", "gitlab git(org/repo) for container")
	// childuserid = flag.Int("childuserid", 0, "gitlab userid for container")

	listenPort = flag.Int("p", 0, "local port number")
	// dialAddr = flag.String("dial", "", "dailing address for test listening (eg. localhost:8081)")
)

func parseArgs(args []string) (user, token string) {
	for _, v := range args[1:] {
		arg := strings.Split(v, "=")
		if len(arg) < 2 {
			continue
		}
		if arg[0] == "token" {
			// fmt.Println("got user ", arg[1])
			ts := strings.Split(v, "token=")
			if len(ts) != 2 {
				log.Printf("parse token err expect prefix: token= ")
				continue
			}
			token = ts[1]
		}
		if arg[0] == "user" {
			// fmt.Println("got user ", arg[1])
			us := strings.Split(v, "user=")
			if len(us) != 2 {
				log.Printf("parse user err expect prefix: user= ")
				continue
			}
			user = us[1]
		}
		continue
	}
	return
}

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
		user, token := parseArgs(os.Args)
		if token == "" {
			fmt.Println("token arg not provided")
			return
		}

		cancelprint := printprogress()
		gitlist, err := GetProjectLists(token, *srcDir)
		cancelprint()
		if err != nil {
			fmt.Printf("get project lists err: %v\n", err)
			return
		}

		git, err := GetProjectFromInput(gitlist, *srcDir)

		if err != nil {
			fmt.Println("get project err: ", err)
			os.Exit(1)
		}

		// check user's permission.
		envs, err := CheckPerm(git, token)
		if err != nil {
			fmt.Printf("check permission err: %v\n", err)
			return
		}

		k8sgit := strings.Replace(git, "_", "-", -1)

		var org, repo string
		giturl := strings.Split(k8sgit, "/")
		if len(giturl) != 2 {
			fmt.Printf("git %v format err, expect lens 2, got  %v\n", git, len(giturl))
			return
		}
		org = giturl[0]
		repo = giturl[1]

		fmt.Printf("\n=== Welcome %v ===\n", user)
		fmt.Printf("logbase: %v, permit envs: %v\n", k8sgit, strings.Join(envs, ","))

		child(org, repo, envs)
		return
	}

	err := run()
	if err != nil {
		fmt.Printf("run err: %v\n", err)
		return
	}

	fmt.Println("exited")
	fmt.Printf("\nTry refresh the page to enter again.\n")
}

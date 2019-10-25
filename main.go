package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

func parseArgs(args []string) (git, user, token, pod string) {
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

		if arg[0] == "pod" {
			// fmt.Println("got user ", arg[1])
			s := strings.Split(v, "pod=")
			if len(s) != 2 {
				log.Printf("parse user err expect prefix: pod= ")
				continue
			}
			pod = s[1]
		}

		continue
	}
	return
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	if *GitlabAccessToken == "" {
		log.Println("gitlab token not set, exit")
		return
	}

	git, user, token, podname := parseArgs(os.Args)
	if user == "" {
		fmt.Println("user arg not provided")
		return
	}
	if token == "" {
		fmt.Println("token arg not provided")
		return
	}
	fmt.Printf("Hi %v\n", strings.TrimSpace(user))

	admin, err := IsAdmin(token)
	if err != nil {
		fmt.Printf("check admin err: %v\n", err)
		return
	}

	// if git == "" {
	//fmt.Printf("\ntry append gitlab info for quicker access:\n")
	// fmt.Printf("\n想快点？ 添加Gitlab项目信息直接进入:\n")
	// fmt.Printf("示例:    http://logs.devops.haodai.net:8001/?git=flow_center/df-openapi\n\n")

	// log.Printf("%#v\n", pod)
	// git = ns + "/" + pod.Name
	// }

	var ns, env string
	if podname == "" {
		// get pod from podname

		pod, err := GetProject(token, admin)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		git = pod.GitName
		env = pod.Env
		ns = pod.Namespace
		podname = pod.PodName
	} else {
		ns, _ = getnsrepo(git)
		env = getenv(podname)
	}

	if git == "" {
		fmt.Printf("git is empty, git and podname should be both provided\n")
		return
	}
	if podname == "" {
		fmt.Printf("podname is empty, git and podname should be both provided\n")
		return
	}
	if ns == "" {
		fmt.Printf("derive ns is empty\n")
		return
	}
	if env == "" {
		fmt.Printf("derive env is empty\n")
		return
	}

	if !admin {
		// check user's permission, need to ignore no-exist error
		envs, err := CheckPerm(git, token)
		if err != nil {
			envs, err = CheckPerm(strings.Replace(git, "-", "_", -1), token)
			if err != nil {
				// if !strings.Contains(err.Error(), "Project Not Found") {
				fmt.Printf("check permission err: %v, for git: %v\n", err, git)
				return
				// }
			}
		}
		if !envok(env, envs) {
			fmt.Printf("env: %v permission not allowed, allowed env: %v\n", env, envs)
			return
		}
	}

	fmt.Printf("\n=== Welcome %v ===\n", user)

	fmt.Printf("Entering ns: %v, pod: %v, env: %v\n", ns, podname, env)
	if env == "test" {
		fmt.Printf("\n\nnote that: pod in test env has separate network, so not reachable\n")
	}
	out, err := runterm(user, ns, podname)
	if err != nil {
		fmt.Printf("run err: %v\noutput: %v\n", err, out)

		return
	}

	fmt.Println("exited")
	fmt.Printf("\nTry refresh the page to enter again.\n")
}

func getnsrepo(git string) (ns, repo string) {
	k8sgit := strings.Replace(git, "_", "-", -1)
	giturl := strings.Split(k8sgit, "/")
	if len(giturl) >= 1 {
		ns = giturl[0]
	}
	if len(giturl) >= 2 {
		repo = giturl[1]
	}
	return
}

func getenv(podname string) (env string) {
	p := strings.Split(podname, "-")
	if len(p) >= 3 {
		env = p[len(p)-3]
		return
	}
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

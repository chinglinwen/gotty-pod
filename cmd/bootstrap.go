package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	srcDir = flag.String("src", "/data/fluentd", "log base directory on the host")
	dstDir = flag.String("dstDir", "/tmp", "destination directory on the host")
	rootFS = flag.String("rootfs", "/data/alpine", "rootfs on the host( eg. alpine rootfs )")

	GitlabAccessToken = flag.String("gitlabtoken", "", "gitlab access token")
	GitlabEndpoint    = flag.String("gitlaburl", "http://g.haodai.net", "gitlab base url")

	// runChild  = flag.Bool("child", false, "run container")
	// user      = flag.String("user", "", "user info")
	// childorg  = flag.String("childorg", "", "org for container")
	// childrepo = flag.String("childrepo", "", "repo for container")
	// childenv  = flag.String("childenv", "", "env for container")

	// listenPort = flag.Int("p", 0, "local port number")
	// dialAddr   = flag.String("dial", "", "dailing address for test listening (eg. localhost:8081)")

	gottyCMD = flag.String("cmd", "./gotty-logs", "gotty-logs command")
)

// the container can't be run directly, there need an parent
func run(user, git string, userid int) error {
	args := []string{"-childgit=" + git, "-childuserid=" + strconv.Itoa(userid), "-user=" + user,
		"-src=" + *srcDir, "-dstDir=" + *dstDir, "-rootfs=" + *rootFS,
		"-gitlabtoken=" + *GitlabAccessToken, "-gitlaburl=" + *GitlabEndpoint}
	// args = append(args, os.Args[1:]...)

	cmd := exec.Command(*gottyCMD, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// fmt.Println("start the gotty-logs")
	return cmd.Run()
}

func main() {
	flag.Parse()

	// fmt.Printf("%#v\n", os.Args)
	var (
		username string
		userid   int
		gitlist  string

		err error
	)
	for _, v := range os.Args[1:] {
		arg := strings.Split(v, "=")
		if len(arg) < 2 {
			continue
		}
		if arg[0] == "user" {
			// fmt.Println("got user ", arg[1])
			us := strings.Split(v, "user=")
			if len(us) != 2 {
				log.Printf("parse user err expect prefix: user= ")
				continue
			}
			username, userid, err = ParseUserInfo(us[1])
			if err != nil {
				log.Printf("parse user err: %v, user: %v\n", err, us[1])
			}
			// fmt.Println("got ", username, userid)
		}
		if arg[0] == "git" {
			gits := strings.Split(v, "git=")
			if len(gits) != 2 {
				log.Printf("parse git err expect prefix: git= ")
				continue
			}
			gitlist, err = UnCompress(gits[1])
			if err != nil {
				log.Printf("parse git err: %v, git: %v\n", err, gits[1])
			}
		}
		// if arg[0] == "env" {
		// 	env = arg[1]
		// }
		continue
	}
	if username == "" {
		fmt.Println("username arg not provided")
		return
	}
	if gitlist == "" {
		fmt.Println("gitlist arg not provided")
		return
	}
	// if env == "" {
	// 	//fmt.Println("env arg not provided, using default env: online")
	// 	env = "online"
	// }
	fmt.Printf("Hi %v, id: %v\n", username, userid)

	//proceed if auth ok
	//auth check here

	// do user init
	// create link? or change user?

	runparent(username, gitlist, *srcDir, userid)

	// check command done

}

func runparent(user, gitlist, srcDir string, userid int) {
	// var i int // to prevent dead loop
	// for {
	// 	i++

	git, err := GetProject(gitlist, srcDir)
	//k8sgit := strings.Replace(git, "_", "-", -1)

	if err != nil {
		fmt.Println("get project err: ", err)
		os.Exit(1)
	}
	// var org, repo string
	// giturl := strings.Split(git, "/")
	// if len(giturl) == 2 {
	// 	org = giturl[0]
	// 	repo = giturl[1]
	// }

	// _, _ = org, repo
	// fmt.Println("starting...")
	err = run(user, git, userid)
	if err != nil {
		fmt.Println("run err: ", err)
		// break
	}
	fmt.Printf("\nTry refresh the page to enter again.\n")
	// if i >= 10 {
	// 	fmt.Println("it's runned too many times, try refresh the page to enter again")
	// 	break
	// }
	// }
}

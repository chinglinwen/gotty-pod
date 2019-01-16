package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
)

var (
	logBase = flag.String("b", "/data/fluentd", "log base directory on the host ( eg. /data/fluentd )")
)

func init() {
	if *logBase == "" {
		*logBase, _ = os.Getwd()
	}
}

func main() {
	flag.Parse()

	fmt.Println("===welcome===")
	fmt.Println("logbase: ", *logBase)
	if len(os.Args) <= 2 {
		fmt.Println("error: not enough arguments")
		os.Exit(1)
	}
	fmt.Println(os.Args[1:])
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
	fmt.Printf("user: %v, pass: %v\n", user, pass)

	//proceed if auth ok
	//auth check here

	// do user init
	// create link? or change user?

	// check command done

	err := CreateLink(*logBase+"/"+git, git)
	if err != nil {
		fmt.Printf("try enter into %v error: %v\n", git, err)
		return
	}
	fmt.Printf("enter into git: %v, org: %v, repo: %v\n", git, org, repo)
	runShell(git)

	//RemoveLink(git)
}

func runShell(git string) {
	pr, pw := io.Pipe()

	//cmd := exec.Command("bash", "-c", "./a.sh")
	cmd := exec.Command("bash")
	cmd.Dir = git
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	// this doesn't work
	// cmd.Env = os.Environ()
	// cmd.Env = append(cmd.Env, "PS1=%m %{${fg_bold[blue]}%}:: %{$reset_color%}%{${fg[green]}%}%3~ $(git_prompt_info)%{${fg_bold[$CARETCOLOR]}%}»%{${reset_color}%}")

	cmd.Stdin = pr
	cmd.Start()

	ch := make(chan string, 100)
	chInput := make(chan string, 100)
	stdinScan := bufio.NewScanner(os.Stdin)
	go func() {
		fmt.Printf("[~] $: ")
		for stdinScan.Scan() {
			//fmt.Printf("[~] $: ")
			line := stdinScan.Text()
			if line == "\n" || line == "\r" || line == "" {
				continue
			}
			err := CommandCheck(line)
			if err != nil {
				fmt.Println(err)
				continue
			}
			chInput <- line
			//fmt.Printf("got input: %v\n", string(line))
		}
	}()

	// stdin
	go func() {
		for v := range chInput {
			fmt.Fprintf(pw, "%v\n", v)
			// if err != nil {
			// 	fmt.Printf("write stdin err: %v\n", err)
			// }
		}
	}()

	stdoutScan := bufio.NewScanner(stdout)
	stderrScan := bufio.NewScanner(stderr)
	go func() {
		for stdoutScan.Scan() {
			line := stdoutScan.Text()
			ch <- line
		}
	}()
	go func() {
		for stderrScan.Scan() {
			line := stderrScan.Text()
			ch <- line
		}
	}()

	go func() {
		cmd.Wait()
		close(ch)
		close(chInput)
	}()

	// stdout
	for line := range ch {
		fmt.Printf("%v\n", line)
	}
}

var forbidList = []string{
	"[.][.]",
	"[<>]",
	";",
}

// check at compiletime
func init() {
	for _, v := range forbidList {
		_ = regexp.MustCompile(v)
	}
}

var whiteList = []string{
	"tail",
	"grep",
	"less",
	"more",
	"cat",
	//"pwd",
	"whoami",
	"ls",
	"echo",
	"exit",
	"head",
	"tail",
}

// make some command can't execute, such as cd to parent directory
func CommandCheck(text string) error {
	cmdname := strings.Split(text, " ")[0]
	for _, v := range whiteList {
		if cmdname == v {
			return charsCheck(text)
		}
	}
	return fmt.Errorf("forbid command %s, allowed cmds: %v\n", cmdname, whiteList)
}

func charsCheck(text string) error {
	for _, v := range forbidList {
		if match := regexp.MustCompile(v).FindString(text); match != "" {
			return fmt.Errorf("forbid character %q in %s\n", match, text)
		}
	}
	return nil
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			// here we just ignore it, to pass it to bash only
		}
	}()
}

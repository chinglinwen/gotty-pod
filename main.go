package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
)

func main() {
	fmt.Println("===welcome===")
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
	fmt.Printf("user: %v, pass: %v\n", user, pass)

	//proceed if auth ok
	//auth check here

	// do user init
	// create link? or change user?

	// check command done

	CreateLink("test", git)
	fmt.Printf("enter into git: %v\n", git)
	runShell(git)
}

func runShell(git string) {
	pr, pw := io.Pipe()

	//cmd := exec.Command("bash", "-c", "./a.sh")
	cmd := exec.Command("bash")
	cmd.Dir = git
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "PS1=%m %{${fg_bold[blue]}%}:: %{$reset_color%}%{${fg[green]}%}%3~ $(git_prompt_info)%{${fg_bold[$CARETCOLOR]}%}»%{${reset_color}%}")

	cmd.Stdin = pr
	cmd.Start()

	ch := make(chan string, 100)
	chInput := make(chan string, 100)
	stdinScan := bufio.NewScanner(os.Stdin)
	stdinScan.Split(bufio.ScanRunes)
	go func() {
		fmt.Printf("[~] $: ")
		var line string
		for stdinScan.Scan() {
			//fmt.Printf("[~] $: ")
			c := stdinScan.Text()
			if c != "\n" {
				line += c
			}
			if c == "\n" {
				err := CommandCheck(line)
				line = ""
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
			chInput <- c
			//fmt.Printf("got input: %v\n", string(line))
		}
	}()

	// stdin
	go func() {
		for v := range chInput {
			fmt.Fprintf(pw, "%v", v)
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
}

// make some command can't execute, such as cd to parent directory
func CommandCheck(text string) error {
	for _, v := range forbidList {
		if matched, _ := regexp.MatchString(v, text); matched {
			return fmt.Errorf("forbid command %s", text)
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

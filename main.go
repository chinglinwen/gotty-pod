package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("error: not enough arguments")
		os.Exit(1)
	}
	fmt.Println(os.Args[1:])
	var user, pass string
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
		continue
	}
	fmt.Printf("user: %v, pass: %v\n", user, pass)

	//proceed if auth ok
	//auth check here

	// do user init
	// create link? or change user?

	// check command done

	CreateLink()
	runShell()
}

func runShell() {
	pr, pw := io.Pipe()

	cmd := exec.Command("bash", "-c", "./a.sh")
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Stdin = pr
	cmd.Start()

	chInput := make(chan string, 100)
	stdinScan := bufio.NewScanner(os.Stdin)
	go func() {
		for stdinScan.Scan() {
			line := stdinScan.Text()
			err := CommandCheck(line)
			if err != nil {
				fmt.Println(err)
				continue
			}
			chInput <- line
			fmt.Printf("got input: %v\n", string(line))
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

	ch := make(chan string, 100)

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
		fmt.Println(line)
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

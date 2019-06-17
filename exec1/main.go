package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	c := exec.Command("/bin/bash", "-c", "kubectl exec -it -n flow-center tangguo-online-597fc44cb4-q8lsc sh")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	go io.Copy(f, os.Stdin)
	io.Copy(os.Stdout, f)

	fmt.Print("exiting\r\n")
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	out, err := runterm("flow-center", "tangguo-pre-6f575795df-fncfb")
	if err != nil {
		log.Fatal("run err", err)
	}
	fmt.Println(out)
}

func run(ns, pod string) (out string, err error) {
	s := fmt.Sprintf("kubectl exec -it -n %v %v sh", ns, pod)
	log.Println("executing: ", s)
	cmd := exec.Command("sh", "-c", s)
	// cmd := exec.Command("sh")
	cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("build execute build err: %v\noutput: %v\n", err, string(output))
		return
	}
	out = string(output)
	return
}

func runtermok(ns, pod string) (out string, err error) {
	s := fmt.Sprintf("kubectl exec -it -n %v %v sh", ns, pod)
	log.Println("executing: ", s)
	c := exec.Command("sh", "-c", s)

	// c := exec.Command("/bin/bash", "-c", "kubectl exec -it -n flow-center tangguo-online-597fc44cb4-q8lsc sh")
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

	return
}

func runterm(ns, pod string) (out string, err error) {
	s := fmt.Sprintf("kubectl exec -it -n %v %v sh", ns, pod)
	log.Println("executing: ", s)
	c := exec.Command("sh", "-c", s)

	// stdinIn, _ := c.StdinPipe()

	pr, pw := io.Pipe()

	// c := exec.Command("/bin/bash", "-c", "kubectl exec -it -n flow-center tangguo-online-597fc44cb4-q8lsc sh")
	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	// go io.Copy(f, os.Stdin)

	// var buf bytes.Buffer
	tee := io.TeeReader(os.Stdin, pw)
	defer pw.Close()

	go func() {
		// defer pw.Close()
		// fmt.Fprintf(pw, "ls\r\n")
		// return

		// copy the data written to the PipeReader via the cmd to stdout
		// if _, err := io.Copy(os.Stdout, pr); err != nil {
		// 	log.Fatal(err)
		// }

		// scanner := bufio.NewScanner(pr)
		// for scanner.Scan() {
		// 	// ucl := strings.ToUpper(scanner.Text())
		// 	text := scanner.Text()
		// 	fmt.Println("got ", text)

		// 	// fmt.Println("got", ucl)
		// 	write(text + "\n")
		// 	// fmt.Fprintf(f, "%v", text)
		// }

		reader := bufio.NewReader(pr)

		for {
			input, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			// fmt.Printf("%c", input)
			if input == '\r' {
				input = '\n'
			}
			write(fmt.Sprintf("%c", input))
			// output = append(output, input)
		}

		// var output []rune

		// for {
		// 	input, _, err := reader.ReadRune()
		// 	if err != nil && err == io.EOF {
		// 		break
		// 	}
		// 	// fmt.Printf("%c", input)
		// 	write(fmt.Sprintf("%c", input))
		// 	// output = append(output, input)
		// }

		// write(fmt.Sprintf("got : %v", len(output)))

		// for j := 0; j < len(output); j++ {
		// 	// fmt.Printf("got %c", output[j])
		// 	write(fmt.Sprintf("got %c", output[j]))
		// }
	}()

	// go io.Copy(f, os.Stdin)
	go io.Copy(f, tee)
	// go io.Copy(f, stdinIn)
	// go func() {
	// 	_, errStdout = io.Copy(stdout, stdoutIn)
	// }()

	io.Copy(os.Stdout, f)

	return
}

func run2(ns, pod string) (out string, err error) {
	// Could read $PAGER rather than hardcoding the path.
	// cmd := exec.Command("/usr/bin/less", file)
	// cmd := exec.Command("/usr/bin/cat")
	s := fmt.Sprintf("kubectl exec -it -n %v %v sh", ns, pod)
	log.Println("executing: ", s)
	cmd := exec.Command("sh", "-c", s)

	pr, pw := io.Pipe()
	// defer pw.Close()

	go func() {
		// defer pr.Close()
		// copy the data written to the PipeReader via the cmd to stdout
		if _, err := io.Copy(os.Stdout, pr); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer pw.Close()
		// copy the data written to the PipeReader via the cmd to stdout
		// if _, err := io.Copy(os.Stdout, pr); err != nil {
		// 	log.Fatal(err)
		// }
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// ucl := strings.ToUpper(scanner.Text())
			text := scanner.Text()
			// fmt.Println("got", ucl)
			write(text + "\n")
			fmt.Fprintf(pw, "%v", text)
		}
	}()

	// Feed it with the string you want to display.
	// cmd.Stdin = strings.NewReader("The text you want to show.")
	// cmd.Stdin = os.Stdin
	cmd.Stdin = pr

	// This is crucial - otherwise it will write to a null device.
	cmd.Stdout = os.Stdout

	// Fork off a process and wait for it to terminate.
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func write(line string) {
	f, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write([]byte(line)); err != nil {
		log.Fatal("write err", err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func runpty(ns, pod string) (out string, err error) {
	// c := exec.Command("grep", "--color=auto", "bar")
	s := fmt.Sprintf("kubectl exec -it -n %v %v sh", ns, pod)
	log.Println("executing: ", s)
	c := exec.Command("sh", "-c", s)
	// pw, s_ := c.StdinPipe()

	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	// oldState, err := terminal.MakeRaw(0)
	// if err != nil {
	// 	panic(err)
	// }
	// defer terminal.Restore(0, oldState)

	// pr, pw := io.Pipe()
	go func() {
		// f.Write([]byte("ls\n"))
		// f.Write([]byte("bar\n"))
		// f.Write([]byte("baz\n"))
		// f.Write([]byte{4}) // EOT

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// ucl := strings.ToUpper(scanner.Text())
			text := scanner.Text()
			// fmt.Println("got", ucl)
			write("got " + text + "\n")
			// fmt.Fprintf(pw, "%v", text)
			// f.Write([]byte("ls\n"))
			f.Write([]byte(text + "\n"))
		}
		f.Write([]byte{4}) // eof
	}()

	// go func() {
	// 	io.Copy(os.Stdin, f)
	// }()

	io.Copy(os.Stdout, f)
	return
}

func runstdin(ns, pod string) (out string, err error) {
	// c := exec.Command("grep", "--color=auto", "bar")
	s := fmt.Sprintf("kubectl exec -it -n %v %v sh", ns, pod)
	log.Println("executing: ", s)
	c := exec.Command("sh", "-c", s)
	pw, _ := c.StdinPipe()
	defer pw.Close()

	f, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	// pr, pw := io.Pipe()
	go func() {
		// f.Write([]byte("ls\n"))
		// f.Write([]byte("bar\n"))
		// f.Write([]byte("baz\n"))
		// f.Write([]byte{4}) // EOT

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			// ucl := strings.ToUpper(scanner.Text())
			text := scanner.Text()
			// fmt.Println("got", ucl)
			// write("got " + text + "\n")
			fmt.Fprintf(pw, "%v", text)
			// f.Write([]byte("ls\n"))
			// f.Write([]byte(text + "\n"))
		}
		f.Write([]byte{4})
	}()

	// go func() {
	// 	io.Copy(os.Stdin, f)
	// }()

	io.Copy(os.Stdout, f)
	return
}

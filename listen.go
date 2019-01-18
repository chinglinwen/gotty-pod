package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func showLocalAddrs() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		fmt.Println(addr.String())
	}
}

// Listen - receive function
func listen(port int) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer lis.Close()

	//go readoutput()

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("Accept error:", err)
		}
		log.Println("accept:", conn.RemoteAddr())

		go func(c net.Conn) {
			//io.Copy(os.Stdin, c)
			io.Copy(os.Stdout, c)
			log.Println("closed:", conn.RemoteAddr())
			defer c.Close()
		}(conn)
	}
}

// Dial - send function
func Dial(addr string) error {
	// need to set timeout
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	go io.Copy(os.Stdout, conn)
	_, err = io.Copy(conn, os.Stdin)
	return err
}

func Listen() {
	// skip listening if not set
	if *listenPort == 0 {
		return
	}
	if !*runChild {
		return
	}
	go func() {
		log.Fatal(listen(*listenPort))
	}()
}

func readoutput() {
	chOutput := make(chan string, 100)
	stdoutScan := bufio.NewScanner(os.Stdout)
	stderrScan := bufio.NewScanner(os.Stderr)
	go func() {
		for stdoutScan.Scan() {
			line := stdoutScan.Text()
			chOutput <- line
		}
	}()
	go func() {
		for stderrScan.Scan() {
			line := stderrScan.Text()
			chOutput <- line
		}
	}()

	defer close(chOutput)

	// stdout
	for line := range chOutput {
		fmt.Printf("%v\n", line)
	}

}

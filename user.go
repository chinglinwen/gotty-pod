package main

import (
	"fmt"
	"syscall"

	"./gitlab"
	"golang.org/x/crypto/ssh/terminal"
)

func UserValidate(user, git string) (err error) {
	// backup stdin
	// oldState, err := terminal.MakeRaw(0)
	// if err != nil {
	// 	panic(err)
	// }
	// defer terminal.Restore(0, oldState)

	//fmt.Printf("user: ")

	// r := bufio.NewReader(os.Stdin)
	// user, _, err := r.ReadLine()
	// if err != nil {
	// 	fmt.Printf("read error: %v\n", err)
	// }

	var p []string
	for i := 0; i < 3; i++ {
		fmt.Printf("password: ")
		pass, e := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if e != nil {
			err = fmt.Errorf("read error: %v\n", e)
			continue
		}

		p, e = gitlab.GetProjects(string(user), string(pass))
		if e != nil {
			fmt.Printf("verify err: %v, try again\n", e)
			err = e
			continue
		}
		err = nil
		break
	}
	if err != nil {
		fmt.Println("too much of try error")
		return
	}
	found := gitlab.ProjectSearch(p, git)
	if !found {
		fmt.Printf("no project found for this user\n")
		// err = fmt.Errorf("no project found for this user\n")
		// return
	}
	return
}

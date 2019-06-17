package main

import (
	"fmt"
	"strconv"
	"strings"

	"./k8s"
	prompt "github.com/c-bata/go-prompt"
)

// type inputentry struct {
// 	id   int
// 	text string
// }

// Get project by user input
func GetProjectFromInput(gitlist []string, admin bool) (pod k8s.Pod, err error) {
	podlist, e := listpods()
	if e != nil {
		err = fmt.Errorf("walk error %v", e)
		return
	}
	var list []string
	for k := range podlist {
		list = append(list, k)
	}
	loglist := list
	if !admin {
		loglist = Filter(list, gitlist)
	}

	n := len(loglist)
	inputlist := make([]string, n+1)

	for i, v := range loglist {
		// fmt.Println("convert gitlist ", i, v)
		inputlist[i+1] = strconv.Itoa(i+1) + " " + v
	}
	text := "0 quit"
	inputlist[0] = text

	// spew.Dump("inputlist", inputlist)

	completer := func(d prompt.Document) []prompt.Suggest {
		s := []prompt.Suggest{}
		for i, v := range inputlist {
			if i == 0 {
				continue
			}
			s = append(s, prompt.Suggest{
				Text: v,
			})
		}
		s = append(s, prompt.Suggest{
			Text: inputlist[0],
		})
		return prompt.FilterFuzzy(s, d.GetWordBeforeCursor(), true)
		// return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
	}

	var p string
	for {
		fmt.Printf("\rPlease select the project ( keyword search is ok ): \n")
		for i, v := range inputlist {
			if i == 0 {
				continue
			}
			fmt.Println("  ", v)
		}
		fmt.Println("  ", inputlist[0])
		t := prompt.Input("> ", completer)
		if t == "" {
			fmt.Println("input is invalid")
			continue
		}
		// spew.Dump("t", t)
		x := searchInput(t, inputlist)
		if x == -1 {
			fmt.Println("input is invalid, not found any")
			continue
		}

		if x == 0 {
			err = fmt.Errorf("exited")
			fmt.Println("exiting...")
			return
		}
		// // index, _ := strconv.Atoi(t)
		// index := x
		// if index <= 0 || index > len(loglist) {
		// 	fmt.Println("input is invalid, input number only")
		// 	continue
		// }
		p = loglist[x-1]
		break
	}
	fmt.Println("You selected " + p)

	pod = podlist[p]
	return
}

func searchInput(t string, inputlist []string) int {
	for i, v := range inputlist {
		// fmt.Println("start ", i, v, t)
		if strings.Contains(v, t) {
			return i
		}
	}
	return -1
}

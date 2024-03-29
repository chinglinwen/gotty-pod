package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	// "./k8s"
	"wen/gotty-pod/k8s"
	// "github.com/AlecAivazis/survey"

	prompt "github.com/c-bata/go-prompt"
	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/AlecAivazis/survey.v1"
)

func VerifyPermission() (err error) {
	token, err := gettokenfrominput()
	if err != nil {
		return
	}
	// fmt.Printf("got token: %q", token)
	err = validateJWT(token)
	if err != nil {
		err = fmt.Errorf("validate token err %v", err)
		return
	}
	return
}

func gettokenfrominput() (token string, err error) {
	prompt := &survey.Input{
		Message: "Enter token for permission( ask yunwei for it): ",
	}
	err = survey.AskOne(prompt, &token, nil)
	if err != nil {
		err = fmt.Errorf("enter token err %v", err)
		return
	}
	// token = answers.Token
	return
}

var SecretKey = "secret"

func validateJWT(token string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}
	t, err := jwt.Parse(token, func(tok *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err == nil && t.Valid {
		return nil
	} else {
		return fmt.Errorf("Invalid")
	}
}

// type inputentry struct {
// 	id   int
// 	text string
// }

func GetProject(token string, admin bool) (pod k8s.Pod, err error) {
	cancelprint := printprogress()
	grouplist, err := GetGroupLists(token)
	cancelprint()
	if err != nil {
		err = fmt.Errorf("get project lists err: %v", err)
		return
	}
	pod, err = GetProjectFromInput(grouplist, admin)
	if err != nil {
		err = fmt.Errorf("get project from input err: %v", err)
		return
	}
	return
}

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
	sort.Slice(loglist, func(i, j int) bool {
		return loglist[i] < loglist[j]
	})

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
		// return prompt.FilterFuzzy(s, d.GetWordBeforeCursor(), true)
		return prompt.FilterContains(s, d.GetWordBeforeCursor(), true)
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
	s := fmt.Sprintf(".*%v.*", strings.ReplaceAll(t, " ", ".*"))
	for i, v := range inputlist {
		// fmt.Println("start ", i, v, t)
		// if strings.Contains(v, t) {
		if matched, _ := regexp.MatchString(s, v); matched {
			return i
		}
	}
	return -1
}

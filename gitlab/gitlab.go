package gitlab

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"

	gitlab "github.com/xanzy/go-gitlab"
)

var (
	GITLAB = "http://g.haodai.net"
)

func getclient(user, pass string) (c *gitlab.Client, err error) {
	c, err = gitlab.NewBasicAuthClient(nil, GITLAB, user, pass)
	if err != nil {
		err = fmt.Errorf("%v clienterr", strings.Split(err.Error(), "\n")[0])
	}
	return
}

func GitlabVerify(user, pass string) (ok bool, err error) {
	_, err = getclient(user, pass)
	if err != nil {
		return
	}
	return true, nil
}

func GetProjects(user, pass string) (projects []string, err error) {
	c, err := getclient(user, pass)
	if err != nil {
		return
	}
	a := true
	ps, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Membership: &a,
	})
	if err != nil {
		err = fmt.Errorf("%v", strings.Split(err.Error(), "\n")[0])

		spew.Dump("err here")
		return
	}
	for _, p := range ps {
		projects = append(projects, p.WebURL)
	}
	return
}

func ProjectSearch(projects []string, git string) bool {
	for _, v := range projects {
		if strings.Contains(v, git) {
			return true
		}
	}
	return false
}

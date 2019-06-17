package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	gitlab "github.com/xanzy/go-gitlab"
)

const (
	EnvOnline    = "online"
	EnvPreOnline = "pre"
	EnvTest      = "test"
)

// var client *gitlab.Client

func client() *gitlab.Client {
	client := gitlab.NewClient(http.DefaultClient, *GitlabAccessToken)
	client.SetBaseURL(*GitlabEndpoint)
	return client
}

func userclient(token string) *gitlab.Client {
	client := gitlab.NewClient(http.DefaultClient, token)
	client.SetBaseURL(*GitlabEndpoint)
	return client
}

func GetUser(token string) (user *gitlab.User, err error) {
	c := userclient(token)
	u, _, err := c.Users.CurrentUser()
	if err != nil {
		log.Println("getuser err", err)
		return
	}
	return u, nil
}

func GetGroups(token string) (c *gitlab.Client, gs []*gitlab.Group, err error) {
	c = userclient(token)
	ps, _, e := c.Groups.ListGroups(&gitlab.ListGroupsOptions{})
	if err != nil {
		err = fmt.Errorf("%v", strings.Split(e.Error(), "\n")[0])
		return
	}
	if len(ps) == 0 {
		err = fmt.Errorf("group: there's no any git group")
		return
	}
	gs = ps
	return
}

func GetGroupLists(token string) (gs []string, err error) {
	// for all group projects
	_, gss, err := GetGroups(token)
	if err != nil {
		log.Println("getgroups err", err)
		return
	}
	for _, g := range gss {
		// log.Println("g", g.Path)
		gs = append(gs, g.Path)
	}
	return
}

// shomehow miss some group, don't use it
// func localexist(list []string, name string) bool {
// 	k8sname := strings.Replace(name, "_", "-", -1)
// 	for _, v := range list {
// 		p := strings.Split(v, "/")[0]
// 		if k8sname == p {
// 			return true
// 		}
// 	}
// 	return false
// }

// // https://gitlab.com/gitlab-org/gitlab-ce/issues/51508, version 11.2 only
// func userIsInGroup(g *gitlab.Group, user string) bool {
// 	_, _, err := client().Groups.ListAllGroupMembers(g.ID, &gitlab.ListGroupMembersOptions{
// 		Query: &user,
// 	})
// 	fmt.Println(err)
// 	return err == nil
// }

func GetProjects(token string) (c *gitlab.Client, pss []*gitlab.Project, err error) {

	// for all group projects
	c, gs, err := GetGroups(token)
	if err != nil {
		log.Println("getgroups err", err)
		return
	}

	var wg sync.WaitGroup

	queue := make(chan []*gitlab.Project, 100) // estimate value
	wg.Add(len(gs))

	// a := gitlab.PrivateVisibility
	// access := gitlab.NoPermissions
	// list := gitlab.ListOptions{PerPage: -1}
	for _, g := range gs {
		go func(g *gitlab.Group) {
			// fmt.Println("got group", g.Path)
			defer wg.Done()

			// ps, _, e := client().Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
			// 	Visibility: &a,
			// })
			// if e != nil {
			// 	return
			// }
			// println("len for ", len(ps), g.Path)
			// printproject(ps, "ham")
			// queue <- ps
			// ps, _, e = client().Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
			// 	MinAccessLevel: &access,
			// })
			// if e != nil {
			// 	return
			// }
			// println("len for ", len(ps), g.Path)
			// printproject(ps, "ham")
			// queue <- ps

			// ps, _, e := c.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
			// 	ListOptions: list,
			// 	Visibility:  &a,
			// })
			// if e != nil {
			// 	return
			// }
			// // spew.Dump("resp", resp)
			// println("len for ", len(ps), g.Path)
			// printproject(ps, "agent")
			// queue <- ps
			// ps, _, e = c.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
			// 	ListOptions:    list,
			// 	MinAccessLevel: &access,
			// })
			// if e != nil {
			// 	return
			// }
			// // spew.Dump("resp", resp)
			// println("len for ", len(ps), g.Path)
			// printproject(ps, "agent")
			// queue <- ps

			listProjects(c, g, queue)
		}(g)
	}

	go func() {
		defer wg.Done()
		for ps := range queue {
			pss = append(pss, ps[:]...)
		}
	}()

	wg.Add(1)
	go func() {
		// for all personal projects inclusion
		// ps, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{
		// 	Visibility: &a,
		// })
		// if err != nil {
		// 	log.Println("listprojects err", err)
		// 	return
		// }
		// // println("len for personal", len(ps))
		// // printproject(ps, "agent")
		// queue <- ps
		listPersonalProjects(c, queue)

		wg.Done()
	}()
	wg.Wait()

	if len(pss) == 0 {
		err = fmt.Errorf("there's no any projects")
		log.Println(err)
		return
	}
	// fmt.Println("len", len(pss))
	return
}

func listPersonalProjects(c *gitlab.Client, queue chan []*gitlab.Project) {
	a := gitlab.PrivateVisibility
	// access := gitlab.NoPermissions
	list := gitlab.ListOptions{Page: 1, PerPage: 1000} //perpage doesn't work

	for i := 1; ; i++ {
		ps, resp, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{
			ListOptions: list,
			Visibility:  &a,
		})
		if err != nil {
			log.Println("listprojects err", err)
			return
		}

		// println("len for personal", len(ps))
		// printproject(ps, "agent")
		queue <- ps

		if resp.TotalPages < i {
			break
		}
		list.Page = i
	}
}

func listProjects(c *gitlab.Client, g *gitlab.Group, queue chan []*gitlab.Project) {
	a := gitlab.PrivateVisibility
	access := gitlab.DeveloperPermissions
	list := gitlab.ListOptions{Page: 1, PerPage: 1000} //perpage doesn't work

	// ps, _, e := client().Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
	// 	Visibility: &a,
	// })
	// if e != nil {
	// 	return
	// }
	// println("len for ", len(ps), g.Path)
	// printproject(ps, "ham")
	// queue <- ps
	// ps, _, e = client().Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
	// 	MinAccessLevel: &access,
	// })
	// if e != nil {
	// 	return
	// }
	// println("len for ", len(ps), g.Path)
	// printproject(ps, "ham")
	// queue <- ps
	for i := 1; ; i++ {
		ps, resp, e := c.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
			ListOptions:    list,
			Visibility:     &a, //this cause need second list
			MinAccessLevel: &access,
		})
		if e != nil {
			return
		}

		// spew.Dump("resp", resp)
		// println("len for ", len(ps), g.Path)
		// printproject(ps, "agent")
		queue <- ps

		ps, _, e = c.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
			ListOptions:    list,
			MinAccessLevel: &access,
		})
		if e != nil {
			return
		}
		// spew.Dump("resp", resp)
		// println("len for ", len(ps), g.Path)
		// printproject(ps, "agent")
		queue <- ps

		if resp.TotalPages < i {
			break
		}
		list.Page = i
	}
}

func printproject(ps []*gitlab.Project, name string) {
	for _, v := range ps {
		if strings.Contains(v.WebURL, name) {
			fmt.Println("ha", v.WebURL)
		}
	}
}

func GetProjectsOld(token string) (c *gitlab.Client, pss []*gitlab.Project, err error) {
	// for all group projects
	c, gs, err := GetGroups(token)
	if err != nil {
		log.Println("getgroups err", err)
		return
	}
	a := gitlab.PrivateVisibility
	for _, g := range gs {
		ps, _, e := c.Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
			Visibility: &a,
		})
		if e != nil {
			continue
		}
		pss = append(pss, ps[:]...)
	}

	// for all personal projects inclusion
	ps, _, err := c.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	if err != nil {
		log.Println("listprojects err", err)
		return
	}
	pss = append(pss, ps[:]...)

	if len(pss) == 0 {
		err = fmt.Errorf("there's no any projects")
		log.Println(err)
		return
	}
	// for _, v := range pss {
	// 	if strings.Contains(v.WebURL, "yunwei/worktile") || strings.Contains(v.WebURL, "yunwei/trx") {
	// 		spew.Dump("v:", v)
	// 		// fmt.Println(v.WebURL, v.RequestAccessEnabled)
	// 	}
	// 	fmt.Println(v.WebURL, v.RequestAccessEnabled)
	// }
	return
}

func GetProjectLists(token string) (admin bool, projects []string, err error) {
	isadmin, e := IsAdmin(token)
	if err != nil {
		err = fmt.Errorf("check admin error %v", e)
		return
	}
	if isadmin {
		// // filter list to reduce project searching time
		// list, e := listpods()
		// if e != nil {
		// 	err = fmt.Errorf("walk error %v", e)
		// 	return
		// }
		admin = true
		return
	}

	_, pss, err := GetProjects(token)
	if err != nil {
		log.Println("getprojects err", err)
		return
	}
	for _, p := range pss {
		// fmt.Println("--", p.PathWithNamespace)
		// spew.Dump("p", p)
		// url := strings.Split(p.WebURL, "/")
		// if len(url) != 5 {
		// 	log.Println("get project list warn: bad format %v", p.WebURL)
		// 	continue
		// }
		// git := fmt.Sprintf("%v/%v", url[3], url[4])
		git := strings.Replace(p.PathWithNamespace, " ", "", -1) //remove empty space
		projects = append(projects, git)
	}
	return false, unique(projects), nil
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// projectPath is org/repo
func GetGitProject(projectPath string) (project *gitlab.Project, err error) {
	project, _, err = client().Projects.GetProject(projectPath)
	return
}

func IsAdmin(token string) (isadmin bool, err error) {
	u, err := GetUser(token)
	if err != nil {
		err = fmt.Errorf("get user err: %v", err)
		return
	}
	if u.IsAdmin {
		// envs = append(envs, EnvPreOnline, EnvTest)
		isadmin = true
		return
	}
	return
}

func CheckPerm(projectPath, token string) (envs []string, err error) {
	u, err := GetUser(token)
	if err != nil {
		err = fmt.Errorf("get user err: %v", err)
		return
	}
	if u.IsAdmin {
		envs = append(envs, EnvOnline, EnvPreOnline, EnvTest)
		return
	}
	// check permissions
	project, err := GetGitProject(projectPath)
	if err != nil {
		err = fmt.Errorf("get git project err: %v", err)
		return
	}

	al, err := getAccessLevel(project, u.ID)
	if err != nil {
		err = fmt.Errorf("get access level err: %v", err)
		return
	}
	envs = getAllowedEnv(al)
	return
}

func getAccessLevel(project *gitlab.Project, userid int) (accessLevel gitlab.AccessLevelValue, err error) {
	var groupAccessLevel, projectAccessLevel gitlab.AccessLevelValue
	groupMember, _, e := client().GroupMembers.GetGroupMember(project.Namespace.ID, userid)
	// fmt.Println("groupmembers err", err)
	if e == nil {
		groupAccessLevel = groupMember.AccessLevel
	}
	projectMember, _, e := client().ProjectMembers.GetProjectMember(project.ID, userid)
	//fmt.Println("projectmember err", err, project, userid)
	if e == nil {
		projectAccessLevel = projectMember.AccessLevel
	}
	if groupAccessLevel > projectAccessLevel {
		accessLevel = groupAccessLevel
	} else {
		accessLevel = projectAccessLevel
	}
	return
}

func getAllowedEnv(accessLevel gitlab.AccessLevelValue) (envs []string) {
	if accessLevel >= gitlab.DeveloperPermissions {
		envs = append(envs, EnvPreOnline, EnvTest)
	}
	if accessLevel >= gitlab.MasterPermissions {
		envs = append(envs, EnvOnline)
	}
	return
}

func isEnvOk(src, env string, envs []string) bool {
	t := filepath.Join(src, env)
	if f, err := os.Stat(t); err != nil || !f.IsDir() {
		// fmt.Printf("target %v is not exist\n", t)
		return false
	}
	for _, v := range envs {
		if env == v {
			return true
		}
	}
	return false
}

func getBinds(src, dst string, envs []string) (okenvs []string, binds map[string]string, err error) {
	binds = map[string]string{}
	var nologs bool
	for _, env := range envs {
		// fmt.Println("check env: ", env)
		if isEnvOk(src, env, envs) {
			t := filepath.Join(dst, "logs", env)
			// if env == EnvOnline {
			// 	t = filepath.Join(dst, "logs") //make online directly bind to logs
			// }
			// fmt.Println("append bind: ", env)
			binds[filepath.Join(src, env)] = t
			okenvs = append(okenvs, env)
			nologs = false
		}
	}
	if nologs {
		err = fmt.Errorf("no any logs been found")
		return
	}
	return
}

// func filterEnvs(envs []string, path string) (err error) {
// 	dirs, err := getDirs(path)
// 	if err != nil {
// 		return
// 	}
// 	for _,v:=range envs{
// 		for _,v:=range dirs{

// 		}
// 	}
// }

// func getDirs(path string) (dirs []string, err error) {
// 	files, err := ioutil.ReadDir(path)
// 	if err != nil {
// 		return
// 	}
// 	for _, file := range files {
// 		if !file.IsDir() {
// 			continue
// 		}
// 		dirs = append(dirs, file.Name())
// 	}
// 	if len(dirs) == 0 {
// 		return nil, fmt.Errorf("no log dirs found")
// 	}
// 	return
// }

// func GetGitGroup(group string) (g *gitlab.Group, err error) {
// 	ps, _, err := client().Groups.ListGroups(&gitlab.ListGroupsOptions{
// 		// Search: &group,
// 	})
// 	for _, v := range ps {
// 		fmt.Println(v.WebURL)
// 	}
// 	if err != nil {
// 		err = fmt.Errorf("%v", strings.Split(err.Error(), "\n")[0])
// 		return
// 	}
// 	if len(ps) == 0 {
// 		err = fmt.Errorf("group: %v not found", group)
// 		return
// 	}
// 	g = ps[0]
// 	return
// }

// func GetGitProject(git string) (p *gitlab.Project, err error) {
// 	gr := strings.Split(git, "/")
// 	if len(gr) != 2 {
// 		err = fmt.Errorf("git: %v invalid format, eg: group/project ", git)
// 		return
// 	}
// 	group, repo := gr[0], gr[1]

// 	var g *gitlab.Group
// 	g, err = GetGitGroup(group)
// 	if err != nil {
// 		return
// 	}
// 	ps, _, e := client().Groups.ListGroupProjects(g.ID, &gitlab.ListGroupProjectsOptions{
// 		// Membership: &a,
// 		Search: &repo,
// 	})
// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	if len(ps) == 0 {
// 		err = fmt.Errorf("repo: %v not found", repo)
// 		return
// 	}
// 	p = ps[0]
// 	return
// }

// func ListUsers() ([]*gitlab.User, error) {
// 	u, _, err := client().Users.ListUsers(&gitlab.ListUsersOptions{})
// 	return u, err
// }

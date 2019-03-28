package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	gitlab "github.com/xanzy/go-gitlab"
)

var (
	GitlabAccessToken = flag.String("gitlabtoken", "", "gitlab access token")
	GitlabEndpoint    = flag.String("gitlaburl", "http://g.haodai.net", "gitlab base url")
)

const (
	EnvOnline    = "online"
	EnvPreOnline = "pre-online"
	EnvTest      = "test"
)

// var client *gitlab.Client

func client() *gitlab.Client {
	client := gitlab.NewClient(http.DefaultClient, *GitlabAccessToken)
	client.SetBaseURL(*GitlabEndpoint)
	return client
}

// projectPath is org/repo
func GetGitProject(projectPath string) (project *gitlab.Project, err error) {
	project, _, err = client().Projects.GetProject(projectPath)
	return
}

func CheckPerm(projectPath string, userid int) (envs []string, err error) {
	// check permissions
	project, err := GetGitProject(projectPath)
	if err != nil {
		return
	}
	al, err := getAccessLevel(project, userid)
	if err != nil {
		return
	}
	envs = getAllowedEnv(al)
	return
}

func getAccessLevel(project *gitlab.Project, userid int) (accessLevel gitlab.AccessLevelValue, err error) {
	var groupAccessLevel, projectAccessLevel gitlab.AccessLevelValue
	groupMember, _, err := client().GroupMembers.GetGroupMember(project.Namespace.ID, userid)
	if err == nil {
		groupAccessLevel = groupMember.AccessLevel
	}
	projectMember, _, err := client().ProjectMembers.GetProjectMember(project.ID, userid)
	if err == nil {
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
		fmt.Printf("target %v is not exist", t)
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

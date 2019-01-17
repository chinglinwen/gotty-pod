package container

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

type Container struct {
	WorkDir  string
	Arg      []string
	Hostname string

	CGroupName string
	Src        string //log path
	Dst        string //mount path
	Rootfs     string //alpine rootfs
}

// // go run main.go run <cmd> <args>
// func main() {
// 	fmt.Println("starting...")
// 	switch os.Args[1] {
// 	case "run":
// 		run()
// 	case "child":
// 		child()
// 	default:
// 		panic("help")
// 	}
// }

// func Run() {
// 	fmt.Printf("Running %v \n", os.Args[2:])

// 	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
// 	cmd.Stdin = os.Stdin
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	cmd.SysProcAttr = &syscall.SysProcAttr{
// 		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
// 		Unshareflags: syscall.CLONE_NEWNS,
// 		//Credential: &Credential{Uid: uid, Gid: gid}
// 	}

// 	must(cmd.Run())
// }

// run in container way
func (c *Container) Run() error {
	err := c.Validate()
	if err != nil {
		return err
	}

	_, err = Mount(c.Src, c.Rootfs, c.Dst)
	if err != nil {
		return err
	}

	//fmt.Printf("Running %v \n", os.Args[2:])
	cmd, err := c.CreateCMD()
	if err != nil {
		return err
	}
	return cmd.Run()
	// ptmx, err := pty.Start(cmd)
	// if err != nil {
	// 	return fmt.Errorf("start cmd error: %v", err)
	// }
	// defer func() { _ = ptmx.Close() }()

	// ch := make(chan os.Signal, 1)
	// signal.Notify(ch, syscall.SIGWINCH)
	// go func() {
	// 	for range ch {
	// 		if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
	// 			fmt.Printf("error resizing pty: %s", err)
	// 		}
	// 	}
	// }()
	// ch <- syscall.SIGWINCH // Initial resize.

	// // Set stdin in raw mode.
	// oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// // Copy stdin to the pty and the pty to stdout.
	// go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	// _, _ = io.Copy(os.Stdout, ptmx)

	_, err = UnMount(c.Dst)
	if err != nil {
		return err
	}

	return nil
}

func (c *Container) CreateCMD() (cmd *exec.Cmd, err error) {

	// cgroup setting
	//createcgroup(c.CGroupName)

	//cmd = exec.Command(c.Arg[0], c.Arg[1:]...)
	cmd = exec.Command(c.Arg[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//cmd.Dir = c.WorkDir

	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	// 	Unshareflags: syscall.CLONE_NEWNS,
	// 	//Credential: &Credential{Uid: uid, Gid: gid}
	// }

	err = syscall.Sethostname([]byte(c.Hostname))
	if err != nil {
		return nil, fmt.Errorf("sethostname: %v error, err: %v", c.Hostname, err)
	}
	err = syscall.Chroot(c.Dst)
	if err != nil {
		return nil, fmt.Errorf("chroot %v error, err: %v", c.Dst, err)
	}
	err = os.Chdir("/")
	if err != nil {
		return nil, fmt.Errorf("change dir to / error, err: %v", err)
	}
	err = syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		return nil, fmt.Errorf("mount proc error, err: %v", err)
	}
	return
}

func Mount(src, rootfs, dst string) (out []byte, err error) {
	err = pathCheck(src, rootfs, dst)
	if err != nil {
		return nil, err
	}
	opts := fmt.Sprintf("ro,lowerdir=%v:%v", src, rootfs)
	cmd := exec.Command("sudo", "mount", "-t", "overlay", "overlay", "-o", opts, dst)
	return cmd.CombinedOutput()
}

func (c *Container) Validate() error {
	if len(c.Arg) == 0 {
		return errors.New("no cmd to run")
	}
	if c.Arg[0] == "" {
		return errors.New("empty cmd to run")
	}

	if c.CGroupName == "" {
		return errors.New("cgroup name is empty")
	}
	if c.Src == "" {
		return errors.New("src name is empty")
	}
	if c.Dst == "" {
		return errors.New("dst name is empty")
	}
	if c.Rootfs == "" {
		return errors.New("rootfs name is empty")
	}
	return nil
}

func pathCheck(src, rootfs, dst string) error {
	_, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("%v does not exist", src)
	}
	_, err = os.Stat(rootfs)
	if err != nil {
		return fmt.Errorf("%v does not exist", rootfs)
	}
	os.MkdirAll(dst, 0755)
	_, err = os.Stat(dst)
	if err != nil {
		return fmt.Errorf("%v does not exist", dst)
	}
	return nil
}

func UnMount(dst string) (out []byte, err error) {
	cmd := exec.Command("sudo", "umount", "-f", dst)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(dst)
	return
}

//sudo mount -t overlay overlay -o ro,lowerdir=/home/wen/log/ /home/wen/alpine/mnt
//sudo mount -t overlay overlay -o ro,lowerdir=/home/wen/log/:/home/wen/alpine/ /home/wen/alpine/mnt
//sudo umount /home/wen/target_root_fs

func createcgroup(name string) {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, name), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, name+"/pids.max"), []byte("20"), 0700))
	// Removes the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, name+"/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, name+"/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

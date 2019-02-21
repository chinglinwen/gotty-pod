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

const (
	EnvOnline    = "online"
	EnvPreOnline = "pre-online"
	EnvTest      = "test"
)

type Container struct {
	WorkDir    string
	Arg        []string
	Hostname   string
	CGroupName string
	Binds      map[string]string //multiple log binds
	Dst        string            //mount path
	// BindDst    string            //bind dst
	Rootfs string //alpine rootfs
	Envs   []string

	// Pr *io.PipeReader
	// Pw *io.PipeWriter
}

// run in container way
func (c *Container) Run() error {
	err := c.Validate()
	if err != nil {
		return err
	}

	// mount rootfs
	_, err = Mount(c.Rootfs, c.Dst)
	if err != nil {
		return err
	}

	// _, err = MountDev(c.Dst)
	// if err != nil {
	// 	return err
	// }

	// // try mount online into logs, need create other directory for bind, which is readonly
	// for src, dst := range c.Binds {
	// 	if filepath.Base(src) == EnvOnline {
	// 		// bind online first
	// 		_, err = MountBind(src, dst)
	// 		if err != nil {
	// 			fmt.Printf("try to bind %v:%v, err: %v\n", src, dst, err)
	// 			return err
	// 		}
	// 	}
	// }

	for src, dst := range c.Binds {
		// if filepath.Base(src) == EnvOnline {
		// 	// skip online
		// 	continue
		// }
		// bind logs directories
		_, err = MountBind(src, dst)
		if err != nil {
			fmt.Printf("try to bind %v:%v, err: %v\n", src, dst, err)
			return err
		}
	}

	// // no need to unmount
	// defer func() {
	// 	_, err = UnMount(c.Dst)
	// 	if err != nil {
	// 		fmt.Printf("umount %v err: %v\n", c.Dst, err)
	// 	}
	// 	_, err = UnMount(c.BindDst)
	// 	if err != nil {
	// 		fmt.Printf("umount %v err: %v\n", c.BindDst, err)
	// 	}

	// 	_, err = UnMountDev(c.Dst)
	// 	if err != nil {
	// 		fmt.Printf("umount %v/dev err: %v\n", c.Dst, err)
	// 	}
	// }()

	//fmt.Printf("Running %v \n", os.Args[2:])
	cmd, err := c.CreateCMD()
	if err != nil {
		return err
	}
	cmd.Run()

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

	return nil
}

func (c *Container) CreateCMD() (cmd *exec.Cmd, err error) {

	// cgroup setting
	//createcgroup(c.CGroupName)

	// var b bytes.Buffer
	// // go func() {
	// // 	time.Sleep(3 * time.Second)
	// // 	os.Stdin.WriteString("pwd\n")
	// // }()
	// b.WriteString("pwd\n")

	//cmd = exec.Command(c.Arg[0], c.Arg[1:]...)

	cmd = exec.Command(c.Arg[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	// cmd.Stdin = c.Pr
	// cmd.Stdout = c.Pw
	cmd.Stderr = os.Stderr
	cmd.Dir = c.WorkDir

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// Unshareflags: syscall.CLONE_NEWNS,
		Credential: &syscall.Credential{Uid: 65534, Gid: 65534},
	}

	// cmd.SysProcAttr = &syscall.SysProcAttr{
	// 	Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	// 	Unshareflags: syscall.CLONE_NEWNS,
	// }
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
	// err = syscall.Mount("proc", "proc", "proc", 0, "")
	// if err != nil {
	// 	return nil, fmt.Errorf("mount proc error, err: %v", err)
	// }
	return
}

func Mount(rootfs, dst string) (out []byte, err error) {
	err = pathCheck(rootfs, dst)
	if err != nil {
		return nil, err
	}
	opts := fmt.Sprintf("rw,lowerdir=%v:%v", dst, rootfs) // we fake dst as lowerdir, to make cmd ok
	//fmt.Println("sudo", "mount", "-t", "overlay", "overlay", "-o", opts, dst)
	cmd := exec.Command("sudo", "mount", "-t", "overlay", "overlay", "-o", opts, dst)
	return cmd.CombinedOutput()
}

// bind logs
func MountBind(src, dst string) (out []byte, err error) {
	err = pathCheck(src, dst)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command("sudo", "mount", "--bind", "-o", "rw", src, dst)
	return cmd.CombinedOutput()
}

// bind logs
func MountDev(dst string) (out []byte, err error) {
	// err = pathCheck(src, dst)
	// if err != nil {
	// 	return nil, err
	// }
	cmd := exec.Command("sudo", "mount", "--bind", "/dev/", dst+"/dev")
	return cmd.CombinedOutput()
}

// bind logs
func UnMountDev(dst string) (out []byte, err error) {
	// err = pathCheck(src, dst)
	// if err != nil {
	// 	return nil, err
	// }
	cmd := exec.Command("sudo", "umount", "-f", dst+"/dev")
	return cmd.CombinedOutput()
}

func (c *Container) Validate() error {
	if len(c.Arg) == 0 {
		return errors.New("no cmd to run")
	}
	if c.Arg[0] == "" {
		return errors.New("empty cmd to run")
	}

	// if c.CGroupName == "" {
	// 	return errors.New("cgroup name is empty")
	// }
	if c.Binds == nil {
		return errors.New("no logs to bind")
	}
	if c.Dst == "" {
		return errors.New("dst name is empty")
	}
	if c.Rootfs == "" {
		return errors.New("rootfs name is empty")
	}
	return nil
}

func pathCheck(src, dst string) error {
	_, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("src %v does not exist", src)
	}
	err = os.MkdirAll(dst, 0755)
	if err != nil {
		fmt.Printf("create dst err: %v\n", err)
	}
	_, err = os.Stat(dst)
	if err != nil {
		fmt.Printf("dst %v does not exist\n", dst)
		//return fmt.Errorf("dst %v does not exist", dst)
		return nil
	}
	return nil
}

func UnMount(dst string) (out []byte, err error) {
	cmd := exec.Command("umount", "-f", dst)
	out, err = cmd.CombinedOutput()
	if err != nil {
		cmd := exec.Command("umount", "-f", dst)
		return cmd.CombinedOutput()
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

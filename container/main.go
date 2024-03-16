package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"log"
)

func main() {
	switch os.Args[1]{
	case "run":
		run()
	case "child":
		child()
	default:
		fmt.Println("COMMAND not found")
	}
}

func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	logerror(cmd.Run())
}

func child() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/home/ubuntu")
	os.Chdir("/home/ubuntu")
	syscall.Mount("proc", "proc", "proc", 0, "")
	syscall.Mount("thing", "mytemp", "tmpfs", 0, "")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func logerror(e error) {
	if e != nil {
		log.Fatalf("Failed with %s\n", e)
	}
}


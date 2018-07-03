package build

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chuckleheads/pkg-export-oci/chmod"
	"github.com/go-cmd/cmd"
)

type BuildSpec struct {
	Hab                 string
	HabLauncher         string
	HabSup              string
	URL                 string
	Channel             string
	BasePackagesURL     string
	BasePackagesChannel string
}

type Naming struct {
	CustomImageName   string
	LatestTag         bool
	VersionTag        bool
	VersionReleaseTag bool
	CustomTag         string
}

func Build(fsroot string, pkg string) {
	installBasePkgs(fsroot)
	installUserPkgs(fsroot, pkg)
	chmod.ChmodR(filepath.Join(fsroot, "hab"))
	binlink(fsroot)

	fmt.Printf("I'm a service?: %v", isAService(fsroot, pkg))
}

func install(fsroot string, pkg string) {
	runCommand(fsroot, "hab", "pkg", "install", pkg)
}

func installBasePkgs(fsroot string) {
	basePkgs := []string{"core/hab", "core/hab-sup", "core/hab-launcher", "core/busybox", "core/cacerts"}
	for _, pkg := range basePkgs {
		install(fsroot, pkg)
	}
}

func installUserPkgs(fsroot string, pkg string) {
	install(fsroot, pkg)
}

func binlink(fsroot string) {
	runCommand(fsroot, "hab", "pkg", "binlink", "core/busybox", "busybox")
	runCommand(fsroot, "hab", "pkg", "binlink", "core/busybox", "sh")
	runCommand(fsroot, "hab", "pkg", "binlink", "core/hab", "hab")
}

func runCommand(fsroot string, args ...string) cmd.Status {
	cmdOptions := cmd.Options{
		Buffered:  true,
		Streaming: true,
	}

	name, args := args[0], args[1:]

	// Create Cmd with options
	habCmd := cmd.NewCmdOptions(cmdOptions, name, args...)
	habCmd.Env = []string{fmt.Sprintf("FS_ROOT=%s", fsroot), "TERM=xterm-256color"}
	// Print STDOUT and STDERR lines streaming from Cmd
	go func() {
		for {
			select {
			case line := <-habCmd.Stdout:
				fmt.Println(line)
			case line := <-habCmd.Stderr:
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	// Run and wait for Cmd to return, discard Status
	finalStatus := <-habCmd.Start()

	// Cmd has finished but wait for goroutine to print all lines
	for len(habCmd.Stdout) > 0 || len(habCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}
	return finalStatus
}

func isAService(fsroot string, ident string) bool {
	status := runCommand(fsroot, "hab", "pkg", "path", ident)
	if status.Error != nil {
		panic(status.Error)
	}
	_, err := os.Stat(filepath.Join(status.Stdout[0], "SVC_USER"))
	// If the `SVC_USER` file doesn't exist we aren't a service
	return !os.IsNotExist(err)
}

package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chuckleheads/habitat/components/hab/pkg/ui"

	"github.com/chuckleheads/pkg-export-oci/account"
	"github.com/chuckleheads/pkg-export-oci/chmod"
	"github.com/chuckleheads/pkg-export-oci/runc"
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
	Entrypoint          string
}

type Naming struct {
	CustomImageName   string
	LatestTag         bool
	VersionTag        bool
	VersionReleaseTag bool
	CustomTag         string
}

func (b *BuildSpec) Build(fsroot string, pkg string) {
	installUserPkgs(fsroot, pkg)
	totesAService := isAService(fsroot, pkg)
	if !totesAService && b.Entrypoint == "" {
		fmt.Printf("No service found and Entrypoint is empty")
		os.Exit(1)
	}
	installBasePkgs(fsroot, totesAService)
	chmod.ChmodR(filepath.Join(fsroot, "hab"))
	binlink(fsroot)

	if totesAService {
		account := serviceAccounts(fsroot, pkg)
		account.Write(fsroot)
	}

	err := runc.Config(fsroot, pkg, b.Entrypoint, totesAService)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func install(fsroot string, pkg string) {
	runHabCommand(fsroot, "pkg", "install", pkg)
}

func installBasePkgs(fsroot string, serviceIs bool) {
	basePkgs := []string{"core/hab", "core/busybox", "core/cacerts"}
	// We don't install the sup and launcher if it's not needed.
	// If pkg is not a service, then we consider it not needed.
	if serviceIs {
		basePkgs = append(basePkgs, "core/hab-sup")
		basePkgs = append(basePkgs, "core/hab-launcher")
	}
	for _, pkg := range basePkgs {
		install(fsroot, pkg)
	}
}

func installUserPkgs(fsroot string, pkg string) {
	install(fsroot, pkg)
}

func binlink(fsroot string) {
	runHabCommand(fsroot, "pkg", "binlink", "core/busybox", "busybox")
	runHabCommand(fsroot, "pkg", "binlink", "core/busybox", "sh")
	runHabCommand(fsroot, "pkg", "binlink", "core/hab", "hab")
}

func runHabCommand(fsroot string, args ...string) cmd.Status {
	cmdOptions := cmd.Options{
		Buffered:  true,
		Streaming: true,
	}

	// Create Cmd with options
	habCmd := cmd.NewCmdOptions(cmdOptions, "hab", args...)
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
	status := runHabCommand(fsroot, "pkg", "path", ident)
	if status.Error != nil {
		ui.Fatal(fmt.Sprintf("%d", status.Exit))
		os.Exit(1)
	}
	_, err := os.Stat(filepath.Join(status.Stdout[0], "SVC_USER"))
	// If the `SVC_USER` file doesn't exist we aren't a service
	return !os.IsNotExist(err)
}

func serviceAccounts(fsroot string, ident string) account.Account {
	status := runHabCommand(fsroot, "pkg", "path", ident)
	if status.Error != nil {
		ui.Fatal(fmt.Sprintf("%d", status.Exit))
		os.Exit(1)
	}
	user, err := ioutil.ReadFile(filepath.Join(status.Stdout[0], "SVC_USER"))
	if err != nil {
		ui.Fatal(err.Error())
		os.Exit(1)
	}

	group, err := ioutil.ReadFile(filepath.Join(status.Stdout[0], "SVC_GROUP"))
	if err != nil {
		ui.Fatal(err.Error())
		os.Exit(1)
	}

	return account.New(strings.TrimSpace(string(user[:])), strings.TrimSpace(string(group[:])))
}

func renderPasswd(rootfs string) {

}

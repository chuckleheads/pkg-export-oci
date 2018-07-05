package rootfs

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

const RESOLVCONF = "defaults/etc/resolv.conf"
const NSSWITCH = "defaults/etc/nsswitch.conf"

func Create() (dir string) {
	dir = createTempFS()
	// force umask to not mask so none of the local build systems umasking bleeds into our container
	syscall.Umask(0)
	os.MkdirAll(filepath.Join(dir, "bin"), 0755)
	os.MkdirAll(filepath.Join(dir, "etc"), 0755)
	os.MkdirAll(filepath.Join(dir, "root"), 0755)
	os.MkdirAll(filepath.Join(dir, "tmp"), 0777)
	os.MkdirAll(filepath.Join(dir, "var/tmp"), 0777)

	resolv, err := ioutil.ReadFile(RESOLVCONF)
	check(err)
	ioutil.WriteFile(filepath.Join(dir, "etc", "resolv.conf"), resolv, 700)
	nsswitch, err := ioutil.ReadFile(NSSWITCH)
	check(err)
	ioutil.WriteFile(filepath.Join(dir, "etc", "nsswitch"), nsswitch, 700)
	// Force umask to be a sane default for all future installs
	syscall.Umask(022)
	return dir
}

func createTempFS() (dir string) {
	dir, err := ioutil.TempDir("", "hab-exporter-temp")
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, "rootfs")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

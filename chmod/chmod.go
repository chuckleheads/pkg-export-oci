package chmod

import (
	"os"
	"path/filepath"
)

func ChmodR(path string) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			current := info.Mode().Perm()
			new := gEqualsU(current)
			err = os.Chmod(name, new)
		}
		return err
	})
}

func gEqualsU(perms os.FileMode) os.FileMode {
	other := perms & 07
	user := (perms >> 6) & 07
	cleared := (perms >> 9) << 9
	newPerms := (((user << 3) | user) << 3) | other
	return cleared | newPerms
}

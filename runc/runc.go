package runc

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/opencontainers/runc/libcontainer/specconv"
)

func GenRuncConfig(fsroot string, isAService bool) error {
	genRuncConfig := specconv.Example()
	genRuncConfig.Process.Env = []string{"PATH=/bin"}

	if isAService {
		// run hab sup run <origin>/<pkg>

	} else {
		// run hab pkg exec <origin>/<pkg>

	}

	data, err := json.MarshalIndent(genRuncConfig, "", "\t")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(fsroot, "..", "config.json"), data, 0666); err != nil {
		return err
	}
	return nil
}

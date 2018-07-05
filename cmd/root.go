// Copyright © 2018 Elliott Davis
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chuckleheads/pkg-export-oci/build"
	"github.com/chuckleheads/pkg-export-oci/rootfs"
	"github.com/spf13/cobra"
)

var b build.BuildSpec

var rootCmd = &cobra.Command{
	Use:   "hab-oci-exporter [flags] <PKG_IDENT_OR_ARTIFACT>",
	Short: "Exporter for runc packages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tmpDir := rootfs.Create()
		defer os.RemoveAll(tmpDir)
		ident := args[0]
		b.Build(tmpDir, ident)
		fIdent := strings.Replace(ident, "/", "-", -1)
		fmt.Println("Creating archive...")
		komand := exec.Command(fmt.Sprintf("tar -cvjSf %s.tar.bz2 %s %s", fIdent, tmpDir, filepath.Join(tmpDir, "..", "config.json")))
		err := komand.Run()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&b.Entrypoint, "entrypoint", "e", "", "Specify an optional default entrypoint for the service. This will cause the container to run in a one-off mode")
}

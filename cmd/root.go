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

	"github.com/chuckleheads/hab-oci-exporter/build"
	"github.com/chuckleheads/hab-oci-exporter/rootfs"
	"github.com/spf13/cobra"
)

var b build.BuildSpec
var n build.Naming

var rootCmd = &cobra.Command{
	Use:   "hab-oci-exporter [flags] <PKG_IDENT_OR_ARTIFACT>",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tmpDir := rootfs.Create()
		// defer os.RemoveAll(tmpDir)
		build.Build(tmpDir, args[0])
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("tag-latest", "", true, "Tag image with: \"latest\"")
	rootCmd.Flags().BoolP("tag-version", "", true, "Tag image with: \"{{pkg_version}}\"")
	rootCmd.Flags().BoolP("tag-latest-release", "", true, "Tag image with: \"{{pkg_version}}-{{pkg_release}}\"")
	rootCmd.Flags().StringVarP(&b.BasePackagesURL, "base-pkgs-url", "", "https://bldr.habitat.sh", "Install base packages from Builder at the specified URL")
	rootCmd.Flags().StringVarP(&b.BasePackagesChannel, "base-pkgs-channel", "", "stable", "Install base packages from the specified release channel")
	rootCmd.Flags().StringVarP(&b.URL, "url", "u", "https://bldr.habitat.sh", "Install packages from Builder at the specified URL")
	rootCmd.Flags().StringVarP(&b.Channel, "channel", "c", "stable", "Install packages from the specified release channel")
	rootCmd.Flags().StringVarP(&b.Hab, "hab-pkg", "", "core/hab", "Habitat CLI package identifier (ex: acme/redis) or filepath to a Habitat artifact (ex: acme-redis-3.0.7-21120102031201-x86_64-linux.hart) to install")
	rootCmd.Flags().StringVarP(&b.HabLauncher, "launcher-pkg", "", "core/hab-launcher", "Launcher package identifier (ex: core/hab-launcher) or filepath to a Habitat artifact (ex: core-hab-launcher-6083-20171101045646-x86_64-linux.hart) to install")
	rootCmd.Flags().StringVarP(&b.HabSup, "sup-pkg", "", "core/hab-sup", "Supervisor package identifier (ex: core/hab-sup) or filepath to a Habitat artifact (ex: core-hab-sup-0.39.1-20171118011657-x86_64-linux.hart) to install")
	rootCmd.Flags().StringVarP(&n.CustomImageName, "img-name", "i", "{{pkg_origin}}/{{pkg_name}}", "Image name")
	rootCmd.Flags().StringVarP(&n.CustomTag, "tag-custom", "", "", " Tag image with additional custom tag")
}

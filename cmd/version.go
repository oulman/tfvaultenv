/*
Copyright © 2024 Dylan Northrup

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tfVaultEnvVersion string

func init() {
	tfVaultEnvVersion = "0.5.1"
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version",
	Long:  "Show version of tfvaultenv command",
	Run: func(cmd *cobra.Command, args []string) {

		logrus.SetFormatter(&logrus.TextFormatter{
			DisableLevelTruncation: true,
			PadLevelText:           true,
			DisableTimestamp:       true,
		})

		fmt.Printf("tfvaultenv version %s\n", tfVaultEnvVersion)
	},
}

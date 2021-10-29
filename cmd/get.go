/*
Copyright © 2021 James Oulman

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
	"math/rand"
	"time"

	config "github.com/oulman/tfvaultenv/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const envPrefix = "TFVAULTENV"

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "queries Vault for secrets and outputs them in environment variables",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

		logrus.SetFormatter(&logrus.TextFormatter{
			DisableLevelTruncation: true,
			PadLevelText:           true,
			DisableTimestamp:       true,
		})

		// There is a random function for the HCL configuration.
		rand.Seed(time.Now().Unix())

		depth, err := cmd.Flags().GetInt("configdepth")
		if err != nil {
			logrus.Fatal(err)
		}

		configFileName, err := cmd.Flags().GetString("config")
		if err != nil {
			logrus.Fatal(err)
		}

		configFilePath, err := config.FindConfigFile(depth, configFileName)
		if err != nil {
			logrus.Fatal(err)
		}

		configParsed, err := config.ParseConfig(configFilePath)
		if err != nil {
			logrus.Fatal(err)
		}

		err = config.ProcessConfig(configParsed)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

// func getDefault(env string, defaultValue interface{}) interface{} {
// 	if env != "" {
// 		formattedEnv := fmt.Sprintf("%s_%s", envPrefix, env)
// 		value := os.Getenv(formattedEnv)
// 		if value != "" {
// 			return value
// 		}
// 	}
// 	return defaultValue
// }

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().String("config", ".tfvaultenv.config.hcl", "path to the configuration file")
	getCmd.Flags().Int("configdepth", 0, "number of parent directories to scan for the -config file")
}

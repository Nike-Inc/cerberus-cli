/*
 *  Copyright (c) 2019 Nike, Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package cmd

import (
	"fmt"
	"os"

	"cerberus-cli/client"
	"cerberus-cli/tool"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cerberus-cli",
	Short: "A CLI for Cerberus",
	Long:  `cerberus-cli is a CLI for Cerberus that can be used to perform basic tasks`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if Quiet {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
		}
		if client.Url == "" && client.Region == "" {
			return fmt.Errorf("Cerberus server url and AWS region not found. See [Global Flags] section below")
		} else if client.Url == "" {
			return fmt.Errorf("Cerberus server url not found. See [Global Flags] section below")
		} else if client.Region == "" {
			return fmt.Errorf("AWS Cerberus region not found. See [Global Flags] section below")
		}
		return nil
	},
}

var Quiet bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&client.Url, "url", "u", tool.GetEnvVariable(tool.EnvCerbUrl), fmt.Sprintf("Cerberus server url / set the %s env variable", tool.EnvCerbUrl))
	rootCmd.PersistentFlags().StringVarP(&client.Region, "region", "r", tool.GetEnvVariable(tool.EnvCerbRegion), fmt.Sprintf("AWS Cerberus region / set the %s env variable", tool.EnvCerbRegion))
	rootCmd.PersistentFlags().StringVarP(&client.Token, "token", "t", tool.GetEnvVariable(tool.EnvCerbToken), fmt.Sprintf("Cerberus token      / set the %s env variable", tool.EnvCerbToken))
	_ = rootCmd.PersistentFlags().MarkHidden("token")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "set this flag for quiet mode (no usage or error messages for valid commands)")
	rootCmd.Version = version
	rootCmd.InitDefaultVersionFlag()
	rootCmd.SetVersionTemplate(fmt.Sprintln("cerberus-cli " + version))
	rootCmd.SetErr(os.Stderr)
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetIn(os.Stdin)
}

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
	"cerberus-cli/tool"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version number of this tool",
	Long:  `Display the version number of this tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(tool.CliVersionMsg)
	},
}


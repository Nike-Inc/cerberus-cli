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
	"cerberus-cli/client"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var filelistCmd = &cobra.Command{
	Use:   "list <secure data path>",
	Short: "list files of a secure data path",
	Long:  `list files of a secure data path`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		path := args[0]
		category := strings.Split(path, "/")[0]

		autocomplete, err := cmd.Flags().GetBool("autocomplete")
		if err != nil {
			return err
		}

		response, err := cl.SecureFile().List(path)
		if err != nil {
			return err
		}

		if response == nil {
			return fmt.Errorf("path %s does not exist", path)
		}

		files := response.Summaries

		sep := "\n"
		if autocomplete {
			sep = "\\n"
		}

		filelist := ""
		for _, file := range files {
			filelist += category + "/" + file.Path + sep
		}
		fmt.Print(filelist)
		return nil
	},
}

func init() {
	fileCmd.AddCommand(filelistCmd)
	filelistCmd.Flags().BoolP("autocomplete", "a", false, "only for use with autocomplete script")
	filelistCmd.Flags().MarkHidden("autocomplete")
}

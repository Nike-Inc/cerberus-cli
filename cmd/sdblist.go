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

var sdblistCmd = &cobra.Command{
	Use:   "list",
	Short: "list sdbs",
	Long:  `list sdbs`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		category, err := cmd.Flags().GetString("category")
		if err != nil {
			return err
		}

		sdbs, err := cl.SDB().List()
		if err != nil {
			return err
		}

		autocomplete, err := cmd.Flags().GetBool("autocomplete")
		if err != nil {
			return err
		}

		sep := "\n"
		if autocomplete {
			sep = " "
		}

		sdblist := ""
		for _, sdb := range sdbs {
			if category != "" {
				if strings.Index(sdb.Path, category) == 0 {
					sdblist += sdb.Path + sep
				}
			} else {
				sdblist += sdb.Path + sep
			}
		}
		fmt.Println(sdblist)
		return nil
	},
}

func init() {
	sdbCmd.AddCommand(sdblistCmd)
	sdblistCmd.Flags().StringP("category", "c", "", "the category of sdb to list")
	sdblistCmd.Flags().BoolP("autocomplete", "a", false, "only for use with autocomplete script")
	sdblistCmd.Flags().MarkHidden("autocomplete")
}

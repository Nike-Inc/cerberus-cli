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
)

var categorylistCmd = &cobra.Command{
	Use:   "list",
	Short: "list Cerberus categories",
	Long:  `list Cerberus categories`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		categories, err := cl.Category().List()
		if err != nil {
			return err
		}
		categorylist := ""
		for _, category := range categories {
			categorylist += category.Path + "/ "
		}
		fmt.Println(categorylist)
		return nil
	},
}

func init() {
	categoryCmd.AddCommand(categorylistCmd)
}

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

var secretlistCmd = &cobra.Command{
	Use:   "list <secure data path>",
	Short: "list secrets of a secure data path",
	Long:  `list secrets of a secure data path`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		path := args[0]
		if strings.LastIndex(path, "/") != len(path)-1 {
			path += "/"
		}

		autocomplete, err := cmd.Flags().GetBool("autocomplete")
		if err != nil {
			return err
		}

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		secret, err := cl.Secret().List(path)
		if err != nil {
			return err
		}

		if secret == nil {
			return fmt.Errorf("path %s does not exist", path)
		}

		keys := secret.Data["keys"]
		s := keys.([]interface{})

		if autocomplete {
			secretlist := ""
			for _, key := range s {
				secretlist += path + key.(string) + "\\n"
			}
			fmt.Print(secretlist)
		} else {
			for _, key := range s {
				fmt.Println(key)
			}
		}

		return nil
	},
}

func init() {
	secretCmd.AddCommand(secretlistCmd)
	secretlistCmd.Flags().BoolP("autocomplete", "a", false, "only for use with autocomplete script")
	secretlistCmd.Flags().MarkHidden("autocomplete")
}

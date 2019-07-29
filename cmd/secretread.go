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
	"cerberus-cli/tool"
	"fmt"
	"github.com/spf13/cobra"
)

var secretreadCmd = &cobra.Command{
	Use:   "read <secure data path>",
	Short: "read a specific secret",
	Long:  `read a specific secret`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		path := args[0]
		err := GetSecretWithFullPath(path)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	secretCmd.AddCommand(secretreadCmd)
}

func GetSecretWithFullPath(path string) error {
	cl, err := client.GetClient()
	if err != nil {
		return err
	}

	secret, err := cl.Secret().Read(path)
	if err != nil || secret == nil {
		return fmt.Errorf("could not read path %s: %v", path, err)
	}

	output, err := tool.ToJSON(secret.Data)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

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

var secretdeleteCmd = &cobra.Command{
	Use:   "delete <secure data path>",
	Short: "delete a specific secret",
	Long:  `delete a specific secret`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		path := args[0]
		err := DeleteSecretWithFullPath(path)
		if err != nil {
			return err
		}
		fmt.Printf("Successfully deleted %s\n", path)
		return nil
	},
}

func init() {
	secretCmd.AddCommand(secretdeleteCmd)
}

func DeleteSecretWithFullPath(path string) error {
	cl, err := client.GetClient()
	if err != nil {
		return err
	}

	_, err = cl.Secret().Delete(path)
	if err != nil {
		return fmt.Errorf("Could not delete %s: %v", path, err)
	}
	return nil
}

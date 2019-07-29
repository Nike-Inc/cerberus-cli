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
	"net/http"
	"path"
)

var filedeleteCmd = &cobra.Command{
	Use:   "delete <secure file path>",
	Short: "delete a specific file",
	Long:  "delete a specific file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		secureFilePath := args[0]

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		resp, err := cl.DoRequest(http.MethodDelete,
			path.Join("/v1/secure-file", secureFilePath),
			map[string]string{},
			nil)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return fmt.Errorf("error while deleting secure file: %v", err)
		}

		if resp.StatusCode == http.StatusNoContent {
			fmt.Printf("Successfully deleted secure file %s\n", secureFilePath)
			return nil
		} else if resp.StatusCode == http.StatusBadRequest {
			return fmt.Errorf("error while deleting secure file %s. Got HTTP status code %d. "+
				"Check for write permissions",
				secureFilePath,
				resp.StatusCode)
		} else {
			return fmt.Errorf("error while deleting secure file %s. Got HTTP status code %d",
				secureFilePath,
				resp.StatusCode)
		}
	},
}

func init() {
	fileCmd.AddCommand(filedeleteCmd)
}

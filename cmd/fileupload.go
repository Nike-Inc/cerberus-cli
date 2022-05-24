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
	"os"
	"path"
)

var fileuploadCmd = &cobra.Command{
	Use:   "upload <destination secure file path> <local source filepath>",
	Short: "upload a specific file to an existing or new secure file path",
	Long:  "upload a specific file to an existing or new secure file path",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		localFilePath := args[1]
		destinationFilePath := args[0]

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		input, err := os.Open(localFilePath)
		if err != nil {
			return fmt.Errorf("failed to open input file: %v", err)
		}
		defer input.Close()

		filename := path.Base(localFilePath)
		err = cl.SecureFile().Put(destinationFilePath, filename, input)
		if err != nil {
			return fmt.Errorf("failed to upload file %s to path %s: %v\n",
				localFilePath,
				destinationFilePath,
				err)
		}

		fmt.Printf("Successfully uploaded file %s to path %s\n",
			localFilePath,
			destinationFilePath)
		return nil
	},
}

func init() {
	fileCmd.AddCommand(fileuploadCmd)
}

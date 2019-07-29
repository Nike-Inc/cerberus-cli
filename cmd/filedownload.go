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
	"path/filepath"
)

var filedownloadCmd = &cobra.Command{
	Use:   "download <secure file path>",
	Short: "download a specific file",
	Long:  `download a specific file. If successful, the download path is printed to terminal.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}

		path := args[0]

		returnpath, err := DownloadFileWithFullPath(path, filepath)
		if err != nil {
			return err
		}

		if returnpath != "" {
			fmt.Println(returnpath)
		}
		return nil
	},
}

func init() {
	fileCmd.AddCommand(filedownloadCmd)
	filedownloadCmd.Flags().StringP("output", "o", "", "(Optional) complete path of a file to write to. If the file does not already exist it will be created.")
}

func DownloadFileWithFullPath(path string, pathToSaveTo string) (string, error) {
	cl, err := client.GetClient()
	if err != nil {
		return "", err
	}

	var output *os.File

	if pathToSaveTo != "" {
		output, err = os.Create(pathToSaveTo)
		if err != nil {
			return "", err
		}
	} else {
		_, filename := filepath.Split(path)
		if filename == "" {
			return "", fmt.Errorf("invalid Cerberus secure file path given")
		}
		pathToSaveTo = "./" + filename
		output, err = os.Create(pathToSaveTo)
		if err != nil {
			return "", err
		}
	}

	defer output.Close()

	err = cl.SecureFile().Get(path, output)
	if err != nil {
		newErr := os.Remove(pathToSaveTo)
		if newErr != nil {
			return "", newErr
		}
		return "", err
	} else {
		absolutePath, err := filepath.Abs(pathToSaveTo)
		if err != nil {
			return pathToSaveTo, nil
		}
		return absolutePath, nil
	}
}

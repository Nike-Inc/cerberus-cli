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
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var fileeditCmd = &cobra.Command{
	Use:   "edit <secure file path>",
	Short: "edit a specific file inline",
	Long:  `edit a specific file inline`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		path := args[0]

		editor, err := cmd.Flags().GetString("editor")
		if err != nil {
			return err
		}

		if editor == "" {
			editor = "vim"
		} else {
			editor = strings.ToLower(editor)
		}

		err = EditFileWithFullPath(path, editor)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	fileCmd.AddCommand(fileeditCmd)
	fileeditCmd.Flags().StringP("editor", "e", tool.GetEnvVariable(tool.EnvPrefEditor), "(Optional) editor to use / set the CERBERUS_EDITOR env variable")
}

func EditFileWithFullPath(path string, editor string) error {
	cl, err := client.GetClient()
	if err != nil {
		return err
	}

	tempDir := os.TempDir()

	suffix := ""
	i := strings.LastIndex(path, ".")
	if i != -1 {
		suffix = path[i:]
	}

	tempfile, err := tool.TempFile(tempDir, "cerberus_temp_", suffix)

	defer os.Remove(tempfile.Name())

	err = cl.SecureFile().Get(path, tempfile)
	if err != nil {
		return err
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = cleanup(tempfile)
		os.Exit(1)
	}()

	fmt.Printf("File temporariliy saved to: %s\n", tempfile.Name())

	var myCmd *exec.Cmd
	if editor != "vi" && editor != "vim" && editor != "nano" && editor != "emacs" {
		myCmd = exec.Command(editor, "--wait", tempfile.Name())
		fmt.Println("A --wait flag has been added to your editor command. If supported, edits will be uploaded upon saving and closing the file.")
	} else {
		myCmd = exec.Command(editor, tempfile.Name())
	}

	myCmd.Stdin = os.Stdin
	myCmd.Stdout = os.Stdout
	myCmd.Stderr = os.Stderr

	err = myCmd.Run()
	if err != nil {
		return err
	}

	err = UploadFileWithFullPath(path, tempfile.Name())
	if err != nil {
		return err
	}
	fmt.Printf("Edits successfully uploaded to %s\n", path)
	return nil
}

// Helper function to upload a file
// path string: the secure file path in Cerberus
// file string: the name of the existing file to upload
func UploadFileWithFullPath(path string, file string) error {
	cl, err := client.GetClient()
	if err != nil {
		return err
	}

	_, nameOfFile := filepath.Split(path)
	if nameOfFile == "/" || nameOfFile == "." || nameOfFile == "" {
		return fmt.Errorf("invalid path given. Is the filename included in the path?")
	}

	input, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	err = cl.SecureFile().Put(path, nameOfFile, input)
	if err != nil {
		return err
	}

	return nil
}

// called when system interrupt is handled
func cleanup(myFile *os.File) error {
	err := os.Remove(myFile.Name())
	if err == nil {
		fmt.Printf("Editor unexpectedly aborted. No edits uploaded. Deleted temporary file %s\n", myFile.Name())
	}
	return err
}

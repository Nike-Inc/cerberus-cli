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
	"bufio"
	"cerberus-cli/client"
	"cerberus-cli/tool"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var secreteditCmd = &cobra.Command{
	Use:   "edit <secure data path>",
	Short: "edit a specific secret",
	Long:  `edit a specific secret`,
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
		output, err := GetSecretStringWithFullPath(path)
		if output == "client_error" {
			return err
		}
		if err != nil {
			fmt.Printf("Given path does not exist. Attempt to create new secret at path %s? [y/N] ", path)
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			text := strings.ToLower(scanner.Text())
			if text == "yes" || text == "y"{
				output = ""
			} else if text == "no" || text == "n" || text == "" {
				return nil
			} else {
				return fmt.Errorf("Invalid option")
			}
		}

		err = EditSecretAndUpload(output, editor, path, nil)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	secretCmd.AddCommand(secreteditCmd)
	secreteditCmd.Flags().StringP("editor", "e", tool.GetEnvVariable(tool.EnvPrefEditor),
		"(Optional) editor to use / set the CERBERUS_EDITOR env variable")
}

func GetSecretStringWithFullPath(path string) (string, error) {
	cl, err := client.GetClient()
	if err != nil {
		return "client_error", err
	}

	secret, err := cl.Secret().Read(path)
	if err != nil || secret == nil {
		return "", fmt.Errorf("could not read path %s: %v", path, err)
	}

	output, err := tool.ToYAML(secret.Data)
	if err != nil {
		return "", err
	}
	return output, nil
}

func EditSecretAndUpload(secret string, editor string, path string, myFile *os.File) error {
	var tempfile *os.File
	scanner := bufio.NewScanner(os.Stdin)
	if myFile == nil {
		tempDir := os.TempDir()

		suffix := ".yaml"
		f, err := tool.TempFile(tempDir, "cerberus_temp_", suffix)
		if err != nil {
			return err
		}
		tempfile = f
		defer os.Remove(tempfile.Name())

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			_ = cleanup(tempfile)
			os.Exit(1)
		}()

		message := "editing an existing secret"
		if len(secret) == 0 {
			message = "creating a new secret"
		}
		data := []byte(secret)
		doc := "# You are now " + message + " at " + path + "\n" +
			"# Make any edits to your secret in this file.\n" +
			"# Save and close the file to upload your edits.\n" +
			"# Follow the YAML guidelines here: https://github.com/go-yaml/yaml#compatibility\n"
		header := []byte(doc)

		_, err = tempfile.Write(header)
		if err != nil {
			return fmt.Errorf("could not write to temporary file: %v", err)
		}
		_, err = tempfile.Write(data)
		if err != nil {
			return fmt.Errorf("could not write to temporary file: %v", err)
		}
		fmt.Printf("Secret temporariliy saved to: %s\n", tempfile.Name())
	} else {
		tempfile = myFile
	}

	var myCmd *exec.Cmd
	if editor != "vi" && editor != "vim" && editor != "nano" && editor != "emacs" {
		myCmd = exec.Command(editor, "--wait", tempfile.Name())
		fmt.Println("A --wait flag has been added to your editor command. If supported, edits will be uploaded " +
			"upon saving and closing the file.")
	} else {
		myCmd = exec.Command(editor, tempfile.Name())
	}

	myCmd.Stdin = os.Stdin
	myCmd.Stdout = os.Stdout
	myCmd.Stderr = os.Stderr

	err := myCmd.Run()
	if err != nil {
		return err
	}

	input, err := os.Open(tempfile.Name())
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	edits, err := ioutil.ReadAll(input)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var kvpairs map[string]interface{}
	err = yaml.Unmarshal(edits, &kvpairs)
	if err != nil {
		fmt.Printf("Failed to parse edited secret. Try editing again? [y/N] ")
		scanner.Scan()
		text := strings.ToLower(scanner.Text())
		if text == "yes" || text == "y" {
			return EditSecretAndUpload("", editor, path, tempfile)
		} else if text == "no" || text == "n" || text == "" {
			fmt.Println("Edit secret aborted. Deleted temporary file.")
			return nil
		} else {
			return fmt.Errorf("failed to parse edited secret: %v", err)
		}
	}

	if len(kvpairs) > 0 {
		err = WriteSecret(path, kvpairs)
		if err != nil {
			fmt.Printf("Failed to write secret to path %s: %v. Try again? [y/N] ", path, err)
			scanner.Scan()
			text := strings.ToLower(scanner.Text())
			if text == "yes" || text == "y" {
				return EditSecretAndUpload("", editor, path, tempfile)
			} else if text == "no" || text == "n" || text == "" {
				fmt.Println("Edit secret aborted. Deleted temporary file.")
				return nil
			} else {
				return fmt.Errorf("failed to write secret to path %s: %v", path, err)
			}
		}
	} else {
		fmt.Printf("No key/value pairs to write. Try again? [y/N] ")
		scanner.Scan()
		text := strings.ToLower(scanner.Text())
		if text == "yes" || text == "y" {
			return EditSecretAndUpload("", editor, path, tempfile)
		} else if text == "no" || text == "n" || text == "" {
			fmt.Println("Edit secret aborted. Deleted temporary file.")
			return nil
		} else {
			return fmt.Errorf("no key/value pairs to write")
		}
	}

	fmt.Printf("Edits successfully uploaded to path %s\n", path)
	return nil
}

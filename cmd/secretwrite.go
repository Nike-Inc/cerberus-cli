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
	"github.com/Nike-Inc/cerberus-go-client/v3/cerberus"
	"github.com/spf13/cobra"
	"strings"
)

const EXPECTED_KV_PAIR = 2

var secretwriteCmd = &cobra.Command{
	Use:   "write <secure data path> <entry> ... ",
	Short: "write a secret to a specific secure data path",
	Long:  `write a secret to a specific secure data path`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		path := args[0]

		entries, err := cmd.Flags().GetStringArray("entry")
		if err != nil {
			return err
		}

		merge, err := cmd.Flags().GetBool("merge")
		if err != nil {
			return err
		}

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		data, err := getFinalData(cl, path, entries, merge)
		if err != nil {
			return err
		}

		err = SecretWrite(cl, path, data)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully wrote secrets to %s\n", path)
		return nil
	},
}

func init() {
	secretCmd.AddCommand(secretwriteCmd)
	secretwriteCmd.Flags().StringArrayP("entry", "e", []string{}, "use this flag for each "+
		"entry to add to the new secret, in the required format of KEY=VALUE")
	secretwriteCmd.MarkFlagRequired("entry")
	secretwriteCmd.Flags().BoolP("merge", "m", false, "merge key-value pairs instead "+
		"of overwriting, off by default")
}

func WriteSecret(path string, data map[string]interface{}) error {
	cl, err := client.GetClient()
	if err != nil {
		return err
	}

	return SecretWrite(cl, path, data)
}

func SecretWrite(cl *cerberus.Client, path string, data map[string]interface{}) error {
	_, err := cl.Secret().Write(path, data)
	if err != nil {
		return err
	}
	return nil
}

func mapEntries(entries []string,
	data map[string]interface{}) (map[string]interface{}, error) {

	for _, item := range entries {
		split := strings.SplitN(item, "=", EXPECTED_KV_PAIR)
		if len(split) != EXPECTED_KV_PAIR {
			return nil, fmt.Errorf("entry %s not in expected format (key=value)", item)
		}
		data[split[0]] = split[1]
	}
	return data, nil
}

func getFinalData(cl *cerberus.Client, path string,
	entries []string, merge bool) (map[string]interface{}, error) {

	existingSecretData := make(map[string]interface{})
	if merge {
		existingSecret, err := cl.Secret().Read(path)
		if err != nil {
			return nil, err
		}
		if existingSecret != nil {
			existingSecretData = existingSecret.Data
		}
	}

	data, err := mapEntries(entries, existingSecretData)
	if err != nil {
		return nil, err
	}

	return data, err
}

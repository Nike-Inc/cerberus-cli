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
	"github.com/Nike-Inc/cerberus-go-client/api"
	"github.com/Nike-Inc/cerberus-go-client/cerberus"
	"github.com/spf13/cobra"
	"strings"
)

var sdbdeleteCmd = &cobra.Command{
	Use:   "delete <path of SDB>",
	Short: "delete an existing SDB",
	Long:  `create an existing SDB`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		path := args[0]

		sdb, err := GetSDBWithFullPath(cl, path)
		if err != nil {
			return err
		}

		err = cl.SDB().Delete(sdb.ID)
		if err != nil {
			return fmt.Errorf("delete SDB %s failed: %s", path, err)
		}
		fmt.Printf("Successfully deleted SDB %s\n", path)

		return nil
	},
}

func init() {
	sdbCmd.AddCommand(sdbdeleteCmd)
}

func GetSDBWithFullPath(cl *cerberus.Client, path string) (*api.SafeDepositBox, error) {
	sdbList, err := cl.SDB().List()
	if err != nil {
		return nil, err
	}
	for _, sdb := range sdbList {
		if strings.TrimSuffix(sdb.Path, "/") == strings.TrimSuffix(path, "/") {
			return sdb, nil
		}
	}
	return nil, fmt.Errorf("could not find SDB %s,", path)
}

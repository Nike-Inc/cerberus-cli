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
	"github.com/zalando/go-keyring"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and clear any existing auth session saved in keyring",
	Long:  "Logout and clear any existing auth session saved in keyring",
	RunE: func(cmd *cobra.Command, args []string) error {
		err1 := keyring.Delete(client.SERVICE, client.CERBTOKEN)
		if err1 != nil {
			fmt.Println("Token not found in keyring")
		}
		err2 := keyring.Delete(client.SERVICE, client.EXPIRYTIME)
		if err2 != nil {
			fmt.Println("Expiry time not found in keyring")
		}
		err3 := keyring.Delete(client.SERVICE, client.CERBURL)
		if err3 != nil {
			fmt.Println("Url not found in keyring")
		}

		if err1 == nil && err2 == nil && err3 == nil{
			fmt.Println("Session cleared")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

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

var tokenStsAuthCmd = &cobra.Command{
	Use:   "sts",
	Short: "Get a Cerberus token with AWS STS auth",
	Long:  `Get a Cerberus token with AWS STS auth`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cl, err := client.GetClient()
		if err != nil {
			return fmt.Errorf("err: %s", err)
		}else{
			token, err := cl.Authentication.GetToken(nil)
			if err != nil{
				return fmt.Errorf("Could not retrive client token! err: %s", err)
			}else{
				fmt.Println(token)
			}
		}
		return nil
	},
}

func init() {
	tokenCmd.AddCommand(tokenStsAuthCmd)
}

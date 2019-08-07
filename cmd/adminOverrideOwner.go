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
	"bytes"
	"cerberus-cli/client"
	"encoding/json"
	"fmt"
	"github.com/Nike-Inc/cerberus-go-client/api"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var adminOverrideOwner = &cobra.Command{
	Use:   "override-owner <secure file path>",
	Short: "override the owner of an SDB",
	Long:  "override the owner of an SDB",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		sdbname, err := cmd.Flags().GetString("sdb")
		if err != nil {
			return err
		}

		newowner, err := cmd.Flags().GetString("owner")
		if err != nil {
			return err
		}

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		params := map[string]string{
			"sdbName": sdbname,
		}

		resp, err := cl.DoRequest(http.MethodGet, "/v1/metadata", params, nil)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			return err
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error while retrieving metadata for SDB %s. Got HTTP status code %d", sdbname, resp.StatusCode)
		}

		p, err := ioutil.ReadAll(resp.Body)
		var metadataResp = &api.MetadataResponse{}
		err = json.Unmarshal(p, metadataResp)
		if err != nil {
			return err
		}

		metadata := metadataResp.Metadata[0]
		curowner := metadata.Owner
		createdBy := metadata.CreatedBy
		createdTs := metadata.Created
		description := metadata.Description
		lastupdatedBy := metadata.LastUpdatedBy
		lastupdatedTs := metadata.LastUpdated

		fmt.Printf("SDB Name: %s\n" +
			"Description: %s\n" +
			"Current Owner: %s\n" +
			"Created by: %s on %s\n" +
			"Last updated by %s on %s\n\n",
			sdbname, description, curowner, createdBy, createdTs, lastupdatedBy, lastupdatedTs)

		fmt.Printf("You are about to change the owner of SDB '%s' from '%s' to '%s'. Are you sure? [y/N] ",
			sdbname, curowner, newowner)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := strings.ToLower(scanner.Text())
		if text == "yes" || text == "y" {
			fmt.Println("Continuing with override of SDB owner")
		} else if text == "no"  || text == "n" || text == ""{
			fmt.Println("Cancelling override of SDB owner")
			return nil
		} else {
			return fmt.Errorf("Invalid option")
		}

		requestBody, err := json.Marshal(map[string]string{
			"name": sdbname,
			"owner": newowner,
		})

		resp, err = cl.DoRequestWithBody(http.MethodPut, "/v1/admin/override-owner", nil, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("error while trying to override owner of SDB %s to %s. Got HTTP status code %d",
				sdbname, newowner, resp.StatusCode)
		}
		fmt.Println("Successfully overrode SDB owner")
		return nil
	},
}

func init() {
	adminCmd.AddCommand(adminOverrideOwner)
	adminOverrideOwner.Flags().StringP("sdb", "s", "", "SDB name")
	adminOverrideOwner.Flags().StringP("owner", "o", "", "new SDB owner")
	_ = adminOverrideOwner.MarkFlagRequired("sdb")
	_ = adminOverrideOwner.MarkFlagRequired("owner")
}

func parseResponse(r io.Reader, parseTo interface{}) error {
	// Decode the body into the provided interface
	return json.NewDecoder(r).Decode(parseTo)
}

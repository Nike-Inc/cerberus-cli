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
	"github.com/Nike-Inc/cerberus-go-client/v2/api"
	"github.com/spf13/cobra"
	"path"
	"strings"
)

var (
	name         string
	owner        string
	category     string
	description  string
	usergroup    []string
	iamprincipal []string
)

const EXPECTED_NAME_ROLE_PAIR = 2

var sdbcreateCmd = &cobra.Command{
	Use:   "create <path of new SDB>",
	Short: "create a new SDB",
	Long:  `create a new SDB`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		cl, err := client.GetClient()
		if err != nil {
			return err
		}

		sdbpath := path.Join(category, name)

		categoryId := ""
		categories, err := cl.Category().List()
		if err != nil {
			return err
		}
		for _, c := range categories {
			if category == c.Path {
				categoryId = c.ID
			}
		}
		if categoryId == "" {
			return fmt.Errorf("category %s not found", category)
		}

		userGroupPermissions, err := mapUserGroupPermissions(usergroup)
		if err != nil {
			return err
		}

		iamPrincipalPermissions, err := mapIAMPrincipalPermissions(iamprincipal)
		if err != nil {
			return err
		}

		newSDB := &api.SafeDepositBox{
			Name:                    name,
			Path:                    sdbpath,
			CategoryID:              categoryId,
			Description:             description,
			Owner:                   owner,
			UserGroupPermissions:    userGroupPermissions,
			IAMPrincipalPermissions: iamPrincipalPermissions,
		}
		_, err = cl.SDB().Create(newSDB)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully created SDB %s\n", sdbpath)
		return nil
	},
}

func init() {
	sdbCmd.AddCommand(sdbcreateCmd)

	sdbcreateCmd.Flags().StringVarP(&category, "category", "c", "", "SDB category")
	sdbcreateCmd.Flags().StringVarP(&name, "name", "n", "", "name of the SDB")
	sdbcreateCmd.Flags().StringVarP(&owner, "owner", "o", "", "owner of the SDB")
	sdbcreateCmd.MarkFlagRequired("category")
	sdbcreateCmd.MarkFlagRequired("name")
	sdbcreateCmd.MarkFlagRequired("owner")

	sdbcreateCmd.Flags().StringVarP(&description, "description", "d", "",
		"(Optional) description of the SDB")
	sdbcreateCmd.Flags().StringArrayVarP(&usergroup, "usergroup", "g", []string{},
		"(Optional) use this flag for each user group permission to add, in the required format of "+
			"'<UserGroup>,<read/write/owner>'")
	sdbcreateCmd.Flags().StringArrayVarP(&iamprincipal, "iam", "i", []string{},
		"(Optional) use this flag for each IAM Principal permission to add, in the required format of "+
			"'<IAM Principal ARN>,<read/write/owner>'")
	sdbcreateCmd.Flags().SortFlags = false
}

func mapUserGroupPermissions(usergroups []string) ([]api.UserGroupPermission, error) {
	userGroupPermissions := make([]api.UserGroupPermission, 0, len(usergroups))
	myMap, err := createNameToRoleMap(usergroups)
	if err != nil {
		return nil, err
	}

	for name, roleID := range myMap {
		userGroupPermissions = append(userGroupPermissions,
			api.UserGroupPermission{
				Name:   name,
				RoleID: roleID,
			})
	}

	return userGroupPermissions, nil
}

func mapIAMPrincipalPermissions(iamprincipals []string) ([]api.IAMPrincipal, error) {
	iamPrincipalPermissions := make([]api.IAMPrincipal, 0, len(iamprincipals))
	myMap, err := createNameToRoleMap(iamprincipals)
	if err != nil {
		return nil, err
	}

	for arn, roleID := range myMap {
		iamPrincipalPermissions = append(iamPrincipalPermissions,
			api.IAMPrincipal{
				IAMPrincipalARN: arn,
				RoleID:          roleID,
			})
	}
	return iamPrincipalPermissions, nil
}

func createNameToRoleMap(entries []string) (map[string]string, error) {
	data := make(map[string]string)
	roleMap, err := getRoleMap()
	if err != nil {
		return nil, err
	}

	for _, item := range entries {
		split := strings.SplitN(item, ",", EXPECTED_NAME_ROLE_PAIR)
		if len(split) != EXPECTED_NAME_ROLE_PAIR {
			return nil, fmt.Errorf("entry %s not in expected format '<name>,<role>'", item)
		}
		name, role := split[0], split[1]
		if role == "read" || role == "write" || role == "owner" {
			roleID, present := roleMap[role]
			if !present {
				return nil, fmt.Errorf("could not find role ID for %s", role)
			}
			data[name] = roleID
		} else {
			return nil, fmt.Errorf("role not valid. Must be read/write")
		}
	}
	return data, nil
}

func getRoleMap() (map[string]string, error) {
	cl, err := client.GetClient()
	if err != nil {
		return nil, err
	}

	roles, err := cl.Role().List()
	if err != nil {
		return nil, err
	}

	roleMap := make(map[string]string)

	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}

	return roleMap, nil
}

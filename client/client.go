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

package client

import (
	"fmt"
	"github.com/Nike-Inc/cerberus-go-client/auth"
	"github.com/Nike-Inc/cerberus-go-client/cerberus"
	"github.com/zalando/go-keyring"
	"time"
)

var Region string
var Token string
var Url string

const SERVICE string = "cerberus-cli"
const CERBTOKEN string = "X-Cerberus-Token"
const EXPIRYTIME string = "Token-Expiry-Time"
const LAYOUT string = "2006-01-02 15:04:05.999999999 -0700 MST"

func GetClient() (*cerberus.Client, error) {
	cerbToken, err := keyring.Get(SERVICE, CERBTOKEN)
	if err != nil {
		// token doesn't exist
		return authenticate()
	} else {
		// token exists
		expiry, err := keyring.Get(SERVICE, EXPIRYTIME)
		if err != nil {
			return authenticate()
		}

		expiryTime, err := time.Parse(LAYOUT, expiry)
		if err != nil {
			return authenticate()
		}
		currentTime := time.Now()

		// token has expired
		if currentTime.After(expiryTime) {
			_ = keyring.Delete(SERVICE, EXPIRYTIME)
			_ = keyring.Delete(SERVICE, CERBTOKEN)
			return authenticate()
		}

		// try auth with existing token
		tokenAuth, err := auth.NewTokenAuth(Url, cerbToken)
		if err != nil {
			return authenticate()
		}

		// try getting client with auth
		cl, err := cerberus.NewClient(tokenAuth, nil)
		if err != nil {
			return authenticate()
		}

		// auth and get client with existing token was successful
		return cl, nil
	}
}

func authenticate() (*cerberus.Client, error) {
	var authMethod auth.Auth
	var err error

	if Token != "" {
		if authMethod, err = auth.NewTokenAuth(Url, Token); err != nil {
			return nil, err
		}
	} else if Region != "" {
		if authMethod, err = auth.NewSTSAuth(Url, Region); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("no authentication provided")
	}

	// try auth
	cl, err := cerberus.NewClient(authMethod, nil)
	if err != nil {
		return nil, err
	}

	// successful auth, save token to keyring
	if cl != nil {
		// ignore error that might occur on Linux with keyring
		_ = saveTokenToKeyring(cl)
		return cl, nil
	}
	return nil, err
}

func saveTokenToKeyring(cl *cerberus.Client) error {
	tok, err := cl.Authentication.GetToken(nil)
	if err != nil {
		return err
	}

	exp, err := cl.Authentication.GetExpiry()
	if err != nil {
		return err
	}

	err = keyring.Set(SERVICE, CERBTOKEN, tok)
	if err != nil {
		return err
	}

	err = keyring.Set(SERVICE, EXPIRYTIME, exp.Format(LAYOUT))
	if err != nil {
		_ = keyring.Delete(SERVICE, CERBTOKEN)
		return err
	}
	return nil
}
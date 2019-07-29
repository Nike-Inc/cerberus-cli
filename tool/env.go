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

package tool

import (
	"os"
)

const EnvCerbToken = "CERBERUS_TOKEN"

const EnvCerbUrl = "CERBERUS_URL"

const EnvCerbRegion = "CERBERUS_REGION"

const EnvPrefEditor = "CERBERUS_EDITOR"

func GetEnvVariable(variable string) string {
	envVariable := os.Getenv(variable)
	if len(envVariable) == 0 {
		return ""
	}
	return envVariable
}
